package src

import (
	"net/http"

	"github.com/birabittoh/dispatcher/src/api"
	"github.com/birabittoh/dispatcher/src/config"
	"github.com/birabittoh/dispatcher/src/models"

	"github.com/joho/godotenv"
)

func Main() int {
	godotenv.Load()

	cfg := config.LoadConfig()
	logger := cfg.GetLogger()

	err := cfg.Validate()
	if err != nil {
		logger.Error("Invalid config", "error", err)
		return 1
	}

	db, err := models.InitDB(logger, cfg)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		return 1
	}

	m := api.NewManager(logger, cfg, db)

	mux := m.GetServeMultiplexer()

	logger.Info("Server starting on " + cfg.ListenAddress)
	logger.Error(http.ListenAndServe(cfg.ListenAddress, mux).Error())
	return 0
}
