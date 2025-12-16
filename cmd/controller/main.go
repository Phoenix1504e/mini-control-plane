package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/Phoenix1504e/mini-control-plane/pkg/informer"
        "github.com/Phoenix1504e/mini-control-plane/pkg/scheduler"
	"github.com/Phoenix1504e/mini-control-plane/pkg/leader"
	"github.com/Phoenix1504e/mini-control-plane/pkg/reconciler"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
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
        sched := scheduler.New(store, []string{
	"node-a",
	"node-b",
	"node-c",
        })


	// Leader election
	elector := leader.New(cli, "/control-plane/leader", controllerID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start leader election loop
	go func() {
		if err := elector.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Start informer
	go inf.Start(ctx)

	// Controller event loop
	for ev := range inf.EventChan {
	if !elector.IsLeader() {
		continue
	}

	if err := rec.Reconcile(ev.Resource); err != nil {
		log.Println("Reconcile error:", err)
		continue
	}

	if err := sched.Schedule(ev.Resource); err != nil {
		log.Println("Schedule error:", err)
	}
}
}
