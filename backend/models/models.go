package models

import "time"

// --- Request/Response Models ---

// ValidateRequest is the incoming request payload for the validate endpoint
type ValidateRequest struct {
	ReleaseName string      `json:"release_name" binding:"required"`
	Version     string      `json:"version" binding:"required"`
	TargetFleet string      `json:"target_fleet" binding:"required"`
	Containers  []Container `json:"containers" binding:"required"`
}

// Container represents a software artifact in a release
type Container struct {
	Name         string       `json:"name" binding:"required"`
	Version      string       `json:"version" binding:"required"`
	Dependencies []Dependency `json:"dependencies"`
	IsOptional   bool         `json:"is_optional"`
}

// Dependency represents a container dependency constraint
type Dependency struct {
	Name             string `json:"name" binding:"required"`
	RequiredVersion  string `json:"required_version" binding:"required"`
}

// ValidationResult is the output of a validation run
type ValidationResult struct {
	ReleaseID   string           `json:"release_id"`
	ReleaseName string           `json:"release_name"`
	Version     string           `json:"version"`
	TargetFleet string           `json:"target_fleet"`
	Risk        string           `json:"risk"` // HIGH, MEDIUM, LOW, SAFE
	Issues      []Issue          `json:"issues"`
	Insight     string           `json:"insight"`
	Recommendations []Recommendation `json:"recommendations"`
	ValidatedAt time.Time        `json:"validated_at"`
}

// Issue represents a detected validation issue
type Issue struct {
	Type      string `json:"type"` // version_mismatch, missing_dependency, duplicate, unsupported_combo
	Severity  string `json:"severity"` // HIGH, MEDIUM, LOW
	Container string `json:"container"`
	Message   string `json:"message"`
}

// Recommendation represents a fix action for an issue
type Recommendation struct {
	IssueType string `json:"issue_type"`
	Action    string `json:"action"`
}

// --- Internal Engine Models ---

// ValidationContext holds all data during validation pipeline
type ValidationContext struct {
	Request   *ValidateRequest
	Issues    []Issue
	Risk      string
	Insight   string
	ReleaseID string
}

// --- Unsupported Combo Pair ---
type UnsupportedPair struct {
	Container1 string
	Container2 string
}

// GetUnsupportedCombos returns hardcoded list of incompatible container pairs
func GetUnsupportedCombos() []UnsupportedPair {
	return []UnsupportedPair{
		{Container1: "IFE_Software", Container2: "Legacy_Display_Driver"},
		{Container1: "Navigation_Module", Container2: "Old_GPS_Firmware"},
		{Container1: "Flight_Control_UI", Container2: "Deprecated_CanBus_Adapter"},
		{Container1: "Cabin_Connectivity_Service", Container2: "Satcom_V1_Modem"},
		{Container1: "Security_Agent_V3", Container2: "Telemetry_Collector_V1"},
		{Container1: "Engine_Monitoring_Pack", Container2: "Sensor_Bridge_Lite"},
		{Container1: "Crew_Apps_Runtime", Container2: "Java8_Base_Image"},
		{Container1: "EFB_Sync_Service", Container2: "Offline_Config_Loader"},
		{Container1: "Diagnostics_Core_V2", Container2: "Log_Forwarder_V0"},
		{Container1: "Map_Render_Engine_HD", Container2: "Memory_Optimizer_Legacy"},
	}
}

// --- Mock Release Models ---

// MockRelease represents a predefined release for demo/testing
type MockRelease struct {
	ReleaseName string      `json:"release_name"`
	Version     string      `json:"version"`
	TargetFleet string      `json:"target_fleet"`
	Containers  []Container `json:"containers"`
}
