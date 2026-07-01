package main

import (
	"log"
	"net/http"

	"booking-api/internal/handlers"
	"booking-api/internal/repositories"
	"booking-api/internal/services"
)

func main() {
	repo := repositories.NewBookingRepository()
	service := services.NewBookingService(repo)
	handler := handlers.NewBookingHandler(service)

	http.HandleFunc("/bookings", handler.PostBooking)

	log.Println("Booking API démarrée sur le port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Erreur : ", err)
	}
}
