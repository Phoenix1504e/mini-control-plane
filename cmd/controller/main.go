package main

import (
	"context"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/Phoenix1504e/mini-control-plane/pkg/informer"
	"github.com/Phoenix1504e/mini-control-plane/pkg/reconciler"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewEtcdStorage("/resources")
	if err != nil {
		log.Fatal(err)
	}

	inf := informer.NewInformer(cli, "/resources")
	ctx := context.Background()
	go inf.Start(ctx)

	rec := reconciler.New(store)

	log.Println("Controller running")

	for ev := range inf.EventChan {
		log.Printf("EVENT %s %s", ev.Type, ev.Resource.Spec.Name)
		rec.Reconcile(ev.Resource)
	}
}
