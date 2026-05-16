package metrics

import (
	"context"
	"log"
	"time"

	"github.com/Phoenix1504e/mini-control-plane/internal/telemetry"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

type StateSampler struct {
	store  storage.Storage
	logger *telemetry.Logger
}

func NewStateSampler(
	store storage.Storage,
	logger *telemetry.Logger,
) *StateSampler {

	return &StateSampler{
		store:  store,
		logger: logger,
	}
}

func (s *StateSampler) Start(
	ctx context.Context,
	interval time.Duration,
) {

	ticker := time.NewTicker(interval)

	go func() {

		defer ticker.Stop()

		for {
			select {

			case <-ctx.Done():
				return

			case <-ticker.C:

				resources, err := s.store.List()

				if err != nil {
					log.Printf(
						"[SAMPLER] failed to list resources: %v",
						err,
					)

					continue
				}

				for _, res := range resources {

					if res == nil {
						continue
					}

					err := s.logger.Emit(telemetry.Event{
						Timestamp:       time.Now(),
						Type:            "state_sample",
						Resource:        res.Metadata.Name,
						DesiredReplicas: res.Spec.Replicas,
						CurrentReplicas: res.Status.CurrentReplicas,
						ResourceVersion: res.Metadata.ResourceVersion,
					})

					if err != nil {
						log.Printf(
							"[SAMPLER] failed to emit telemetry: %v",
							err,
						)
					}
				}
			}
		}
	}()
}
