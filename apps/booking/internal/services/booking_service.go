package services

import (
	"booking-api/internal/messaging"
	"booking-api/internal/models"
	"booking-api/internal/repositories"
)

type BookingService interface {
	CreateBooking(tripID int, email string) (*models.Booking, error)
}

type bookingService struct {
	repo      repositories.BookingRepository
	publisher messaging.EventPublisher
}

func NewBookingService(repo repositories.BookingRepository, pub messaging.EventPublisher) BookingService {
	return &bookingService{repo: repo, publisher: pub}
}

func (s *bookingService) CreateBooking(tripID int, email string) (*models.Booking, error) {
	booking := &models.Booking{
		TripID:    tripID,
		UserEmail: email,
		Status:    "PENDING",
	}

	id, err := s.repo.Create(booking)
	if err != nil {
		return nil, err
	}
	booking.ID = id

	// Publication asynchrone RabbitMQ via l'interface d'abstraction
	_ = s.publisher.PublishBookingCreated(booking.ID, booking.UserEmail)

	return booking, nil
}