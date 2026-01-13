package src

import (
	"log/slog"
	"net/http"

	"webhook-dispatcher/src/api"
	"webhook-dispatcher/src/config"

	"github.com/joho/godotenv"
)

func Main() int {
	godotenv.Load()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return 1
	}

	dispatcher := api.NewDispatcher()

	// API handlers
	http.Handle("/api/webhook", api.HandleWebhook(dispatcher, cfg))
	http.HandleFunc("/health", api.HandleHealth)

	slog.Info("Server starting on " + cfg.ListenAddress)
	slog.Error(http.ListenAndServe(cfg.ListenAddress, nil).Error())
	return 0
}
