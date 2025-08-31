package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	botToken   string
	chatID     int64
	httpClient *http.Client
	baseURL    string
}

type SendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}

func NewClient(botToken string, chatID int64) *Client {
	return &Client{
		botToken:   botToken,
		chatID:     chatID,
		baseURL:    "https://api.telegram.org/bot",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) SendMessage(message string) error {
	return c.SendMessageWithParseMode(message, "")
}

func (c *Client) SendMessageWithParseMode(message, parseMode string) error {
	cleanToken := strings.Trim(c.botToken, "\"")
	url := fmt.Sprintf("%s%s/sendMessage", c.baseURL, cleanToken)

	payload := SendMessageRequest{
		ChatID:    c.chatID,
		Text:      message,
		ParseMode: parseMode,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var telegramResp TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %s", telegramResp.Description)
	}

	return nil
}

func (c *Client) SendImageChangeNotification(imageURL string, timestamp time.Time) error {
	message := fmt.Sprintf(
		"🔄 *Image Change Detected*\n\n"+
			"📸 URL: %s\n"+
			"⏰ Time: %s\n"+
			"🤖 Watchdog is monitoring your image!",
		imageURL,
		timestamp.Format("2006-01-02 15:04:05 UTC"),
	)

	return c.SendMessageWithParseMode(message, "Markdown")
}

func (c *Client) TestConnection() error {
	cleanToken := strings.Trim(c.botToken, "\"")
	url := fmt.Sprintf("%s%s/getMe", c.baseURL, cleanToken)
	fmt.Printf("%s", url)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to test connection: %w", err)
	}
	defer resp.Body.Close()

	var telegramResp TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %s", telegramResp.Description)
	}

	return nil
}
