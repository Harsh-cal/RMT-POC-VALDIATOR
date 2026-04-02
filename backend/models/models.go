package models

import "time"

// --- Request/Response Models ---

// ValidateRequest is the incoming request payload for the validate endpoint
type ValidateRequest struct {
	ReleaseName string      `json:"release_name" binding:"required"`
	Version     string      `json:"version" binding:"required"`
	TargetFleet string      `json:"target_fleet" binding:"required"`
	Aircraft    Aircraft    `json:"aircraft" binding:"required"`
	Containers  []Container `json:"containers" binding:"required"`
}

// Container represents a software artifact in a release
type Container struct {
	Name         string       `json:"name" binding:"required"`
	Version      string       `json:"version" binding:"required"`
	PartNumber   string       `json:"partNumber" binding:"required"`
	SystemType   string       `json:"systemType" binding:"required"`
	Maturity     string       `json:"maturity"` // experimental, beta, stable
	Dependencies []Dependency `json:"dependencies"`
	IsOptional   bool         `json:"is_optional"`
}

// Aircraft represents the target aircraft state for validation
type Aircraft struct {
	TailNumber      string                    `json:"tailNumber" binding:"required"`
	Type            string                    `json:"type" binding:"required"`
	System          string                    `json:"system" binding:"required"`
	CurrentSoftware map[string]InstalledImage `json:"currentSoftware"`
}

// InstalledImage represents currently installed container metadata on aircraft
type InstalledImage struct {
	Version    string `json:"version"`
	PartNumber string `json:"partNumber"`
}

// Dependency represents a container dependency constraint
type Dependency struct {
	Name            string `json:"name" binding:"required"`
	RequiredVersion string `json:"required_version" binding:"required"`
}

// ValidationResult is the output of a validation run
type ValidationResult struct {
	ReleaseID       string           `json:"release_id" bson:"release_id"`
	ReleaseName     string           `json:"release_name" bson:"release_name"`
	Version         string           `json:"version" bson:"version"`
	TargetFleet     string           `json:"target_fleet" bson:"target_fleet"`
	Risk            string           `json:"risk" bson:"risk"`     // HIGH, MEDIUM, LOW, SAFE
	Status          string           `json:"status" bson:"status"` // PASS / FAILED
	Issues          []Issue          `json:"issues" bson:"issues"`
	Insight         Insight          `json:"insight" bson:"insight"`
	Recommendations []Recommendation `json:"recommendations" bson:"recommendations"`
	ValidatedAt     time.Time        `json:"validated_at" bson:"validated_at"`
}

// Insight represents structured AI output.
type Insight struct {
	Summary string `json:"summary" bson:"summary"`
	Impact  string `json:"impact" bson:"impact"`
}

// ChatMessage represents one chat turn.
type ChatMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// ValidationChatRequest is the payload for the validation chat endpoint.
type ValidationChatRequest struct {
	Question string           `json:"question" binding:"required"`
	Result   ValidationResult `json:"result" binding:"required"`
	History  []ChatMessage    `json:"history"`
}

// ValidationChatResponse is the response for validation chat endpoint.
type ValidationChatResponse struct {
	Answer string `json:"answer"`
}

// ExportRequest is the payload for exporting validation reports.
type ExportRequest struct {
	ReleaseID   string `json:"release_id"`
	ReleaseName string `json:"release_name"`
	Format      string `json:"format" binding:"required"`
}

// Issue represents a detected validation issue
type Issue struct {
	Type      string `json:"type" bson:"type"`         // version_mismatch, missing_dependency, duplicate, unsupported_combo
	Severity  string `json:"severity" bson:"severity"` // HIGH, MEDIUM, LOW
	Container string `json:"container" bson:"container"`
	Message   string `json:"message" bson:"message"`
}

// Recommendation represents a fix action for an issue
type Recommendation struct {
	IssueType string `json:"issue_type" bson:"issue_type"`
	Action    string `json:"action" bson:"action"`
}

// --- Internal Engine Models ---

// ValidationContext holds all data during validation pipeline
type ValidationContext struct {
	Request   *ValidateRequest
	Issues    []Issue
	Risk      string
	Status    string
	Insight   Insight
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
	Aircraft    Aircraft    `json:"aircraft"`
	Containers  []Container `json:"containers"`
}
