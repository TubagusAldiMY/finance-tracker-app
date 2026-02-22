package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra"
	"github.com/spf13/viper"
)

type telegramUpdate struct {
	UpdateID int64            `json:"update_id"`
	Message  *telegramMessage `json:"message"`
}

type telegramMessage struct {
	Chat telegramChat `json:"chat"`
	Text string       `json:"text"`
}

type telegramChat struct {
	ID int64 `json:"id"`
}

type sendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type deleteWebhookRequest struct {
	DropPendingUpdates bool `json:"drop_pending_updates"`
}

type telegramBasicResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

type getUpdatesResponse struct {
	OK          bool             `json:"ok"`
	Result      []telegramUpdate `json:"result"`
	Description string           `json:"description"`
}

func main() {
	cfg := infra.NewViper("config.telegram")
	log := infra.NewLogger(cfg)

	token := resolveBotToken(cfg)
	if token == "" {
		log.Fatal("FINANCE_BOT_TOKEN is not configured")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpClient := &http.Client{Timeout: 70 * time.Second}

	dropPending := cfg.GetBool("telegram.drop_pending_updates")
	if err := deleteTelegramWebhook(ctx, httpClient, token, dropPending); err != nil {
		log.Fatalf("failed to switch bot to polling mode: %v", err)
	}

	pollTimeoutSeconds := cfg.GetInt("telegram.poll_timeout_seconds")
	if pollTimeoutSeconds <= 0 {
		pollTimeoutSeconds = 30
	}

	log.Info("Telegram bot polling started")
	if err := runPollingLoop(ctx, httpClient, token, pollTimeoutSeconds, log); err != nil {
		log.Fatalf("telegram polling stopped with error: %v", err)
	}

	log.Info("Telegram bot exited properly")
}

func newTelegramViper() *viper.Viper {
	v := viper.New()

	v.SetConfigName("config.telegram")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AddConfigPath("./backend")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			v.SetConfigName("config")
			if err2 := v.ReadInConfig(); err2 != nil {
				if _, ok2 := err2.(viper.ConfigFileNotFoundError); !ok2 {
					panic(fmt.Errorf("fatal error config file: %w", err2))
				}
			}
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return v
}

func runPollingLoop(ctx context.Context, client *http.Client, token string, pollTimeoutSeconds int, log interface {
	Info(...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}) error {
	var offset int64

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		updates, err := getTelegramUpdates(ctx, client, token, offset, pollTimeoutSeconds)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Errorf("failed getUpdates: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, upd := range updates {
			if upd.UpdateID >= offset {
				offset = upd.UpdateID + 1
			}

			if upd.Message == nil || upd.Message.Chat.ID == 0 {
				continue
			}

			text := strings.TrimSpace(upd.Message.Text)
			if text == "" {
				continue
			}

			reply := buildReply(text)
			if reply == "" {
				continue
			}

			if err := sendTelegramMessage(ctx, client, token, upd.Message.Chat.ID, reply); err != nil {
				log.Warnf("failed sendMessage to chat %d: %v", upd.Message.Chat.ID, err)
			}
		}
	}
}

func resolveBotToken(cfg interface{ GetString(string) string }) string {
	token := strings.TrimSpace(cfg.GetString("finance_bot_token"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("FINANCE_BOT_TOKEN"))
	}
	return token
}

func buildReply(text string) string {
	switch strings.ToLower(text) {
	case "/start":
		return "Halo, saya bot Finance Tracker. Gunakan /help untuk lihat perintah."
	case "/help":
		return "Perintah yang tersedia: /start, /help, /ping"
	case "/ping":
		return "pong"
	default:
		return "Perintah belum dikenali. Ketik /help"
	}
}

func getTelegramUpdates(ctx context.Context, client *http.Client, token string, offset int64, timeoutSeconds int) ([]telegramUpdate, error) {
	params := url.Values{}
	params.Set("offset", fmt.Sprintf("%d", offset))
	params.Set("timeout", fmt.Sprintf("%d", timeoutSeconds))

	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?%s", token, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build getUpdates request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request getUpdates: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("getUpdates status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp getUpdatesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("decode getUpdates response: %w", err)
	}
	if !apiResp.OK {
		return nil, fmt.Errorf("getUpdates rejected: %s", apiResp.Description)
	}

	return apiResp.Result, nil
}

func sendTelegramMessage(ctx context.Context, client *http.Client, token string, chatID int64, text string) error {
	payload := sendMessageRequest{ChatID: chatID, Text: text}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal sendMessage payload: %w", err)
	}

	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build sendMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request sendMessage: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("sendMessage status %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp telegramBasicResponse
	if err := json.Unmarshal(respBody, &apiResp); err == nil && !apiResp.OK {
		return fmt.Errorf("sendMessage rejected: %s", apiResp.Description)
	}

	return nil
}

func deleteTelegramWebhook(ctx context.Context, client *http.Client, token string, dropPending bool) error {
	payload := deleteWebhookRequest{DropPendingUpdates: dropPending}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal deleteWebhook payload: %w", err)
	}

	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/deleteWebhook", token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build deleteWebhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request deleteWebhook: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("deleteWebhook status %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp telegramBasicResponse
	if err := json.Unmarshal(respBody, &apiResp); err == nil && !apiResp.OK {
		return fmt.Errorf("deleteWebhook rejected: %s", apiResp.Description)
	}

	return nil
}
