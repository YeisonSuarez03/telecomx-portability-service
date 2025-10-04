package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"telecomx-portability-service/internal/application/service"
	"telecomx-portability-service/internal/domain/model"
)

type PortabilityHandler struct {
	service *service.PortabilityService
}

func NewPortabilityHandler(s *service.PortabilityService) *PortabilityHandler {
	return &PortabilityHandler{service: s}
}

func (h *PortabilityHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/portability", h.handlePortability)
}

func (h *PortabilityHandler) handlePortability(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	switch r.Method {
	case http.MethodGet:
		data, err := h.service.GetAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)

	case http.MethodPost:
		var p model.Portability
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.service.Create(ctx, &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
