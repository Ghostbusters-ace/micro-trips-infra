package main

import (
	"log"
	"net/http"
	"os"

	"Trip-api/internal/handlers"
	"Trip-api/internal/middlewares"
	"Trip-api/internal/repositories"
	"Trip-api/internal/services"
	"Trip-api/internal/database"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("Initialisation du service Trip...")

	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "stable"
	}

	dbConnection := database.InitDB()
	defer dbConnection.Close()

	TripRepo := repositories.NewPostgresTripRepository(dbConnection)
	TripService := services.NewTripService(TripRepo)
	TripHandler := handlers.NewTripHandler(TripService)

	http.HandleFunc("/api/v1/trips", middlewares.PrometheusMiddleware(appVersion, TripHandler.HandleGetTrips))
	http.Handle("/metrics", promhttp.Handler())

	port := "8080"
	log.Printf("Catalog API en écoute sur le port %s (Version: %s)", port, appVersion)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Crash : %v", err)
	}
}