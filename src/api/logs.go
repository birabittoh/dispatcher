package api

import (
	"time"
)

type Log struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Args      map[string]string `json:"args,omitempty"`
	Source    string            `json:"source,omitempty"`

	Sent bool `json:"-"`
}
