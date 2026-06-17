package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/Phoenix1504e/mini-control-plane/pkg/informer"
	"github.com/Phoenix1504e/mini-control-plane/pkg/leader"
	"github.com/Phoenix1504e/mini-control-plane/pkg/reconciler"
	"github.com/Phoenix1504e/mini-control-plane/pkg/scheduler"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
	"github.com/Phoenix1504e/mini-control-plane/pkg/workqueue"
)

func main() {
	// Unique controller identity
	controllerID := uuid.NewString()
	log.Println("Controller starting with ID", controllerID)

	// etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Storage
	store, err := storage.NewEtcdStorage("/resources")
	if err != nil {
		log.Fatal(err)
	}

	// Informer
	inf := informer.NewInformer(cli, "/resources")

	// Reconciler
	rec := reconciler.New(store)

	// Scheduler
	sched := scheduler.New(store, []string{
		"node-a",
		"node-b",
		"node-c",
	})

	// Workqueue
	queue := workqueue.New(100)

	// Leader election
	elector := leader.New(cli, "/control-plane/leader", controllerID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start leader election
	go func() {
		if err := elector.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Start informer
	go inf.Start(ctx)

	log.Println("Controller waiting for events...")

	// Startup resync after leadership acquisition
	go func() {
		for !elector.IsLeader() {
			time.Sleep(500 * time.Millisecond)
		}

		log.Println("Leader acquired, performing startup resync")

		items, err := store.List()
		if err != nil {
			log.Println("Resync failed:", err)
			return
		}

		for _, res := range items {

			log.Printf(
				"RESYNC: resource=%s",
				res.Spec.Name,
			)

			queue.Add(res.Spec.Name)
		}
	}()

	// Worker
	go func() {

		for {

			key := queue.Get()

			res, err := store.Get(key)
			if err != nil {

				// Resource already deleted
				queue.Forget(key)
				continue
			}

			log.Printf(
				"WORKER: processing %s",
				key,
			)

			if err := rec.Reconcile(res); err != nil {

				log.Printf(
					"Worker reconcile failed for %s: %v",
					key,
					err,
				)

				queue.AddRateLimited(key)
				continue
			}

			if err := sched.Schedule(res); err != nil {

				log.Printf(
					"Worker schedule failed for %s: %v",
					key,
					err,
				)

				queue.AddRateLimited(key)
				continue
			}

			queue.Forget(key)
		}
	}()

	// Controller event loop
	for ev := range inf.EventChan {

		log.Printf(
			"EVENT: %s resource=%s",
			ev.Type,
			ev.Resource.Spec.Name,
		)

		if !elector.IsLeader() {
			log.Printf(
				"Skipping event for %s: not leader",
				ev.Resource.Spec.Name,
			)
			continue
		}

		switch ev.Type {

		case informer.Deleted:
			log.Printf(
				"Resource deleted: %s",
				ev.Resource.Spec.Name,
			)
			continue

		case informer.Added, informer.Updated:
			log.Printf(
				"QUEUE: %s",
				ev.Resource.Spec.Name,
			)

			queue.Add(ev.Resource.Spec.Name)
		}
	}
}
