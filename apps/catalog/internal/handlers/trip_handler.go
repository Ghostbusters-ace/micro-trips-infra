package handlers

import (
	"encoding/json"
	"net/http"
	"catalog-api/internal/services"
)

type TripHandler struct {
	service *services.TripService
}

func NewTripHandler(service *services.TripService) *TripHandler {
	return &TripHandler{service: service}
}

// L'équivalent d'un @GetMapping en Spring
func (handler *TripHandler) GetTrips(writter http.ResponseWriter, request *http.Request) {
	writter.Header().Set("Content-Type", "application/json")
	
	trips := handler.service.GetAllAvailableTrips()
	
	writter.WriteHeader(http.StatusOK)
	json.NewEncoder(writter).Encode(trips)
}
