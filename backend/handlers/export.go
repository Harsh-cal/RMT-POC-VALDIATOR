package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harsh-cal/rmt-poc-validator/models"
	"github.com/harsh-cal/rmt-poc-validator/services"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var exportIssueTypeLabel = map[string]string{
	"version_mismatch":                "Version Mismatch",
	"missing_dependency":              "Missing Dependency",
	"duplicate":                       "Duplicate Container",
	"unsupported_combo":               "Unsupported Combination",
	"duplicate_part_number":           "Duplicate Part Number",
	"system_incompatible":             "System Incompatible",
	"aircraft_state_version_conflict": "Aircraft State Version Conflict",
	"version_conflict":                "Version Conflict",
	"maturity_risk":                   "Maturity Risk",
}

// ExportValidationHandler handles POST /validate/export requests.
func ExportValidationHandler(c *gin.Context) {
	var req models.ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	format := strings.ToLower(strings.TrimSpace(req.Format))
	if format != "csv" && format != "pdf" && format != "xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format must be one of: csv, pdf, xlsx"})
		return
	}

	releaseID := strings.TrimSpace(req.ReleaseID)
	releaseName := strings.TrimSpace(req.ReleaseName)
	if releaseID == "" && releaseName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "release_id or release_name is required"})
		return
	}

	var (
		result *models.ValidationResult
		err    error
	)

	if releaseID != "" {
		result, err = services.GetValidationResult(releaseID)
		if err == mongo.ErrNoDocuments && releaseName != "" {
			result, err = services.GetLatestValidationResultByReleaseName(releaseName)
		}
	} else {
		result, err = services.GetLatestValidationResultByReleaseName(releaseName)
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "validation result not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch validation result"})
		return
	}

	fileStem := buildFileStem(result.ReleaseName)

	switch format {
	case "csv":
		content, fileName, err := buildCSV(result, fileStem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate csv"})
			return
		}
		streamFile(c, "text/csv", fileName, content)
	case "xlsx":
		content, fileName, err := buildXLSX(result, fileStem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate xlsx"})
			return
		}
		streamFile(c, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileName, content)
	case "pdf":
		content, fileName, err := buildPDF(result, fileStem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate pdf"})
			return
		}
		streamFile(c, "application/pdf", fileName, content)
	}
}

func streamFile(c *gin.Context, contentType, fileName string, content []byte) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, contentType, content)
}

func buildFileStem(releaseName string) string {
	normalized := strings.ToLower(strings.TrimSpace(releaseName))
	normalized = strings.ReplaceAll(normalized, " ", "-")
	if normalized == "" {
		normalized = "release"
	}
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-validation-%s", normalized, timestamp)
}

func buildSummaryRows(result *models.ValidationResult) [][2]string {
	high := 0
	medium := 0
	low := 0
	for _, issue := range result.Issues {
		switch issue.Severity {
		case "HIGH":
			high++
		case "MEDIUM":
			medium++
		case "LOW":
			low++
		}
	}

	validatedAt := "-"
	if !result.ValidatedAt.IsZero() {
		validatedAt = result.ValidatedAt.Local().Format(time.RFC1123)
	}

	return [][2]string{
		{"Release", result.ReleaseName},
		{"Release ID", result.ReleaseID},
		{"Version", result.Version},
		{"Target Fleet", result.TargetFleet},
		{"Status", result.Status},
		{"Risk", result.Risk},
		{"Validated At", validatedAt},
		{"Total Issues", fmt.Sprintf("%d", len(result.Issues))},
		{"High Issues", fmt.Sprintf("%d", high)},
		{"Medium Issues", fmt.Sprintf("%d", medium)},
		{"Low Issues", fmt.Sprintf("%d", low)},
	}
}

func getRecommendation(issueType string, recs []models.Recommendation) string {
	for _, rec := range recs {
		if rec.IssueType == issueType {
			return rec.Action
		}
	}
	return "-"
}

func formatIssueType(issueType string) string {
	if label, ok := exportIssueTypeLabel[issueType]; ok {
		return label
	}
	return issueType
}

func buildCSV(result *models.ValidationResult, fileStem string) ([]byte, string, error) {
	buf := bytes.NewBuffer(nil)
	writer := csv.NewWriter(buf)

	if err := writer.Write([]string{"Summary", "Value"}); err != nil {
		return nil, "", err
	}
	for _, row := range buildSummaryRows(result) {
		if err := writer.Write([]string{row[0], row[1]}); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Write([]string{}); err != nil {
		return nil, "", err
	}

	headers := []string{"#", "Type", "Severity", "Message", "Recommendation"}
	if err := writer.Write(headers); err != nil {
		return nil, "", err
	}

	for idx, issue := range result.Issues {
		recommendation := getRecommendation(issue.Type, result.Recommendations)
		if err := writer.Write([]string{
			fmt.Sprintf("%d", idx+1),
			formatIssueType(issue.Type),
			issue.Severity,
			issue.Message,
			recommendation,
		}); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), fileStem + ".csv", nil
}

func buildXLSX(result *models.ValidationResult, fileStem string) ([]byte, string, error) {
	f := excelize.NewFile()
	defer f.Close()

	summarySheet := "Summary"
	issuesSheet := "Issues"
	f.SetSheetName("Sheet1", summarySheet)
	f.NewSheet(issuesSheet)

	summaryRows := buildSummaryRows(result)
	for i, row := range summaryRows {
		cellA := fmt.Sprintf("A%d", i+1)
		cellB := fmt.Sprintf("B%d", i+1)
		f.SetCellValue(summarySheet, cellA, row[0])
		f.SetCellValue(summarySheet, cellB, row[1])
	}

	headers := []string{"#", "Type", "Severity", "Message", "Recommendation"}
	for i, h := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(issuesSheet, col+"1", h)
	}

	for idx, issue := range result.Issues {
		row := idx + 2
		f.SetCellValue(issuesSheet, fmt.Sprintf("A%d", row), idx+1)
		f.SetCellValue(issuesSheet, fmt.Sprintf("B%d", row), formatIssueType(issue.Type))
		f.SetCellValue(issuesSheet, fmt.Sprintf("C%d", row), issue.Severity)
		f.SetCellValue(issuesSheet, fmt.Sprintf("D%d", row), issue.Message)
		f.SetCellValue(issuesSheet, fmt.Sprintf("E%d", row), getRecommendation(issue.Type, result.Recommendations))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), fileStem + ".xlsx", nil
}

func buildPDF(result *models.ValidationResult, fileStem string) ([]byte, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(12, 12, 12)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Validation Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	for _, row := range buildSummaryRows(result) {
		pdf.MultiCell(0, 6, fmt.Sprintf("%s: %s", row[0], row[1]), "", "L", false)
	}

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Issues")
	pdf.Ln(9)
	pdf.SetFont("Arial", "", 10)

	if len(result.Issues) == 0 {
		pdf.MultiCell(0, 6, "No issues detected.", "", "L", false)
	} else {
		for idx, issue := range result.Issues {
			rec := getRecommendation(issue.Type, result.Recommendations)
			line := fmt.Sprintf("%d. [%s] %s - %s | Recommendation: %s", idx+1, issue.Severity, formatIssueType(issue.Type), issue.Message, rec)
			pdf.MultiCell(0, 6, line, "", "L", false)
			pdf.Ln(1)
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := pdf.Output(buf); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), fileStem + ".pdf", nil
}
