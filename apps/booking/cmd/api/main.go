package main

import (
	"log"
	"net/http"
	"os"

	"booking-api/internal/database"
	"booking-api/internal/handlers"
	"booking-api/internal/messaging"
	"booking-api/internal/middlewares"
	"booking-api/internal/repositories"
	"booking-api/internal/services"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("Initialisation du microservice Booking...")

	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "stable"
	}

	// ... (Initialisation DB, RabbitMQ, Repo, Publisher, Service, Handler identiques) ...
	dbConnection := database.InitDB()
	defer dbConnection.Close()
	rabbitClient := messaging.InitRabbitMQ()
	defer rabbitClient.Close()

	bookingRepo := repositories.NewPostgresBookingRepository(dbConnection)
	eventPublisher := messaging.NewRabbitEventPublisher(rabbitClient.Channel)
	bookingService := services.NewBookingService(bookingRepo, eventPublisher)
	bookingHandler := handlers.NewBookingHandler(bookingService)


	http.HandleFunc("/api/v1/bookings", middlewares.PrometheusMiddleware(appVersion, bookingHandler.HandleCreateBooking))

	// Route pour Prometheus
	http.Handle("/metrics", promhttp.Handler())

	port := "8080"
	log.Printf("API Booking à l'écoute sur le port %s (Version: %s)", port, appVersion)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Crash du serveur HTTP : %v", err)
	}
}