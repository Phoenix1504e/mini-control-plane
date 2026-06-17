package apiserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
)

func (s *APIServer) handleResource(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGet(w, r)
	case http.MethodPost:
		s.handleCreate(w, r)
	case http.MethodDelete:
		s.handleDelete(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *APIServer) handleGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing resource name", http.StatusBadRequest)
		return
	}

	res, err := s.store.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (s *APIServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	var res api.Resource

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if res.Metadata.Name == "" {
		res.Metadata.Name = res.Spec.Name
	}

	if res.Metadata.Name != res.Spec.Name {
		http.Error(w, "metadata.name and spec.name must match", http.StatusBadRequest)
		return
	}

	if err := s.store.Create(&res); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *APIServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing resource name", http.StatusBadRequest)
		return
	}

	res, err := s.store.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Finalizers present -> soft delete
	if len(res.Metadata.Finalizers) > 0 {
		if res.Metadata.DeletionTimestamp == "" {
			res.Metadata.DeletionTimestamp = time.Now().UTC().Format(time.RFC3339)

			if err := s.store.Update(res); err != nil {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
		}

		w.WriteHeader(http.StatusAccepted)
		return
	}

	// No finalizers -> hard delete
	if err := s.store.Delete(name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
