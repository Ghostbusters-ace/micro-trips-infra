package services

import (
	"catalog-api/internal/models"
	"catalog-api/internal/repositories"
)

type TripService struct {
	repo *repositories.TripRepository
}

// L'équivalent de l'injection de dépendance par constructeur en Java
func NewTripService(repo *repositories.TripRepository) *TripService {
	return &TripService{repo: repo}
}

func (s *TripService) GetAllAvailableTrips() []models.Trip {
	// Ici on pourrait ajouter de la logique métier (ex: filtrer les trajets pleins)
	return s.repo.FindAll()
}
