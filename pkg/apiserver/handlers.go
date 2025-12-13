package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/Phoenix1504e/mini-control-plane/pkg/admission"
	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

type Handler struct {
	Store storage.Storage
}

func (h *Handler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var res api.Resource
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Admission
	if err := admission.Validate(&res); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if err := h.Store.Create(&res); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListResources(w http.ResponseWriter, _ *http.Request) {
	resources, err := h.Store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(resources)
}
func (h *Handler) WatchResources(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	for event := range h.Store.(*storage.FileStorage).WatchCh {
		_ = json.NewEncoder(w).Encode(event)
		flusher.Flush()
	}
}
