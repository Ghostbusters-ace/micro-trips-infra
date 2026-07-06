package handlers

import (
	"encoding/json"
	"net/http"

	"catalog-api/internal/services"
)

type TripHandler struct {
	service services.CatalogService
}

func NewTripHandler(service services.CatalogService) *TripHandler {
	return &TripHandler{service: service}
}

func (h *TripHandler) HandleGetTrips(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	trips, err := h.service.GetTrips()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}