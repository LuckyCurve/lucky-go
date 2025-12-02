package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const telegramAPIBaseURL = "https://api.telegram.org/bot"

// TelegramMessage 表示发送到 Telegram 的消息
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// TelegramResponse 表示 Telegram API 的响应
type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}

// SendTelegramMessage 发送消息到 Telegram
func SendTelegramMessage(message string) error {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN 环境变量未设置")
	}

	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if chatID == "" {
		return fmt.Errorf("TELEGRAM_CHAT_ID 环境变量未设置")
	}

	// 构建请求体
	msg := TelegramMessage{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	// 发送请求
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIBaseURL, botToken)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var telegramResp TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if !telegramResp.OK {
		return fmt.Errorf("Telegram API 返回错误: %s", telegramResp.Description)
	}

	return nil
}
