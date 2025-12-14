package apiserver

import (
	"net/http"

	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

type APIServer struct {
	store  storage.Storage
	Router *http.ServeMux
}

func NewAPIServer(store storage.Storage) *APIServer {
	s := &APIServer{
		store:  store,
		Router: http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *APIServer) routes() {
	s.Router.HandleFunc("/resources", s.handleList)
	s.Router.HandleFunc("/resource", s.handleCreate)
	s.Router.HandleFunc("/resource/status", s.handleStatusUpdate)
}
