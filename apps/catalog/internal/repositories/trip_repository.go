package repositories

import "catalog-api/internal/models"

type TripRepository struct{}

func NewTripRepository() *TripRepository {
	return &TripRepository{}
}

// FindAll simule un appel à la base de données PostgreSQL
func (r *TripRepository) FindAll() []models.Trip {
	return []models.Trip{
		{ID: "TRIP-123", Origin: "Paris", Destination: "Lyon", Price: 45.50},
		{ID: "TRIP-456", Origin: "Marseille", Destination: "Toulouse", Price: 30.00},
		{ID: "TRIP-789", Origin: "Lille", Destination: "Bordeaux", Price: 75.20},
	}
}
