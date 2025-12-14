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
	log.Printf(
		"Reconciling %s: desired=%d",
		res.Spec.Name,
		res.Spec.Replicas,
	)

	// Simulated reconciliation
	res.Status.CurrentReplicas = res.Spec.Replicas

	err := r.Store.UpdateStatus(res)
	if err != nil {
		log.Println("Status update failed:", err)
		return err
	}

	return nil
}
