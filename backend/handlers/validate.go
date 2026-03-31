package handlers

import (
    "net/http"
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
        Issues:          ctx.Issues,
        Insight:         insight,
        Recommendations: recommendations,
        ValidatedAt:     time.Now(),
    }

  go func() {
    services.SaveRelease(&req, releaseID)
    services.SaveValidationResult(&result)
}()


    // Return response
    c.JSON(http.StatusOK, result)
}