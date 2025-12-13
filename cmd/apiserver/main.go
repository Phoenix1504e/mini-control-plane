package main

import (
	"net/http"

	"github.com/Phoenix1504e/mini-control-plane/pkg/apiserver"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

func main() {
	store := storage.NewFileStorage("specs")

	handler := &apiserver.Handler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/resources", handler.CreateResource)
	mux.HandleFunc("/resources/list", handler.ListResources)
        mux.HandleFunc("/watch/resources", handler.WatchResources)

	server := &apiserver.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.Run()
}
