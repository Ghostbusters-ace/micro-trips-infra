package main

import (
	"log"
	"net/http"

	"catalog-api/internal/handlers"
	"catalog-api/internal/repositories"
	"catalog-api/internal/services"
)

func main() {
	// 1. Initialisation des couches (Injection de dépendances manuelle)
	repo := repositories.NewTripRepository()
	service := services.NewTripService(repo)
	handler := handlers.NewTripHandler(service)

	// 2. Définition des routes
	http.HandleFunc("/trips", handler.GetTrips)

	// 3. Démarrage du serveur web
	log.Println("Catalog API démarrée sur le port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur : ", err)
	}
}
