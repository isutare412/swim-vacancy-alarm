package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiSendMessage = "/sendMessage"
)

type Client struct {
	botToken string
	chatID   string
}

func NewClient(cfg Config) *Client {
	return &Client{
		botToken: cfg.BotToken,
		chatID:   cfg.ChatID,
	}
}

func (c *Client) SendMessage(ctx context.Context, message string) error {
	if len(message) == 0 {
		return nil
	}
	message = truncateLargeMessage(message)

	req := sendMessageRequest{
		ChatID: c.chatID,
		Text:   message,
	}

	reqBytes, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("marshaling message body: %w", err)
	}

	url := telegramBotAPIBase(c.botToken) + apiSendMessage
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return fmt.Errorf("building http request: %w", err)
	}

	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("doing http request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code %s; body (%s)", resp.Status, string(bodyBytes))
	}

	return nil
}

func truncateLargeMessage(msg string) string {
	if len(msg) > 4096 {
		return msg[:4096]
	}
	return msg
}

func telegramBotAPIBase(botToken string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", botToken)
}
