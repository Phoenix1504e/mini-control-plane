package reconciler

import (
	"log"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

type Reconciler struct {
	Store storage.Storage
}

func New(store storage.Storage) *Reconciler {
	return &Reconciler{
		Store: store,
	}
}

// Reconcile converges actual state toward desired state
func (r *Reconciler) Reconcile(res *api.Resource) error {

	// Always fetch latest version for MVCC correctness
	current, err := r.Store.Get(res.Spec.Name)
	if err != nil {
		return err
	}

	// ----------------------------------
	// Finalizer / Deletion handling
	// ----------------------------------
	if current.Metadata.DeletionTimestamp != "" {

		// Cleanup phase
		if len(current.Metadata.Finalizers) > 0 {

			log.Printf(
				"Running cleanup for resource %s",
				current.Spec.Name,
			)

			// Simulated cleanup complete
			current.Metadata.Finalizers = nil

			if err := r.Store.Update(current); err != nil {
				return err
			}

			log.Printf(
				"Removed finalizers from resource %s",
				current.Spec.Name,
			)

			return nil
		}

		// Garbage collection phase
		log.Printf(
			"Garbage collecting resource %s",
			current.Spec.Name,
		)

		return r.Store.Delete(current.Spec.Name)
	}

	// ----------------------------------
	// Normal reconciliation
	// ----------------------------------
	if current.Status.CurrentReplicas == current.Spec.Replicas {
		return nil
	}

	current.Status.CurrentReplicas = current.Spec.Replicas

	if err := r.Store.UpdateStatus(current); err != nil {
		return err
	}

	return nil
}
