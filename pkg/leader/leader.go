package leader

import (
	"fmt"
	"os"
	"time"
)

const lockFile = "leader.lock"

func TryAcquire(id string) (bool, error) {
	f, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return false, nil // someone else is leader
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("leader=%s\n", id))
	if err != nil {
		return false, err
	}

	return true, nil
}

func Release() {
	_ = os.Remove(lockFile)
}

func HoldLeadership(id string) {
	for {
		time.Sleep(2 * time.Second)
		// heartbeat placeholder (real systems renew lease here)
	}
}
