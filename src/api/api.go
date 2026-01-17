package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/birabittoh/dispatcher/src/config"
	"gorm.io/gorm"

	gitlabwebhook "github.com/flc1125/go-gitlab-webhook/v2"
)

type Manager struct {
	logger     *slog.Logger
	cfg        *config.Config
	db         *gorm.DB
	dispatcher *gitlabwebhook.Dispatcher
}

func NewManager(logger *slog.Logger, cfg *config.Config, db *gorm.DB) *Manager {
	return &Manager{
		logger:     logger,
		cfg:        cfg,
		db:         db,
		dispatcher: newDispatcher(),
	}
}

func (m Manager) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func newDispatcher() *gitlabwebhook.Dispatcher {
	return gitlabwebhook.NewDispatcher(
		gitlabwebhook.RegisterListeners(
			&telegramListener{},
		),
	)
}

func (m Manager) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, "TelegramBotToken", m.cfg.TelegramBotToken)
	ctx = context.WithValue(ctx, "TelegramChatID", m.cfg.TelegramChatID)
	ctx = context.WithValue(ctx, "TelegramThreadID", m.cfg.TelegramThreadID)

	opts := []gitlabwebhook.DispatchRequestOption{gitlabwebhook.DispatchRequestWithContext(ctx)}
	if m.cfg.GitLabSecretToken != "" {
		opts = append(opts, gitlabwebhook.DispatchRequestWithToken(m.cfg.GitLabSecretToken))
	}

	if err := m.dispatcher.DispatchRequest(r, opts...); err != nil {
		m.logger.Error("failed to dispatch webhook request", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (m Manager) GetServeMultiplexer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/webhook", m.HandleWebhook)
	mux.HandleFunc("/api/log", m.HandleLog)
	mux.HandleFunc("/health", m.HandleHealth)

	return mux
}
