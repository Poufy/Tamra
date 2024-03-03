package models

import (
	"time"
)

type User struct {
	ID                int       `json:"id"`
	Longitude         float64   `json:"longitude" validate:"required"`
	Latitude          float64   `json:"latitude" validate:"required"`
	IsActive          bool      `json:"is_active"`
	Phone             string    `json:"phone" validate:"required,e164"`
	Radius            float64   `json:"radius" validate:"required"`
	LastOrderReceived time.Time `json:"last_order_received"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
