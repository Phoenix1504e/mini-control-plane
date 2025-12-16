package reconciler

import (

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
	// Always fetch the latest version for MVCC correctness
	current, err := r.Store.Get(res.Metadata.Name)
	if err != nil {
		return err
	}

	current.Status.CurrentReplicas = current.Spec.Replicas

	if err := r.Store.UpdateStatus(current); err != nil {
		return err
	}

	return nil
}
