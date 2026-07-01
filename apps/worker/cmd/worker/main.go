package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println("Worker démarré : connexion à RabbitMQ en cours...")
	log.Println("✅ Connecté ! En attente de messages...")

	// Simulation d'une écoute continue de la file d'attente
	for {
		time.Sleep(10 * time.Second)
		
		// Simulation : on prétend avoir reçu le message de la Booking API
		fmt.Println("--------------------------------------------------")
		fmt.Println("[RabbitMQ] Message reçu : 'Traiter la réservation RES-1782939836'")
		fmt.Println("Traitement en cours (vérification paiement...)")
		
		time.Sleep(2 * time.Second) // Simule un traitement long
		
		fmt.Println("✅ [DB] Réservation RES-1782939836 mise à jour avec le statut CONFIRMED")
		fmt.Println("--------------------------------------------------")
	}
}
