package events

import (
	"fmt"
	"os"
	"time"
)

const eventLogFile = "events.log"

func Record(resource string, message string) error {
	entry := fmt.Sprintf(
		"[%s] %s: %s\n",
		time.Now().Format(time.RFC3339),
		resource,
		message,
	)

	f, err := os.OpenFile(eventLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}
