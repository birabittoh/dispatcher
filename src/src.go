package src

import (
	"log/slog"
	"net/http"
	"os"

	"backend-example/src/api"
	"backend-example/src/config"
	"backend-example/src/ui"

	"github.com/joho/godotenv"
)

// Configuration holds all app settings

func Main() int {
	godotenv.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	config := config.LoadConfig()

	// API handlers
	http.HandleFunc("/api/sum", api.HandleSum)
	http.HandleFunc("/health", api.HandleHealth)

	// UI handlers
	http.HandleFunc("/", ui.HandleIndex)

	logger.Info("Server starting on " + config.ListenAddress)
	logger.Error(http.ListenAndServe(config.ListenAddress, nil).Error())
	return 0
}
