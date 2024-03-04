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

type CreateRestaurantRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	ImageURL  string  `json:"image_url" validate:"required,url"`
	Name      string  `json:"name" validate:"required"`
}

type UpdateRestaurantRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	ImageURL  string  `json:"image_url" validate:"required,url"`
	Name      string  `json:"name" validate:"required"`
}

type RestaurantResponse struct {
	ID        int       `json:"id"`
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	ImageURL  string    `json:"image_url"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
