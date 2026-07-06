package models

type Booking struct {
	ID     string `json:"id"`
	TripID string `json:"trip_id"`
	User   string `json:"user"`
	Status string `json:"status"` // PENDING, CONFIRMED, CANCELLED
}

type BookingRequest struct {
	TripID      int    `json:"trip_id"`
	UserEmail   string `json:"user_email"`
}
