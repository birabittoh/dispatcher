package models

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/birabittoh/dispatcher/src/config"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsnFmt = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

type Log struct {
	ID        uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level" gorm:"index"`
	Message   string            `json:"message"`
	Args      map[string]string `json:"args,omitempty" gorm:"type:jsonb"`
	Source    string            `json:"source,omitempty" gorm:"index"`

	Sent bool `json:"-"`
}

func InitDB(logger *slog.Logger, cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	if cfg.PostgresHost != "" {
		// PostgreSQL
		dsn := fmt.Sprintf(dsnFmt, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

		dialector = postgres.Open(dsn)
		logger.Info("Using PostgreSQL database", "dsn", dsn)
	} else {
		// SQLite
		dataDir := filepath.Dir(cfg.DBPath)
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}

		dialector = sqlite.Open(cfg.DBPath)
		logger.Info("Using SQLite database", "path", cfg.DBPath)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&Log{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
