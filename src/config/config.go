package config

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/birabittoh/logs"
	"github.com/lmittmann/tint"
)

type Config struct {
	LogLevel string

	DBPath           string
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

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

func LoadConfig() *Config {
	return &Config{
		ListenAddress:     getEnv("LISTEN_ADDRESS", ":8080"),
		GitLabSecretToken: getEnv("GITLAB_SECRET_TOKEN", ""), // optional
		LogAPIKey:         getEnv("LOG_API_KEY", ""),         // optional
		TelegramBotToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:    getEnv("TELEGRAM_CHAT_ID", ""),
		TelegramThreadID:  getEnv("TELEGRAM_THREAD_ID", ""), // optional
	}
}

func (c *Config) GetLogger() *slog.Logger {
	level := logs.ParseLogLevel(c.LogLevel)
	handler := tint.NewHandler(os.Stdout, &tint.Options{Level: level, TimeFormat: time.RFC3339})
	return slog.New(handler)
}

func (c *Config) Validate() error {
	if c.TelegramBotToken == "" {
		return ErrMissingTelegramBotToken
	}
	if c.TelegramChatID == "" {
		return ErrMissingTelegramChatID
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
