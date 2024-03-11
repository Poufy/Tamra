package models

import (
	"time"
)

type Order struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id" validate:"required"`
	RestaurantID string    `json:"restaurant_id" validate:"required"`
	Code         string    `json:"code" validate:"required"`
	Description  string    `json:"description"`
	State        string    `json:"state" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
	Description  string `json:"description"`
	RestaurantID string `json:"restaurant_id" validate:"required"`
}

type UpdateOrderRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	RestaurantID string `json:"restaurant_id" validate:"required"`
	Code         string `json:"code" validate:"required"`
	Description  string `json:"description"`
	State        string `json:"state" validate:"required"`
}

type OrderResponse struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	RestaurantID string    `json:"restaurant_id"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
