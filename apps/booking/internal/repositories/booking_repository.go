package repositories

import (
	"booking-api/internal/models"
	"fmt"
)

type BookingRepository struct {
	// vraie connexion PostgreSQL plus tard
	bookings map[string]models.Booking
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{
		bookings: make(map[string]models.Booking),
	}
}

func (r *BookingRepository) Save(booking models.Booking) {
	r.bookings[booking.ID] = booking
	fmt.Printf("[DB] Réservation %s sauvegardée avec le statut %s\n", booking.ID, booking.Status)
}
