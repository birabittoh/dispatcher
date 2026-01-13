package api

import (
	"errors"
	"net/http"
	"strconv"
)

func toURLValues(data map[string]string) map[string][]string {
	values := make(map[string][]string)
	for k, v := range data {
		values[k] = []string{v}
	}
	return values
}

func sendTelegramMessage(botToken, chatID, threadID, message string) error {
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

	payload := map[string]string{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "MarkdownV2",
	}
	if threadID != "" {
		payload["message_thread_id"] = threadID
	}

	// send the POST request and check that status is 2xx
	resp, err := http.PostForm(url, toURLValues(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("got status code " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
