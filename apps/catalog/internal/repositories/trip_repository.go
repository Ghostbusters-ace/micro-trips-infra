package repositories

import (
	"database/sql"
	"catalog-api/internal/models"
)

type TripRepository interface {
	GetAll() ([]models.Trip, error)
}

type postgresTripRepository struct {
	db *sql.DB
}

func NewPostgresTripRepository(db *sql.DB) TripRepository {
	return &postgresTripRepository{db: db}
}

func (r *postgresTripRepository) GetAll() ([]models.Trip, error) {
	rows, err := r.db.Query("SELECT id, destination, price, description FROM trips")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var t models.Trip
		if err := rows.Scan(&t.ID, &t.Destination, &t.Price, &t.Description); err != nil {
			return nil, err
		}
		trips = append(trips, t)
	}
	return trips, nil
}