package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/harsh-cal/rmt-poc-validator/models"
)

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// SendValidationAlert sends a Telegram notification for validation result
func SendValidationAlert(result *models.ValidationResult, releaseName string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		return nil
	}

	statusIcon := "✅"
	if result.Status != "PASS" {
		statusIcon = "❌"
	}

	riskColor := "🟢"
	switch result.Risk {
	case "HIGH":
		riskColor = "🔴"
	case "MEDIUM":
		riskColor = "🟡"
	case "LOW":
		riskColor = "🟠"
	}

	var messageText strings.Builder
	messageText.WriteString(fmt.Sprintf("<b>%s Validation Result</b>\n\n", statusIcon))
	messageText.WriteString(fmt.Sprintf("<b>Release:</b> %s\n", releaseName))
	messageText.WriteString(fmt.Sprintf("<b>Version:</b> %s\n", result.Version))
	messageText.WriteString(fmt.Sprintf("<b>ID:</b> <code>%s</code>\n", result.ReleaseID))
	messageText.WriteString(fmt.Sprintf("<b>Status:</b> %s %s\n", statusIcon, result.Status))
	messageText.WriteString(fmt.Sprintf("<b>Risk Level:</b> %s %s\n", riskColor, result.Risk))

	// Add Fleet & Targeting context
	messageText.WriteString("\n<b>📍 Target Information:</b>\n")
	if result.TargetFleet != "" {
		messageText.WriteString(fmt.Sprintf("• <b>Fleet:</b> %s\n", result.TargetFleet))
	}
	if result.TailNumber != "" {
		messageText.WriteString(fmt.Sprintf("• <b>Tail Number:</b> %s\n", result.TailNumber))
	}
	if result.AircraftType != "" {
		messageText.WriteString(fmt.Sprintf("• <b>Aircraft Type:</b> %s\n", result.AircraftType))
	}
	if result.AircraftSystem != "" {
		messageText.WriteString(fmt.Sprintf("• <b>System:</b> %s\n", result.AircraftSystem))
	}

	messageText.WriteString(fmt.Sprintf("\n<b>Issues Found:</b> %d\n", len(result.Issues)))

	if len(result.Issues) > 0 {
		messageText.WriteString("\n<b>All Issues:</b>\n")
		for _, issue := range result.Issues {
			messageText.WriteString(fmt.Sprintf("• <b>%s</b> (%s) - %s\n", issue.Type, issue.Severity, issue.Container))
		}
	}

	msg := TelegramMessage{
		ChatID:    chatID,
		Text:      messageText.String(),
		ParseMode: "HTML",
	}

	return sendTelegramMessage(botToken, msg)
}

// sendTelegramMessage sends a message via Telegram Bot API
func sendTelegramMessage(botToken string, msg TelegramMessage) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal error: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("post error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.ReadAll(resp.Body)
		return fmt.Errorf("telegram api error: status %d", resp.StatusCode)
	}

	return nil
}
