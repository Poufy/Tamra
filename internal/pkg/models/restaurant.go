package models

import (
	"time"
)

type Restaurant struct {
	ID        int       `json:"id"`
	Longitude float64   `json:"longitude" validate:"required"`
	Latitude  float64   `json:"latitude" validate:"required"`
	ImageURL  string    `json:"image_url" validate:"required,url"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
