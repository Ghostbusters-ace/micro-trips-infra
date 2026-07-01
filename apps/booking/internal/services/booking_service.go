package services

import (
	"booking-api/internal/models"
	"booking-api/internal/repositories"
	"fmt"
	"time"
)

type BookingService struct {
	repo *repositories.BookingRepository
}

func NewBookingService(repo *repositories.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) CreateBooking(tripID string, user string) models.Booking {
	// 1. Création de l'objet avec statut PENDING
	booking := models.Booking{
		ID:     fmt.Sprintf("RES-%d", time.Now().Unix()),
		TripID: tripID,
		User:   user,
		Status: "PENDING",
	}

	// 2. Sauvegarde en BDD
	s.repo.Save(booking)

	// 3. Envoi d'un événement à RabbitMQ (Simulé pour l'instant)
	fmt.Printf("[RabbitMQ] Message envoyé : 'Traiter la réservation %s'\n", booking.ID)

	return booking
}
