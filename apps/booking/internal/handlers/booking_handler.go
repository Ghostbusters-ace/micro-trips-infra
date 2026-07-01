package handlers

import (
	"encoding/json"
	"net/http"
	"booking-api/internal/services"
)

// Requete entrante
type CreateBookingRequest struct {
	TripID string `json:"trip_id"`
	User   string `json:"user"`
}

type BookingHandler struct {
	service *services.BookingService
}

func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (handler *BookingHandler) PostBooking(responseBooking http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseBooking, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var bookingRequest CreateBookingRequest
	if err := json.NewDecoder(request.Body).Decode(&bookingRequest); err != nil {
		http.Error(responseBooking, err.Error(), http.StatusBadRequest)
		return
	}

	booking := handler.service.CreateBooking(bookingRequest.TripID, bookingRequest.User)

	// Réponse au client
	responseBooking.Header().Set("Content-Type", "application/json")
	responseBooking.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseBooking).Encode(booking)
}
