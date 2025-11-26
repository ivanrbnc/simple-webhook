package domain

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}

type Subscription struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

type TriggerEventRequest struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
