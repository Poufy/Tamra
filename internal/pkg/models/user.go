package models

import (
	"time"
)

type User struct {
	ID                string    `json:"id"`
	Longitude         float64   `json:"longitude" validate:"required"`
	Latitude          float64   `json:"latitude" validate:"required"`
	IsActive          bool      `json:"is_active" validate:"required"`
	Phone             string    `json:"phone" validate:"required,e164"`
	Radius            int       `json:"radius" validate:"required"`
	FCMToken          string    `json:"fcm_token" validate:"required"`
	LastOrderReceived time.Time `json:"last_order_received"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	IsActive  *bool   `json:"is_active" validate:"required"` // Pointer to a bool so the validation library doesn't complain if the value is false
	Phone     string  `json:"phone" validate:"required,e164"`
	Radius    int     `json:"radius" validate:"required"`
	FCMToken  string  `json:"fcm_token" validate:"required"`
}

type UpdateUserRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	IsActive  *bool   `json:"is_active" validate:"required"`
	Phone     string  `json:"phone" validate:"required,e164"`
	Radius    int     `json:"radius" validate:"required"`
	FCMToken  string  `json:"fcm_token" validate:"required"`
}

type UserResponse struct {
	ID                string    `json:"id"`
	Longitude         float64   `json:"longitude"`
	Latitude          float64   `json:"latitude"`
	IsActive          bool      `json:"is_active"`
	Phone             string    `json:"phone"`
	Radius            int       `json:"radius"`
	LastOrderReceived time.Time `json:"last_order_received"`
}
