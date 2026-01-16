package logs

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

var logLevels = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

func Init(logLevelStr string) {
	logLevel, ok := logLevels[strings.ToUpper(logLevelStr)]
	if !ok {
		logLevel = slog.LevelInfo
	}

	slog.SetDefault(
		slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      logLevel,
				TimeFormat: time.RFC3339,
			}),
		),
	)
}
