package main

import (
	"log"
	"net/http"
	"os"

	"booking-api/internal/handlers"
	"booking-api/internal/middlewares"
	"booking-api/internal/repositories"
	"booking-api/internal/services"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// Récupération de la version depuis l'environnement (fournie par le manifeste K8s)
	// Si elle n'est pas définie, on met "stable" par défaut
	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "stable"
	}

	repo := repositories.NewBookingRepository()
	service := services.NewBookingService(repo)
	handler := handlers.NewBookingHandler(service)

	// Route Prometheus
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/bookings", middlewares.PrometheusMiddleware(appVersion, handler.PostBooking))

	log.Println("Booking API démarrée sur le port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Erreur : ", err)
	}
}
