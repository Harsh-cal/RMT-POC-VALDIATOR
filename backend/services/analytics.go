package services

import (
	"context"
	"fmt"
	"time"

	"github.com/harsh-cal/rmt-poc-validator/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GetReleaseHistory retrieves paginated release validation history
func GetReleaseHistory(limit, offset int) (*models.HistoryResponse, error) {
	if MongoDatabase == nil {
		return nil, fmt.Errorf("MongoDB not initialized")
	}

	collection := MongoDatabase.Collection("validation_results")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get total count
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Fetch releases sorted by validated_at descending
	opts := options.Find().SetSort(bson.M{"validated_at": -1}).SetSkip(int64(offset)).SetLimit(int64(limit))
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	results, err := decodeCursorToValidationResults(cursor, ctx)
	if err != nil {
		return nil, err
	}

	// Convert to history items
	items := make([]models.ReleaseHistoryItem, len(results))
	for i, result := range results {
		topIssues := getTopIssueTypes(result.Issues, 3)
		items[i] = models.ReleaseHistoryItem{
			ReleaseID:   result.ReleaseID,
			ReleaseName: result.ReleaseName,
			Version:     result.Version,
			Status:      result.Status,
			Risk:        result.Risk,
			IssueCount:  len(result.Issues),
			TopIssues:   topIssues,
			ValidatedAt: result.ValidatedAt,
			ValidatedBy: "System", // Could track user later
		}
	}

	// Calculate metrics
	metrics := calculateHistoryMetrics(results)

	return &models.HistoryResponse{
		Total:    int(total),
		Releases: items,
		Metrics:  metrics,
	}, nil
}

// GetTrendData retrieves daily trend statistics for the past N days
func GetTrendData(days int) (*models.TrendResponse, error) {
	if MongoDatabase == nil {
		return nil, fmt.Errorf("MongoDB not initialized")
	}

	collection := MongoDatabase.Collection("validation_results")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch all results from the past N days
	since := time.Now().AddDate(0, 0, -days)
	cursor, err := collection.Find(ctx, bson.M{
		"validated_at": bson.M{"$gte": since},
	})
	if err != nil {
		return nil, err
	}

	results, err := decodeCursorToValidationResults(cursor, ctx)
	if err != nil {
		return nil, err
	}

	// Group by date
	trendMap := make(map[string]*models.TrendData)
	for _, result := range results {
		dateKey := result.ValidatedAt.Format("2006-01-02")
		if _, exists := trendMap[dateKey]; !exists {
			trendMap[dateKey] = &models.TrendData{
				Date:     dateKey,
				Passes:   0,
				Failures: 0,
				AvgRisk:  0,
			}
		}

		if result.Status == "PASS" {
			trendMap[dateKey].Passes++
		} else {
			trendMap[dateKey].Failures++
		}

		// Accumulate risk (convert to numeric)
		riskScore := riskToScore(result.Risk)
		trendMap[dateKey].AvgRisk += riskScore
	}

	// Average risk scores
	for _, trend := range trendMap {
		totalCount := trend.Passes + trend.Failures
		if totalCount > 0 {
			trend.AvgRisk = trend.AvgRisk / float64(totalCount)
		}
	}

	// Convert map to sorted array
	trends := make([]models.TrendData, 0, len(trendMap))
	for _, trend := range trendMap {
		trends = append(trends, *trend)
	}

	// Sort by date
	for i := 0; i < len(trends)-1; i++ {
		for j := i + 1; j < len(trends); j++ {
			if trends[i].Date > trends[j].Date {
				trends[i], trends[j] = trends[j], trends[i]
			}
		}
	}

	return &models.TrendResponse{Data: trends}, nil
}

// GetRecurringIssues returns the most common issues in the past N days
func GetRecurringIssues(days int) (*models.IssueRecurrenceResponse, error) {
	if MongoDatabase == nil {
		return nil, fmt.Errorf("MongoDB not initialized")
	}

	collection := MongoDatabase.Collection("validation_results")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	since := time.Now().AddDate(0, 0, -days)
	cursor, err := collection.Find(ctx, bson.M{
		"validated_at": bson.M{"$gte": since},
	})
	if err != nil {
		return nil, err
	}

	results, err := decodeCursorToValidationResults(cursor, ctx)
	if err != nil {
		return nil, err
	}

	// Count issues by type
	issueMap := make(map[string]*models.RecurringIssueStats)
	successMap := make(map[string]int) // Track fixes

	for _, result := range results {
		for _, issue := range result.Issues {
			if _, exists := issueMap[issue.Type]; !exists {
				issueMap[issue.Type] = &models.RecurringIssueStats{
					Type:       issue.Type,
					Count:      0,
					Containers: []string{},
				}
			}
			issueMap[issue.Type].Count++
			issueMap[issue.Type].LastOccurrence = result.ValidatedAt

			// Add container if not already present
			found := false
			for _, c := range issueMap[issue.Type].Containers {
				if c == issue.Container {
					found = true
					break
				}
			}
			if !found && issue.Container != "" {
				issueMap[issue.Type].Containers = append(issueMap[issue.Type].Containers, issue.Container)
			}
		}

		// Track successful validations (implicit fix)
		if result.Status == "PASS" {
			for _, issue := range result.Issues {
				successMap[issue.Type]++
			}
		}
	}

	// Calculate fix rates and convert to array
	issues := make([]models.RecurringIssueStats, 0)
	for _, stat := range issueMap {
		total := stat.Count + successMap[stat.Type]
		if total > 0 {
			stat.FixRate = float64(successMap[stat.Type]) / float64(total)
		}
		issues = append(issues, *stat)
	}

	// Sort by count descending
	for i := 0; i < len(issues)-1; i++ {
		for j := i + 1; j < len(issues); j++ {
			if issues[i].Count < issues[j].Count {
				issues[i], issues[j] = issues[j], issues[i]
			}
		}
	}

	return &models.IssueRecurrenceResponse{Issues: issues}, nil
}

// CompareReleases compares two validation results
func CompareReleases(releaseID1, releaseID2 string) (*models.ReleaseComparisonData, error) {
	result1, err := GetValidationResult(releaseID1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release 1: %v", err)
	}

	result2, err := GetValidationResult(releaseID2)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release 2: %v", err)
	}

	return &models.ReleaseComparisonData{
		Release1:          *result1,
		Release2:          *result2,
		NewContainers:     getNewContainers(result1, result2),
		RemovedContainers: getRemovedContainers(result1, result2),
		UpdatedContainers: getUpdatedContainers(result1, result2),
		FixedIssues:       getFixedIssues(result1, result2),
		NewIssues:         getNewIssues(result1, result2),
	}, nil
}

// Helper functions

func getTopIssueTypes(issues []models.Issue, limit int) []string {
	typeCount := make(map[string]int)
	for _, issue := range issues {
		typeCount[issue.Type]++
	}

	types := make([]string, 0)
	for t := range typeCount {
		types = append(types, t)
	}

	// Limit to top N
	if len(types) > limit {
		types = types[:limit]
	}
	return types
}

func calculateHistoryMetrics(results []models.ValidationResult) models.HistoryMetrics {
	if len(results) == 0 {
		return models.HistoryMetrics{}
	}

	passes := 0
	totalIssues := 0
	highRisk := 0
	mediumRisk := 0
	safe := 0

	for _, result := range results {
		if result.Status == "PASS" {
			passes++
		}
		totalIssues += len(result.Issues)

		switch result.Risk {
		case "HIGH":
			highRisk++
		case "MEDIUM":
			mediumRisk++
		case "SAFE":
			safe++
		}
	}

	passRate := float64(passes) / float64(len(results))
	avgIssues := float64(totalIssues) / float64(len(results))

	return models.HistoryMetrics{
		PassRate:        passRate,
		AvgIssues:       avgIssues,
		HighRiskCount:   highRisk,
		MediumRiskCount: mediumRisk,
		SafeCount:       safe,
	}
}

func riskToScore(risk string) float64 {
	switch risk {
	case "HIGH":
		return 80.0
	case "MEDIUM":
		return 50.0
	case "LOW":
		return 30.0
	case "SAFE":
		return 0.0
	default:
		return 50.0
	}
}

func getNewContainers(result1, result2 *models.ValidationResult) []string {
	// Containers in result2 but not in result1
	// (simplified - in a real scenario would compare by name/version)
	return []string{}
}

func getRemovedContainers(result1, result2 *models.ValidationResult) []string {
	return []string{}
}

func getUpdatedContainers(result1, result2 *models.ValidationResult) map[string]interface{} {
	return make(map[string]interface{})
}

func getFixedIssues(result1, result2 *models.ValidationResult) []string {
	// Issues in result1 but not in result2
	fixed := make([]string, 0)
	for _, issue1 := range result1.Issues {
		found := false
		for _, issue2 := range result2.Issues {
			if issue1.Type == issue2.Type {
				found = true
				break
			}
		}
		if !found {
			fixed = append(fixed, issue1.Type)
		}
	}
	return fixed
}

func getNewIssues(result1, result2 *models.ValidationResult) []string {
	// Issues in result2 but not in result1
	newIssues := make([]string, 0)
	for _, issue2 := range result2.Issues {
		found := false
		for _, issue1 := range result1.Issues {
			if issue2.Type == issue1.Type {
				found = true
				break
			}
		}
		if !found {
			newIssues = append(newIssues, issue2.Type)
		}
	}
	return newIssues
}
