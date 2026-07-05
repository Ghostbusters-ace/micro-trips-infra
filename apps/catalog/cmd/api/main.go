package main

import (
	"log"
	"net/http"
	"os"

	"catalog-api/internal/handlers"
	"catalog-api/internal/middlewares"
	"catalog-api/internal/repositories"
	"catalog-api/internal/services"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// Récupération de la version depuis l'environnement (fournie par le manifeste K8s)
	// Si elle n'est pas définie, on met "stable" par défaut
	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "stable"
	}

	// 1. Initialisation des couches (Injection de dépendances manuelle)
	repo := repositories.NewTripRepository()
	service := services.NewTripService(repo)
	handler := handlers.NewTripHandler(service)

	// Route Prometheus
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/trips", middlewares.PrometheusMiddleware(appVersion, handler.GetTrips))

	// 2. Démarrage du serveur web
	log.Println("Catalog API démarrée sur le port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur : ", err)
	}
}