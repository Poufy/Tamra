package services

import (
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/pkg/models"
	"fmt"
)

type RestaurantService struct {
	restaurantRepository repositories.RestaurantRepository
}

func NewRestaurantService(restaurantRepository repositories.RestaurantRepository) *RestaurantService {
	return &RestaurantService{restaurantRepository: restaurantRepository}
}

func (s *RestaurantService) CreateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	createdRestaurant, err := s.restaurantRepository.CreateRestaurant(restaurant)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to create restaurant: %w", err)
	}
	return createdRestaurant, nil
}

func (s *RestaurantService) GetRestaurant(id int) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepository.GetRestaurant(id)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}
	return restaurant, nil
}

func (s *RestaurantService) GetRestaurants() ([]*models.Restaurant, error) {
	restaurants, err := s.restaurantRepository.GetRestaurants()
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get restaurants: %w", err)
	}
	return restaurants, nil
}

func (s *RestaurantService) UpdateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	updatedRestaurant, err := s.restaurantRepository.UpdateRestaurant(restaurant)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to update restaurant: %w", err)
	}
	return updatedRestaurant, nil
}

func (s *RestaurantService) DeleteRestaurant(id int) error {
	err := s.restaurantRepository.DeleteRestaurant(id)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to delete restaurant: %w", err)
	}
	return nil
}
