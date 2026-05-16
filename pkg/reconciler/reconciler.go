package reconciler

import (
	"log"
	"time"

	"github.com/Phoenix1504e/mini-control-plane/internal/telemetry"
	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

type Reconciler struct {
	Store  storage.Storage
	Logger *telemetry.Logger
}

func New(store storage.Storage) *Reconciler {

	logger, err := telemetry.NewLogger("reconcile.jsonl")
	if err != nil {
		log.Fatalf(
			"failed to initialize telemetry logger: %v",
			err,
		)
	}

	return &Reconciler{
		Store:  store,
		Logger: logger,
	}
}

// Reconcile converges actual state toward desired state
func (r *Reconciler) Reconcile(res *api.Resource) error {

	start := time.Now()

	resourceName := ""

	if res != nil {
		resourceName = res.Metadata.Name
	}

	err := r.Logger.Emit(telemetry.Event{
		Timestamp: start,
		Type:      "reconcile_start",
		Resource:  resourceName,
	})

	if err != nil {
		log.Printf(
			"failed to emit reconcile_start telemetry: %v",
			err,
		)
	}

	// Always fetch the latest version for MVCC correctness
	current, err := r.Store.Get(res.Metadata.Name)
	if err != nil {

		emitErr := r.Logger.Emit(telemetry.Event{
			Timestamp: time.Now(),
			Type:      "reconcile_error",
			Resource:  resourceName,
			Message:   err.Error(),
		})

		if emitErr != nil {
			log.Printf(
				"failed to emit reconcile_error telemetry: %v",
				emitErr,
			)
		}

		return err
	}

	current.Status.CurrentReplicas = current.Spec.Replicas

	if err := r.Store.UpdateStatus(current); err != nil {

		emitErr := r.Logger.Emit(telemetry.Event{
			Timestamp: time.Now(),
			Type:      "reconcile_error",
			Resource:  resourceName,
			Message:   err.Error(),
		})

		if emitErr != nil {
			log.Printf(
				"failed to emit reconcile_error telemetry: %v",
				emitErr,
			)
		}

		return err
	}

	duration := time.Since(start)

	err = r.Logger.Emit(telemetry.Event{
		Timestamp: time.Now(),
		Type:      "reconcile_complete",
		Resource:  resourceName,
		Message:   duration.String(),
	})

	if err != nil {
		log.Printf(
			"failed to emit reconcile_complete telemetry: %v",
			err,
		)
	}

	return nil
}
