package api

import (
	"context"
	"log/slog"
	"net/http"
	"webhook-dispatcher/src/config"

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

		// add cfg.TelegramBotToken, cfg.TelegramChatID, cfg.TelegramThreadID to context
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
