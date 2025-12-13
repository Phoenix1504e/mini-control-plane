package reconciler

import (
	"fmt"
	"log"
	"time"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
	"github.com/Phoenix1504e/mini-control-plane/pkg/events"
	"github.com/Phoenix1504e/mini-control-plane/pkg/runtime"
)

func Reconcile(res *api.Resource) error {
	spec := res.Spec

	current, err := runtime.ListInstances(spec.Name)
	if err != nil {
		return err
	}

	diff := spec.Replicas - len(current)
        // Drift detection
if len(current) != spec.Replicas {
	message := fmt.Sprintf(
		"Drift detected for %s: desired=%d actual=%d",
		spec.Name,
		spec.Replicas,
		len(current),
	)

	log.Println(message)
	_ = events.Record(spec.Name, message)
}


	// Scale UP
	if diff > 0 {
		message := fmt.Sprintf(
			"Scaling up %s from %d → %d",
			spec.Name,
			len(current),
			spec.Replicas,
		)

		log.Println(message)
		_ = events.Record(spec.Name, message)

		for i := 0; i < diff; i++ {
			count, err := runtime.CountInstances(spec.Name)
			if err != nil {
				return err
			}
			if err := runtime.CreateInstance(spec.Name, count); err != nil {
				return err
			}
		}
	}

	// Scale DOWN
	if diff < 0 {
		message := fmt.Sprintf(
			"Scaling down %s from %d → %d",
			spec.Name,
			len(current),
			spec.Replicas,
		)

		log.Println(message)
		_ = events.Record(spec.Name, message)

		for i := 0; i < -diff; i++ {
			if err := runtime.DeleteInstance(current[i]); err != nil {
				return err
			}
		}
	}

	// Update STATUS (observed state)
	count, err := runtime.CountInstances(spec.Name)
	if err != nil {
		return err
	}

	res.Status.CurrentReplicas = count
	res.Status.LastReconciled = time.Now().Format(time.RFC3339)

	return nil
}
