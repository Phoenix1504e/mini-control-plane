package main

import (
	"log"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
	"github.com/Phoenix1504e/mini-control-plane/pkg/informer"
	"github.com/Phoenix1504e/mini-control-plane/pkg/leader"
	"github.com/Phoenix1504e/mini-control-plane/pkg/reconciler"
	"github.com/Phoenix1504e/mini-control-plane/pkg/store"
)

func main() {
	// Unique controller identity
	id := uuid.NewString()
	log.Printf("Controller starting with ID %s", id)

	// =========================
	// LEADER ELECTION
	// =========================
	isLeader, err := leader.TryAcquire(id)
	if err != nil {
		log.Fatal("leader election failed:", err)
	}

	if !isLeader {
		log.Println("Not leader. Running in standby mode.")
		select {} // follower does nothing
	}

	log.Println("I am the leader")
	defer leader.Release()

	// =========================
	// EVENT-DRIVEN CONTROLLER
	// =========================
	log.Println("Starting informer (watch-based reconciliation)")

	informer.Watch("http://localhost:8080/watch/resources", func(event api.WatchEvent) {
		res := event.Resource

		// Safety: admission must be approved
		if !api.IsConditionTrue(res.Status.Conditions, "AdmissionApproved") {
			log.Printf(
				"Skipping reconcile for %s (admission not approved)",
				res.Spec.Name,
			)
			return
		}

		// =========================
		// RECONCILIATION
		// =========================
		if err := reconciler.Reconcile(res); err != nil {
			log.Println("reconcile error:", err)
			return
		}

		// =========================
		// PERSIST STATUS
		// =========================
		resourcePath := filepath.Join("specs", res.Spec.Name+".yaml")
		if err := store.SaveResource(resourcePath, res); err != nil {
			log.Println("failed to save resource status:", err)
		}
	})
}
