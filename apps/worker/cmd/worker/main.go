package main

import (
	"encoding/json"
	"log"

	"worker-service/internal/mailer"
	"worker-service/internal/messaging"
)

// EventPayload correspond à ce que Booking a envoyé
type EventPayload struct {
	BookingID int    `json:"booking_id"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}

func main() {
	log.Println("Démarrage du Notification Worker...")

	rabbit := messaging.InitRabbitMQ()
	defer rabbit.Close()

	queueName := "bookings_queue"
	q, err := rabbit.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la déclaration de la queue : %v", err)
	}

	msgs, err := rabbit.Channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la consommation de la queue : %v", err)
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Message RabbitMQ reçu : %s", string(d.Body))

			var event EventPayload
			if err := json.Unmarshal(d.Body, &event); err == nil {
				mailer.SendNotification(event.Email, event.BookingID)
			} else {
				log.Printf("Impossible de décoder le message : %v", err)
			}
		}
	}()

	log.Printf("Worker en attente de messages sur la queue [%s]...", queueName)
	<-forever
}