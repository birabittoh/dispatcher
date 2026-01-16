package config

import (
	"errors"
	"log/slog"
	"os"
)

type Config struct {
	ListenAddress     string
	GitLabSecretToken string
	LogAPIKey         string
	TelegramBotToken  string
	TelegramChatID    string
	TelegramThreadID  string
}

var (
	ErrMissingTelegramBotToken = errors.New("TELEGRAM_BOT_TOKEN is not set")
	ErrMissingTelegramChatID   = errors.New("TELEGRAM_CHAT_ID is not set")
)

func LoadConfig() (config Config, err error) {
	config = Config{
		ListenAddress:     getEnv("LISTEN_ADDRESS", ":8080"),
		GitLabSecretToken: getEnv("GITLAB_SECRET_TOKEN", ""), // optional
		LogAPIKey:         getEnv("LOG_API_KEY", ""),         // optional
		TelegramBotToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:    getEnv("TELEGRAM_CHAT_ID", ""),
		TelegramThreadID:  getEnv("TELEGRAM_THREAD_ID", ""), // optional
	}

	if config.TelegramBotToken == "" {
		slog.Error("TELEGRAM_BOT_TOKEN is not set")
		err = ErrMissingTelegramBotToken
		return
	}

	if config.TelegramChatID == "" {
		slog.Error("TELEGRAM_CHAT_ID is not set")
		err = ErrMissingTelegramChatID
		return
	}

	return
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
