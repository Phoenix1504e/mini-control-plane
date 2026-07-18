package events

import (
	"encoding/json"
	"os"
	"time"
)

// 1. Change extension to .jsonl
const eventLogFile = "events.jsonl"

// 2. Define a structured struct for clean marshalling
type EventEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Resource  string    `json:"resource"`
	Message   string    `json:"message"`
}

func Record(resource string, message string) error {
	entry := EventEntry{
		Timestamp: time.Now(),
		Resource:  resource,
		Message:   message,
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	// Append a newline character so each event is its own line
	jsonData = append(jsonData, '\n')

	f, err := os.OpenFile(eventLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(jsonData)
	return err
}
