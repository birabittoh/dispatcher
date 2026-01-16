package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	telegramAPISendMessageURL = "https://api.telegram.org/bot%s/sendMessage"
	parseModeMarkdown         = "Markdown"
)

func isStatusOK(code int) bool {
	return code >= 200 && code < 300
}

func toURLValues(data map[string]string) map[string][]string {
	values := make(map[string][]string)
	for k, v := range data {
		values[k] = []string{v}
	}
	return values
}

func sendTelegramMessage(botToken, chatID, threadID, message string, silent bool) error {
	url := fmt.Sprintf(telegramAPISendMessageURL, botToken)

	payload := map[string]string{
		"chat_id":              chatID,
		"message_thread_id":    threadID,
		"text":                 message,
		"parse_mode":           parseModeMarkdown,
		"disable_notification": strconv.FormatBool(silent),
	}

	resp, err := http.PostForm(url, toURLValues(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !isStatusOK(resp.StatusCode) {
		return errors.New("got status code " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
