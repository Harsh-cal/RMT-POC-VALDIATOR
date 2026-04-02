package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/harsh-cal/rmt-poc-validator/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// GenerateInsight creates an AI-generated insight for validation results
// If OpenAI API fails or times out, returns fallback insight with nil error

func GenerateInsight(releaseInfo string, issues []models.Issue, risk string) (models.Insight, error) {
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
		return models.Insight{
			Summary: response.Choices[0].Message.Content,
			Impact:  buildImpactFromRisk(risk),
		}, nil
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
func generateFallbackInsight(releaseInfo string, risk string, issueCount int) models.Insight {
	switch risk {
	case "HIGH":
		return models.Insight{
			Summary: fmt.Sprintf("Release contains critical issues (%d). All HIGH severity items must be resolved.", issueCount),
			Impact:  "GO/NO-GO decision: NO-GO. Deployment may cause system instability and operational risk.",
		}
	case "MEDIUM":
		return models.Insight{
			Summary: fmt.Sprintf("Release has medium risks (%d issues). Address MEDIUM severity items and conduct thorough testing.", issueCount),
			Impact:  "Deployment should be gated until medium-risk findings are remediated.",
		}
	case "LOW":
		return models.Insight{
			Summary: fmt.Sprintf("Release has low risks (%d minor issues). May proceed with standard validation.", issueCount),
			Impact:  "Low operational impact expected with standard verification checks.",
		}
	case "SAFE":
		return models.Insight{
			Summary: "Release passed all checks. Safe for deployment.",
			Impact:  "No blocking impact identified for target aircraft deployment.",
		}
	default:
		return models.Insight{
			Summary: "Validation complete. Review issues for deployment readiness.",
			Impact:  "Final deployment decision should be based on issue severity and compatibility checks.",
		}
	}
}

func buildImpactFromRisk(risk string) string {
	switch risk {
	case "HIGH":
		return "GO/NO-GO decision: NO-GO. Critical risks can impact aircraft operations."
	case "MEDIUM":
		return "GO/NO-GO decision: CONDITIONAL NO-GO until medium-risk issues are mitigated."
	case "LOW":
		return "GO/NO-GO decision: GO with caution and standard post-deployment monitoring."
	default:
		return "GO/NO-GO decision: GO. No critical impact identified."
	}
}

// GenerateValidationChatAnswer returns a grounded answer about one validation result.
func GenerateValidationChatAnswer(question string, result models.ValidationResult, history []models.ChatMessage) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "OpenAI key is not configured. Please set OPENAI_API_KEY on the backend.", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := openai.NewClient(option.WithAPIKey(apiKey))

	systemPrompt := buildValidationChatSystemPrompt(result)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	for _, msg := range history {
		role := strings.ToLower(strings.TrimSpace(msg.Role))
		if role == "assistant" {
			messages = append(messages, openai.AssistantMessage(msg.Content))
			continue
		}
		messages = append(messages, openai.UserMessage(msg.Content))
	}

	messages = append(messages, openai.UserMessage(question))

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    "gpt-4o-mini",
		Messages: messages,
	})
	if err != nil {
		return "I could not generate a response right now. Please try again.", nil
	}

	if len(resp.Choices) == 0 || strings.TrimSpace(resp.Choices[0].Message.Content) == "" {
		return "I do not have enough context to answer that from this validation result.", nil
	}

	return resp.Choices[0].Message.Content, nil
}

func buildValidationChatSystemPrompt(result models.ValidationResult) string {
	issues := "None"
	if len(result.Issues) > 0 {
		var lines []string
		for _, i := range result.Issues {
			lines = append(lines, fmt.Sprintf("- [%s] %s: %s", i.Severity, i.Type, i.Message))
		}
		issues = strings.Join(lines, "\n")
	}

	recs := "None"
	if len(result.Recommendations) > 0 {
		var lines []string
		for _, r := range result.Recommendations {
			lines = append(lines, fmt.Sprintf("- %s: %s", r.IssueType, r.Action))
		}
		recs = strings.Join(lines, "\n")
	}

	return fmt.Sprintf(`You are an RMT (Release Management Tool) validation assistant for aircraft software deployment decisions.

You can answer ONLY from the validation context below.
If asked anything outside this context, say you do not have that data.
Be concise and actionable for airline engineers.

Validation Context:
Release: %s v%s
Target Fleet: %s
Status: %s
Risk: %s
Issues (%d):
%s
Recommendations:
%s
Insight Summary: %s
Insight Impact: %s`,
		result.ReleaseName,
		result.Version,
		result.TargetFleet,
		result.Status,
		result.Risk,
		len(result.Issues),
		issues,
		recs,
		result.Insight.Summary,
		result.Insight.Impact,
	)
}
