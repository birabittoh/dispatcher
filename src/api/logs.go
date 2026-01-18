package api

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/birabittoh/dispatcher/src/models"

	"github.com/birabittoh/logs"
)

func (m Manager) HandleLog(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if m.cfg.LogAPIKey != "" {
		apiKey := r.Header.Get("X-API-Key")
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(m.cfg.LogAPIKey)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	defer r.Body.Close()
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		m.logger.Error("failed to read log payload", "error", err)
		http.Error(w, "Failed to read log payload", http.StatusInternalServerError)
		return
	}

	var logsBatch []models.Log
	err = json.Unmarshal(buf, &logsBatch)
	if err != nil {
		m.logger.Error("failed to unmarshal log payload", "error", err)
		http.Error(w, "Failed to unmarshal log payload", http.StatusBadRequest)
		return
	}

	for i, log := range logsBatch {
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

		if log.Level == logs.WARN || log.Level == logs.ERROR {
			err = sendTelegramMessage(m.cfg.TelegramBotToken, m.cfg.TelegramChatID, m.cfg.TelegramThreadID, msg.String(), false)
			logsBatch[i].Sent = err == nil
		}
	}

	err = m.db.Create(logsBatch).Error
	if err != nil {
		m.logger.Error("failed to store logs in database", "error", err)
		http.Error(w, "Failed to store logs in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
