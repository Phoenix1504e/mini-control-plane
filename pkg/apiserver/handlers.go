package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
)

func (s *APIServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	var res api.Resource
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.store.Create(&res); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *APIServer) handleList(w http.ResponseWriter, _ *http.Request) {
	items, err := s.store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

func (s *APIServer) handleStatusUpdate(w http.ResponseWriter, r *http.Request) {
	var res api.Resource
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.store.UpdateStatus(&res); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}
