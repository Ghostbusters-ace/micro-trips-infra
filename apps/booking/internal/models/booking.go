package models

type Booking struct {
	ID        int    `json:"id"`
	TripID    int    `json:"trip_id"`
	UserEmail string `json:"user_email"`
	Status    string `json:"status"` // PENDING, CONFIRMED, CANCELLED
}

type BookingRequest struct {
	TripID      int    `json:"trip_id"`
	UserEmail   string `json:"user_email"`
}
