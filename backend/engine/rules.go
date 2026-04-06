package engine

import (
	"fmt"
	"strings"

	"github.com/harsh-cal/rmt-poc-validator/models"
)

// DetectIssues runs all validation rules and returns issues found
func DetectIssues(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue

	// Check for duplicates first
	duplicateIssues := checkDuplicates(ctx)
	issues = append(issues, duplicateIssues...)

	// Check for missing dependencies
	missingDepIssues := checkMissingDependencies(ctx)
	issues = append(issues, missingDepIssues...)

	// Check for version mismatches
	versionIssues := checkVersionMismatches(ctx)
	issues = append(issues, versionIssues...)

	// Check for unsupported combinations
	comboIssues := checkUnsupportedCombos(ctx)
	issues = append(issues, comboIssues...)

	// Check for duplicate part numbers
	partNumberIssues := checkDuplicatePartNumbers(ctx)
	issues = append(issues, partNumberIssues...)

	// Check system compatibility against aircraft
	compatibilityIssues := checkSystemCompatibility(ctx)
	issues = append(issues, compatibilityIssues...)

	// Check version conflicts against currently installed software
	stateConflictIssues := checkAircraftStateVersionConflicts(ctx)
	issues = append(issues, stateConflictIssues...)

	// Add maturity-based risk signals
	maturityIssues := checkMaturityRisk(ctx)
	issues = append(issues, maturityIssues...)

	return issues
}

// checkDuplicates detects if same container appears more than once
func checkDuplicates(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue
	seen := make(map[string]int)

	for _, container := range ctx.Request.Containers {
		seen[container.Name]++
	}

	for name, count := range seen {
		if count > 1 {
			issues = append(issues, models.Issue{
				Type:      "duplicate",
				Severity:  "MEDIUM",
				Container: name,
				Message:   fmt.Sprintf("Container '%s' appears %d times in release", name, count),
			})
		}
	}

	return issues
}

// checkMissingDependencies detects if required dependency is not in release
func checkMissingDependencies(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue

	// Build map of available containers
	containerMap := make(map[string]bool)
	for _, container := range ctx.Request.Containers {
		containerMap[container.Name] = true
	}

	// Check each container's dependencies
	for _, container := range ctx.Request.Containers {
		for _, dep := range container.Dependencies {
			if !containerMap[dep.Name] {
				issues = append(issues, models.Issue{
					Type:      "missing_dependency",
					Severity:  "HIGH",
					Container: container.Name,
					Message:   fmt.Sprintf("Required dependency '%s' not found in release", dep.Name),
				})
			}
		}
	}

	return issues
}

// checkVersionMismatches detects version constraint violations
func checkVersionMismatches(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue

	// Build map of container versions
	containerVersions := make(map[string]string)
	for _, container := range ctx.Request.Containers {
		containerVersions[container.Name] = container.Version
	}

	// Check each container's dependencies
	for _, container := range ctx.Request.Containers {
		for _, dep := range container.Dependencies {
			if availableVersion, exists := containerVersions[dep.Name]; exists {
				if !satisfiesVersionConstraint(availableVersion, dep.RequiredVersion) {
					issues = append(issues, models.Issue{
						Type:      "version_mismatch",
						Severity:  "HIGH",
						Container: dep.Name,
						Message:   fmt.Sprintf("'%s' requires '%s' %s, found %s", container.Name, dep.Name, dep.RequiredVersion, availableVersion),
					})
				}
			}
		}
	}

	return issues
}

// checkUnsupportedCombos detects known incompatible container pairs
func checkUnsupportedCombos(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue

	// Build map of containers in release
	containerNames := make(map[string]bool)
	for _, container := range ctx.Request.Containers {
		containerNames[container.Name] = true
	}

	// Get unsupported pairs from config
	unsupportedPairs := models.GetUnsupportedCombos()

	// Check each pair
	for _, pair := range unsupportedPairs {
		if containerNames[pair.Container1] && containerNames[pair.Container2] {
			issues = append(issues, models.Issue{
				Type:      "unsupported_combo",
				Severity:  "HIGH",
				Container: pair.Container1,
				Message:   fmt.Sprintf("'%s' and '%s' are incompatible and cannot be deployed together", pair.Container1, pair.Container2),
			})
			// Also add for second container for visibility
			issues = append(issues, models.Issue{
				Type:      "unsupported_combo",
				Severity:  "HIGH",
				Container: pair.Container2,
				Message:   fmt.Sprintf("'%s' and '%s' are incompatible and cannot be deployed together", pair.Container1, pair.Container2),
			})
		}
	}

	return issues
}

// checkDuplicatePartNumbers detects duplicate part numbers in release payload
func checkDuplicatePartNumbers(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue
	seen := make(map[string]int)

	for _, container := range ctx.Request.Containers {
		if strings.TrimSpace(container.PartNumber) == "" {
			continue
		}
		seen[container.PartNumber]++
	}

	for partNumber, count := range seen {
		if count > 1 {
			issues = append(issues, models.Issue{
				Type:      "duplicate_part_number",
				Severity:  "HIGH",
				Container: partNumber,
				Message:   fmt.Sprintf("Part number '%s' appears %d times in release", partNumber, count),
			})
		}
	}

	return issues
}

// checkSystemCompatibility validates whether container system matches aircraft supported system
func checkSystemCompatibility(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue
	aircraftSystem := strings.TrimSpace(strings.ToLower(ctx.Request.Aircraft.System))
	if aircraftSystem == "" {
		return issues
	}

	for _, container := range ctx.Request.Containers {
		systemType := strings.TrimSpace(strings.ToLower(container.SystemType))
		if systemType == "" {
			continue
		}

		if systemType != aircraftSystem {
			issues = append(issues, models.Issue{
				Type:      "system_incompatible",
				Severity:  "HIGH",
				Container: container.Name,
				Message:   fmt.Sprintf("Container '%s' system '%s' is incompatible with aircraft system '%s'", container.Name, container.SystemType, ctx.Request.Aircraft.System),
			})
		}
	}

	return issues
}

// checkAircraftStateVersionConflicts catches downgrade or conflict vs installed software on aircraft
func checkAircraftStateVersionConflicts(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue
	if len(ctx.Request.Aircraft.CurrentSoftware) == 0 {
		return issues
	}

	for _, container := range ctx.Request.Containers {
		installed, exists := ctx.Request.Aircraft.CurrentSoftware[container.Name]
		if !exists || strings.TrimSpace(installed.Version) == "" {
			continue
		}

		if compareVersions(container.Version, installed.Version) < 0 {
			issues = append(issues, models.Issue{
				Type:      "version_conflict",
				Severity:  "HIGH",
				Container: container.Name,
				Message:   fmt.Sprintf("Release version %s for '%s' is older than installed version %s on aircraft", container.Version, container.Name, installed.Version),
			})
		}
	}

	return issues
}

// checkMaturityRisk adds risk signals based on container maturity level
func checkMaturityRisk(ctx *models.ValidationContext) []models.Issue {
	var issues []models.Issue

	for _, container := range ctx.Request.Containers {
		switch strings.ToLower(strings.TrimSpace(container.Maturity)) {
		case "experimental":
			issues = append(issues, models.Issue{
				Type:      "maturity_risk",
				Severity:  "HIGH",
				Container: container.Name,
				Message:   fmt.Sprintf("Container '%s' is marked experimental and is not recommended for deployment", container.Name),
			})
		case "beta":
			issues = append(issues, models.Issue{
				Type:      "maturity_risk",
				Severity:  "MEDIUM",
				Container: container.Name,
				Message:   fmt.Sprintf("Container '%s' is marked beta; additional validation is recommended", container.Name),
			})
		case "candidate", "preview", "rc":
			issues = append(issues, models.Issue{
				Type:      "maturity_risk",
				Severity:  "LOW",
				Container: container.Name,
				Message:   fmt.Sprintf("Container '%s' is pre-release (%s); monitor closely after deployment", container.Name, container.Maturity),
			})
		}
	}

	return issues
}

// satisfiesVersionConstraint checks if availableVersion meets the constraint
// For simplicity, support >= operator (e.g., ">=4.0.0")
func satisfiesVersionConstraint(availableVersion, constraint string) bool {
	// Parse constraint: expected format ">=X.Y.Z"
	if strings.HasPrefix(constraint, ">=") {
		requiredVersion := strings.TrimPrefix(constraint, ">=")
		return compareVersions(availableVersion, requiredVersion) >= 0
	}

	// If no operator specified, treat as exact match or >= for now
	return compareVersions(availableVersion, constraint) >= 0
}

// compareVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
// Simple comparison: split by . and compare as integers
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Compare each part
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		p1 := "0"
		p2 := "0"

		if i < len(parts1) {
			p1 = parts1[i]
		}
		if i < len(parts2) {
			p2 = parts2[i]
		}

		// Extract numeric part (in case of prerelease tags)
		p1Num := extractNumeric(p1)
		p2Num := extractNumeric(p2)

		if p1Num < p2Num {
			return -1
		} else if p1Num > p2Num {
			return 1
		}
	}

	return 0
}

// extractNumeric extracts numeric prefix from a version part
func extractNumeric(s string) int {
	var numStr string
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			numStr += string(ch)
		} else {
			break
		}
	}

	if numStr == "" {
		return 0
	}

	var result int
	fmt.Sscanf(numStr, "%d", &result)
	return result
}
