package utils

// In this file we have mappers that map between request/response and domain models.

import (
	"Tamra/internal/pkg/models"
)

// TODO: is there a better way of mapping these? Maybe use a library like mapstruct?
// MapCreateUserRequestToUser maps a CreateUserRequest to a User.
func MapCreateUserRequestToUser(req *models.CreateUserRequest) *models.User {
	return &models.User{
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
		Phone:     req.Phone,
		Radius:    req.Radius,
		FCMToken:  req.FCMToken,
		IsActive:  *req.IsActive,
	}
}

// MapUserToUserResponse maps a User to a UserResponse.
func MapUserToUserResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:                user.ID,
		Longitude:         user.Longitude,
		Latitude:          user.Latitude,
		IsActive:          user.IsActive,
		Phone:             user.Phone,
		Radius:            user.Radius,
		LastOrderReceived: user.LastOrderReceived,
	}
}

func MapUsersToUserResponses(users []*models.User) []*models.UserResponse {
	// Pre-allocate the array to the correct length to avoid unnecessary allocations when appending
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = MapUserToUserResponse(user)
	}
	return userResponses
}

func MapUpdateUserRequestToUser(req *models.UpdateUserRequest) *models.User {
	return &models.User{
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
		IsActive:  *req.IsActive,
		Phone:     req.Phone,
		Radius:    req.Radius,
		FCMToken:  req.FCMToken,
	}
}

// MapCreateRestaurantRequestToRestaurant maps a CreateRestaurantRequest to a Restaurant.
func MapCreateRestaurantRequestToRestaurant(req *models.CreateRestaurantRequest) *models.Restaurant {
	return &models.Restaurant{
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
		LogoURL:   req.LogoURL,
		Name:      req.Name,
	}
}

// MapRestaurantToRestaurantResponse maps a Restaurant to a RestaurantResponse.
func MapRestaurantToRestaurantResponse(restaurant *models.Restaurant) *models.RestaurantResponse {
	return &models.RestaurantResponse{
		ID:        restaurant.ID,
		Longitude: restaurant.Longitude,
		Latitude:  restaurant.Latitude,
		LogoURL:   restaurant.LogoURL,
		Name:      restaurant.Name,
		CreatedAt: restaurant.CreatedAt,
		UpdatedAt: restaurant.UpdatedAt,
	}
}

func MapRestaurantsToRestaurantResponses(restaurants []*models.Restaurant) []*models.RestaurantResponse {
	// Pre-allocate the array to the correct length to avoid unnecessary allocations when appending
	restaurantResponses := make([]*models.RestaurantResponse, len(restaurants))
	for i, restaurant := range restaurants {
		restaurantResponses[i] = MapRestaurantToRestaurantResponse(restaurant)
	}
	return restaurantResponses
}

func MapUpdateRestaurantRequestToRestaurant(req *models.UpdateRestaurantRequest) *models.Restaurant {
	return &models.Restaurant{
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
		LogoURL:   req.LogoURL,
		Name:      req.Name,
	}
}

// MapCreateOrderRequestToOrder maps a CreateOrderRequest to a Order.
func MapCreateOrderRequestToOrder(req *models.CreateOrderRequest) *models.Order {
	return &models.Order{
		Description: req.Description,
	}
}

// MapOrderToOrderResponse maps a Order to a OrderResponse.
func MapOrderToOrderResponse(order *models.Order) *models.OrderResponse {
	return &models.OrderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		RestaurantID: order.RestaurantID,
		Code:         order.Code,
		Description:  order.Description,
		State:        order.State,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}
}

func MapOrdersToOrderResponses(orders []*models.Order) []*models.OrderResponse {
	// Pre-allocate the array to the correct length to avoid unnecessary allocations when appending
	orderResponses := make([]*models.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = MapOrderToOrderResponse(order)
	}
	return orderResponses
}

func MapUpdateOrderRequestToOrder(req *models.UpdateOrderRequest) *models.Order {
	return &models.Order{
		UserID:       req.UserID,
		RestaurantID: req.RestaurantID,
		Code:         req.Code,
		Description:  req.Description,
		State:        req.State,
	}
}
