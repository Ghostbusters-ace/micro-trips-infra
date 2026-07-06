package repositories

import (
	"database/sql"
	"booking-api/internal/models"
)

type BookingRepository interface {
	Create(booking *models.Booking) (int, error)
}

type postgresBookingRepository struct {
	db *sql.DB
}

func NewPostgresBookingRepository(db *sql.DB) BookingRepository {
	return &postgresBookingRepository{db: db}
}

func (r *postgresBookingRepository) Create(b *models.Booking) (int, error) {
	query := `INSERT INTO bookings (trip_id, user_email, status) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRow(query, b.TripID, b.UserEmail, b.Status).Scan(&id)
	return id, err
}
