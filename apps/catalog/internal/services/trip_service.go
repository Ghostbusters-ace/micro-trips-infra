package services

import (
	"catalog-api/internal/models"
	"catalog-api/internal/repositories"
)

type CatalogService interface {
	GetTrips() ([]models.Trip, error)
}

type catalogService struct {
	repo repositories.CatalogRepository
}

func NewCatalogService(repo repositories.CatalogRepository) CatalogService {
	return &catalogService{repo: repo}
}

func (s *catalogService) GetTrips() ([]models.Trip, error) {
	return s.repo.GetAll()
}