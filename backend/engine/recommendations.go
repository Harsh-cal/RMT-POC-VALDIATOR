package engine

import (
	"fmt"

	"github.com/harsh-cal/rmt-poc-validator/models"
)

// GenerateRecommendations creates fix actions for detected issues
func GenerateRecommendations(issues []models.Issue) []models.Recommendation {
	var recommendations []models.Recommendation
	seen := make(map[string]bool) // Avoid duplicate recommendations

	for _, issue := range issues {
		key := issue.Type + ":" + issue.Container
		if seen[key] {
			continue // Skip if we already have a rec for this issue type+container combo
		}

		rec := buildRecommendation(issue)
		if rec != nil {
			recommendations = append(recommendations, *rec)
			seen[key] = true
		}
	}

	return recommendations
}

// buildRecommendation creates a specific recommendation for an issue
func buildRecommendation(issue models.Issue) *models.Recommendation {
	switch issue.Type {
	case "version_mismatch":
		return &models.Recommendation{
			IssueType: "version_mismatch",
			Action:    fmt.Sprintf("Upgrade '%s' to satisfy the required version constraint", issue.Container),
		}

	case "missing_dependency":
		return &models.Recommendation{
			IssueType: "missing_dependency",
			Action:    fmt.Sprintf("Add the missing dependency container to the release or mark '%s' as optional if not needed", issue.Container),
		}

	case "duplicate":
		return &models.Recommendation{
			IssueType: "duplicate",
			Action:    fmt.Sprintf("Remove duplicate entries of '%s' from the release package", issue.Container),
		}

	case "unsupported_combo":
		return &models.Recommendation{
			IssueType: "unsupported_combo",
			Action:    fmt.Sprintf("Remove one of the incompatible containers or select compatible versions of '%s'", issue.Container),
		}

	default:
		return nil
	}
}
