package telemetry

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`

	Resource string `json:"resource,omitempty"`

	Message string `json:"message,omitempty"`
}

type Logger struct {
	mu  sync.Mutex
	enc *json.Encoder
}

func NewLogger(path string) (*Logger, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &Logger{
		enc: json.NewEncoder(f),
	}, nil
}

func (l *Logger) Emit(e Event) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.enc.Encode(e)
}
