package handlers

import (
	"encoding/json"
	"net/http"
	"booking-api/internal/services"
	"booking-api/internal/models"
)

type BookingHandler struct {
	service services.BookingService
}

func NewBookingHandler(service services.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) HandleCreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload JSON invalide", http.StatusBadRequest)
		return
	}

	booking, err := h.service.CreateBooking(req.TripID, req.UserEmail)
	if err != nil {
		http.Error(w, "Erreur interne lors du traitement", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}