package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/harsh-cal/rmt-poc-validator/engine"
	"github.com/harsh-cal/rmt-poc-validator/handlers"
	"github.com/harsh-cal/rmt-poc-validator/models"
	"github.com/harsh-cal/rmt-poc-validator/services"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	godotenv.Load()

	if err := services.InitMongo(); err != nil {
		fmt.Printf("Warning: MongoDB init failed: %v\n", err)
	}
	defer services.CloseMongo()

	// CLI flags
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	filePath := validateCmd.String("file", "", "Path to release JSON file")
	releaseName := validateCmd.String("release", "", "Name of mock release")
	output := validateCmd.String("output", "json", "Output format: json or text")

	if len(os.Args) < 2 {
		startServer()
		return
	}

	if os.Args[1] == "validate" {
		validateCmd.Parse(os.Args[2:])

		if *filePath != "" {
			validateFromFile(*filePath, *output)
		} else if *releaseName != "" {
			validateFromMockRelease(*releaseName, *output)
		} else {
			fmt.Println("Usage: rmt-validator validate --file <path> OR --release <name>")
		}
	} else {
		startServer()
	}
}

func startServer() {
	router := gin.Default()

	router.Use(cors.Default())

	// Routes
	router.POST("/api/dev/v1/validate", handlers.ValidateHandler)
	router.POST("/api/dev/v1/validate/chat", handlers.ChatHandler)
	router.POST("/api/dev/v1/validate/export", handlers.ExportValidationHandler)

	// History & Analytics Routes
	router.GET("/api/dev/v1/releases/history", handlers.HistoryHandler)
	router.GET("/api/dev/v1/releases/trends", handlers.TrendHandler)
	router.GET("/api/dev/v1/issues/recurring", handlers.IssuesHandler)
	router.GET("/api/dev/v1/releases/:id1/compare/:id2", handlers.CompareHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	router.Run(":" + port)
}

func validateFromFile(filePath string, outputFormat string) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Parse JSON
	var req models.ValidateRequest
	if err := json.Unmarshal(data, &req); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	result := performValidation(&req)

	// Save to MongoDB (synchronous for CLI)
	services.SaveRelease(&req, result.ReleaseID)
	services.SaveValidationResult(result)

	// Output
	outputResult(result, outputFormat)
}

func validateFromMockRelease(releaseName string, outputFormat string) {
	// Load mock release from JSON folder
	mockReleases := getMockReleases()

	// Normalize: trim .json extension if accidentally passed
	releaseName = strings.TrimSuffix(releaseName, ".json")

	var mockReq *models.ValidateRequest
	for _, mock := range mockReleases {
		if strings.EqualFold(mock.ReleaseName, releaseName) {
			mockReq = &models.ValidateRequest{
				ReleaseName: mock.ReleaseName,
				Version:     mock.Version,
				TargetFleet: mock.TargetFleet,
				Aircraft:    mock.Aircraft,
				Containers:  mock.Containers,
			}
			break
		}
	}

	if mockReq == nil {
		fmt.Printf("Mock release '%s' not found\n", releaseName)
		return
	}

	// Validate
	result := performValidation(mockReq)

	services.SaveRelease(mockReq, result.ReleaseID)
	services.SaveValidationResult(result)

	// Output
	outputResult(result, outputFormat)
}

func performValidation(req *models.ValidateRequest) *models.ValidationResult {
	ctx := &models.ValidationContext{
		Request:   req,
		ReleaseID: uuid.New().String(),
	}

	// Run validation pipeline
	ctx.Issues = engine.DetectIssues(ctx)
	ctx.Risk = engine.CalculateRisk(ctx.Issues)
	ctx.Status = engine.CalculateStatus(ctx.Issues)
	recommendations := engine.GenerateRecommendations(ctx.Issues)

	releaseInfo := req.ReleaseName + " v" + req.Version + " for " + req.TargetFleet
	insight, _ := services.GenerateInsight(releaseInfo, ctx.Issues, ctx.Risk)

	return &models.ValidationResult{
		ReleaseID:       ctx.ReleaseID,
		ReleaseName:     req.ReleaseName,
		Version:         req.Version,
		TargetFleet:     req.TargetFleet,
		Risk:            ctx.Risk,
		Status:          ctx.Status,
		Issues:          ctx.Issues,
		Insight:         insight,
		Recommendations: recommendations,
		ValidatedAt:     time.Now(),
	}
}

func outputResult(result *models.ValidationResult, format string) {
	if format == "json" {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else if format == "text" {
		fmt.Printf("Release: %s v%s\n", result.ReleaseName, result.Version)
		fmt.Printf("Status: %s\n", result.Status)
		fmt.Printf("Risk: %s\n", result.Risk)
		fmt.Printf("Issues: %d\n", len(result.Issues))
		fmt.Printf("Insight Summary: %s\n", result.Insight.Summary)
		fmt.Printf("Insight Impact: %s\n", result.Insight.Impact)
	}
}

func getMockReleases() []models.MockRelease {
	mockDir := "mock"
	var releases []models.MockRelease

	entries, err := os.ReadDir(mockDir)
	if err != nil {
		fmt.Printf("Warning: could not read mock directory: %v\n", err)
		return releases
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := mockDir + "/" + entry.Name()
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: could not read mock file %s: %v\n", filePath, err)
			continue
		}

		var mock models.MockRelease
		if err := json.Unmarshal(data, &mock); err != nil {
			fmt.Printf("Warning: could not parse mock file %s: %v\n", filePath, err)
			continue
		}

		releases = append(releases, mock)
	}

	return releases
}
