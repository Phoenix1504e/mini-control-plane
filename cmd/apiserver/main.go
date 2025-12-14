package main

import (
	"log"
	"net/http"

	"github.com/Phoenix1504e/mini-control-plane/pkg/apiserver"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

func main() {
	// Initialize etcd-backed storage
	store, err := storage.NewEtcdStorage("/resources")
	if err != nil {
		log.Fatalf("failed to connect to etcd: %v", err)
	}

	// Create API server with router
	server := apiserver.NewAPIServer(store)

	log.Println("listening on :8080")

	// IMPORTANT: use server.Router, NOT nil
	if err := http.ListenAndServe(":8080", server.Router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
