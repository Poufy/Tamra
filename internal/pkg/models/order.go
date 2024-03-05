package models

import (
	"time"
)

type Order struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id" validate:"required"`
	RestaurantID int       `json:"restaurant_id" validate:"required"`
	Code         string    `json:"code" validate:"required"`
	Description  string    `json:"description"`
	State        string    `json:"state" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
	UserID       int    `json:"user_id" validate:"required"`
	RestaurantID int    `json:"restaurant_id" validate:"required"`
	Code         string `json:"code" validate:"required,len=6"`
	Description  string `json:"description"`
}

type UpdateOrderRequest struct {
	UserID       int    `json:"user_id" validate:"required"`
	RestaurantID int    `json:"restaurant_id" validate:"required"`
	Code         string `json:"code" validate:"required"`
	Description  string `json:"description"`
	State        string `json:"state" validate:"required"`
}

type OrderResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	RestaurantID int       `json:"restaurant_id"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
