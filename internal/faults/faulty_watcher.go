package faults

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/Phoenix1504e/mini-control-plane/internal/telemetry"
)

var logger *telemetry.Logger

func init() {
	var err error

	logger, err = telemetry.NewLogger("faults.jsonl")
	if err != nil {
		log.Fatalf("failed to initialize telemetry logger: %v", err)
	}
}

func WrapWatchChannel(
	ctx context.Context,
	upstream clientv3.WatchChan,
	injector Injector,
) clientv3.WatchChan {

	downstream := make(chan clientv3.WatchResponse)

	go func() {
		defer close(downstream)

		for {
			select {

			case <-ctx.Done():
				return

			case resp, ok := <-upstream:
				if !ok {
					return
				}

				log.Println("[FAULT] wrapper received watch response")

				decision := injector.Decide()

				switch decision.Action {

				case Drop:

					resourceKey := ""

					if len(resp.Events) > 0 && resp.Events[0].Kv != nil {
						resourceKey = string(resp.Events[0].Kv.Key)
					}

					log.Printf(
						"[FAULT] dropped watch response events=%d",
						len(resp.Events),
					)

					err := logger.Emit(telemetry.Event{
						Timestamp: time.Now(),
						Type:      "watch_event_dropped",
						Resource:  resourceKey,
						Message:   "watch response intentionally dropped",
					})

					if err != nil {
						log.Printf(
							"[FAULT] failed to emit telemetry: %v",
							err,
						)
					}

					continue

				default:
					downstream <- resp
				}
			}
		}
	}()

	return downstream
}
