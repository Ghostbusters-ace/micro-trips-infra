package services

import (
	"catalog-api/internal/models"
	"catalog-api/internal/repositories"
)

type CatalogService interface {
	GetTrips() ([]models.Trip, error)
}

type catalogService struct {
	repo repositories.TripRepository
}

func NewTripService(repo repositories.TripRepository) CatalogService {
	return &catalogService{repo: repo}
}

func (s *catalogService) GetTrips() ([]models.Trip, error) {
	return s.repo.GetAll()
}