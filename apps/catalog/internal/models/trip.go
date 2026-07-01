package models

type Trip struct {
	ID          string  `json:"id"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Price       float64 `json:"price"`
}
