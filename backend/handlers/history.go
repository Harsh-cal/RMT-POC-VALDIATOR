package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harsh-cal/rmt-poc-validator/services"
)

// HistoryHandler handles GET /releases/history requests
func HistoryHandler(c *gin.Context) {
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	history, err := services.GetReleaseHistory(limit, offset)
	if err != nil {
		fmt.Printf("HistoryHandler error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch history: %v", err)})
		return
	}

	c.JSON(http.StatusOK, history)
}

// TrendHandler handles GET /releases/trends requests
func TrendHandler(c *gin.Context) {
	days := 90

	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	trends, err := services.GetTrendData(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trends"})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// IssuesHandler handles GET /issues/recurring requests
func IssuesHandler(c *gin.Context) {
	days := 90

	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	issues, err := services.GetRecurringIssues(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recurring issues"})
		return
	}

	c.JSON(http.StatusOK, issues)
}

// CompareHandler handles GET /releases/:id1/compare/:id2 requests
func CompareHandler(c *gin.Context) {
	id1 := c.Param("id1")
	id2 := c.Param("id2")

	if id1 == "" || id2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both release IDs are required"})
		return
	}

	comparison, err := services.CompareReleases(id1, id2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comparison)
}
