package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/birabittoh/dispatcher/src/config"

	gitlabwebhook "github.com/flc1125/go-gitlab-webhook/v2"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewDispatcher() *gitlabwebhook.Dispatcher {
	return gitlabwebhook.NewDispatcher(
		gitlabwebhook.RegisterListeners(
			&telegramListener{},
		),
	)
}

func HandleWebhook(dispatcher *gitlabwebhook.Dispatcher, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, "TelegramBotToken", cfg.TelegramBotToken)
		ctx = context.WithValue(ctx, "TelegramChatID", cfg.TelegramChatID)
		ctx = context.WithValue(ctx, "TelegramThreadID", cfg.TelegramThreadID)

		opts := []gitlabwebhook.DispatchRequestOption{gitlabwebhook.DispatchRequestWithContext(ctx)}
		if cfg.GitLabSecretToken != "" {
			opts = append(opts, gitlabwebhook.DispatchRequestWithToken(cfg.GitLabSecretToken))
		}

		if err := dispatcher.DispatchRequest(r, opts...); err != nil {
			slog.Error("failed to dispatch webhook request", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleLog(cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if cfg.LogAPIKey != "" {
			apiKey := r.Header.Get("X-API-Key")
			if subtle.ConstantTimeCompare([]byte(apiKey), []byte(cfg.LogAPIKey)) != 1 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		defer r.Body.Close()
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("failed to read log payload", "error", err)
			http.Error(w, "Failed to read log payload", http.StatusInternalServerError)
			return
		}

		var log Log
		err = json.Unmarshal(buf, &log)
		if err != nil {
			slog.Error("failed to unmarshal log payload", "error", err)
			http.Error(w, "Failed to unmarshal log payload", http.StatusBadRequest)
			return
		}

		msg := strings.Builder{}
		if log.Source != "" {
			msg.WriteString("*Source:* " + log.Source + "\n")
		}
		msg.WriteString("*Level:* " + log.Level + "\n")
		msg.WriteString("*Message:* " + log.Message + "\n")

		if len(log.Args) > 0 {
			msg.WriteString("*Args:*\n")
			for k, v := range log.Args {
				msg.WriteString("  - *" + k + "*: " + v + "\n")
			}
		}

		if log.Level != "WARN" && log.Level != "ERROR" {
			err = sendTelegramMessage(cfg.TelegramBotToken, cfg.TelegramChatID, cfg.TelegramThreadID, msg.String(), false)
			log.Sent = err == nil
		}

		// TODO: insert log in db
		// err = db.Create(&log).Error

		w.WriteHeader(http.StatusNoContent)
	})
}
