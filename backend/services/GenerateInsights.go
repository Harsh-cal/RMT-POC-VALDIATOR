package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/harsh-cal/rmt-poc-validator/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// GenerateInsight creates an AI-generated insight for validation results
// If OpenAI API fails or times out, returns fallback insight with nil error

func GenerateInsight(releaseInfo string, issues []models.Issue, risk string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return generateFallbackInsight(releaseInfo, risk, len(issues)), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	prompt := buildInsightPrompt(releaseInfo, issues, risk)

	response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		// 🔑 IMPORTANT: omit MaxTokens completely (avoids param issue)
	})

	if err != nil {
		return generateFallbackInsight(releaseInfo, risk, len(issues)), nil
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content != "" {
		return response.Choices[0].Message.Content, nil
	}

	return generateFallbackInsight(releaseInfo, risk, len(issues)), nil
}

// buildInsightPrompt creates the prompt for OpenAI
func buildInsightPrompt(releaseInfo string, issues []models.Issue, risk string) string {
	issuesSummary := ""
	for _, issue := range issues {
		issuesSummary += fmt.Sprintf("- %s (%s): %s\n", issue.Type, issue.Severity, issue.Message)
	}

	if issuesSummary == "" {
		issuesSummary = "No issues detected"
	}

	prompt := fmt.Sprintf(`You are an expert aviation validation engineer. 
Analyze this release and provide concise assessment (2-3 sentences).
Focus on critical issues without speculation.

Release: %s
Issues: %s
Risk: %s`, releaseInfo, issuesSummary, risk)

	return prompt
}

// generateFallbackInsight creates deterministic insight when API unavailable
func generateFallbackInsight(releaseInfo string, risk string, issueCount int) string {
	switch risk {
	case "HIGH":
		return fmt.Sprintf("Release contains critical issues (%d). All HIGH severity items must be resolved. Do not deploy.", issueCount)
	case "MEDIUM":
		return fmt.Sprintf("Release has medium risks (%d issues). Address MEDIUM severity items and conduct thorough testing.", issueCount)
	case "LOW":
		return fmt.Sprintf("Release has low risks (%d minor issues). May proceed with standard validation.", issueCount)
	case "SAFE":
		return "Release passed all checks. Safe for deployment."
	default:
		return "Validation complete. Review issues for deployment readiness."
	}
}

