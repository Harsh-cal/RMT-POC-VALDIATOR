package engine

import (
	"github.com/harsh-cal/rmt-poc-validator/models"
)

// CalculateRisk determines overall risk level based on issues found
// Logic: highest severity issue determines the risk level
func CalculateRisk(issues []models.Issue) string {
	if len(issues) == 0 {
		return "SAFE"
	}

	// Severity priority: HIGH > MEDIUM > LOW
	highestSeverity := "LOW"

	for _, issue := range issues {
		if issue.Severity == "HIGH" {
			highestSeverity = "HIGH"
			break // No need to check further, HIGH is highest
		}
		if issue.Severity == "MEDIUM" && highestSeverity != "HIGH" {
			highestSeverity = "MEDIUM"
		}
	}

	// Map severity to risk
	switch highestSeverity {
	case "HIGH":
		return "HIGH"
	case "MEDIUM":
		return "MEDIUM"
	case "LOW":
		return "LOW"
	default:
		return "SAFE"
	}
}

// CalculateStatus returns final release approval status based on GO/NO-GO checks.
// FAILED if any HIGH severity issue exists; otherwise PASS.
func CalculateStatus(issues []models.Issue) string {
	for _, issue := range issues {
		if issue.Severity == "HIGH" {
			return "FAILED"
		}
	}
	return "PASS"
}
