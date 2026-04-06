package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harsh-cal/rmt-poc-validator/engine"
	"github.com/harsh-cal/rmt-poc-validator/models"
	"github.com/harsh-cal/rmt-poc-validator/services"
)

// ValidateHandler handles POST /validate requests
func ValidateHandler(c *gin.Context) {
	// Parse request
	var req models.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Generate release ID
	releaseID := uuid.New().String()

	// Create validation context
	ctx := &models.ValidationContext{
		Request:   &req,
		ReleaseID: releaseID,
	}

	// Run validation pipeline

	// 1. Detect issues
	ctx.Issues = engine.DetectIssues(ctx)

	// 2. Calculate risk
	ctx.Risk = engine.CalculateRisk(ctx.Issues)
	ctx.Status = engine.CalculateStatus(ctx.Issues)

	// 3. Generate recommendations
	recommendations := engine.GenerateRecommendations(ctx.Issues)

	// 4. Get AI insight
	releaseInfo := req.ReleaseName + " v" + req.Version + " for " + req.TargetFleet
	insight, _ := services.GenerateInsight(releaseInfo, ctx.Issues, ctx.Risk)

	// 5. Build response
	result := models.ValidationResult{
		ReleaseID:       releaseID,
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

	if err := services.SaveRelease(&req, releaseID); err != nil {
		fmt.Printf("SaveRelease failed for release_id=%s: %v\n", releaseID, err)
	}
	if err := services.SaveValidationResult(&result); err != nil {
		fmt.Printf("SaveValidationResult failed for release_id=%s: %v\n", releaseID, err)
	}

	// Send Telegram alert asynchronously using an immutable snapshot
	resultSnapshot := result
	releaseNameSnapshot := req.ReleaseName
	go func(snapshot models.ValidationResult, releaseName string, rid string) {
		if err := services.SendValidationAlert(&snapshot, releaseName); err != nil {
			fmt.Printf("Telegram alert failed for release_id=%s: %v\n", rid, err)
		}
	}(resultSnapshot, releaseNameSnapshot, releaseID)

	// Return response
	c.JSON(http.StatusOK, result)
}

// ChatHandler handles POST /validate/chat requests.
func ChatHandler(c *gin.Context) {
	var req models.ValidationChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if strings.TrimSpace(req.Question) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question is required"})
		return
	}

	answer, err := services.GenerateValidationChatAnswer(req.Question, req.Result, req.History)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate chat answer"})
		return
	}

	c.JSON(http.StatusOK, models.ValidationChatResponse{Answer: answer})
}
