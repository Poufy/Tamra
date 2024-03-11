package models

import (
	"time"
)

type Restaurant struct {
	ID        string    `json:"id"`
	Longitude float64   `json:"longitude" validate:"required"`
	Latitude  float64   `json:"latitude" validate:"required"`
	LogoURL   string    `json:"logo_url" validate:"required,url"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRestaurantRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	LogoURL   string  `json:"logo_url" validate:"required,url"`
	Name      string  `json:"name" validate:"required"`
}

type UpdateRestaurantRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	LogoURL   string  `json:"logo_url" validate:"required,url"`
	Name      string  `json:"name" validate:"required"`
}

type RestaurantResponse struct {
	ID        string    `json:"id"`
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	LogoURL   string    `json:"logo_url"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RestaurantLogoUploadResponse struct {
	PresignedURL  string `json:"presigned_url"`
	StoredFileURL string `json:"stored_file_url"`
	Description   string `json:"description"`
}
