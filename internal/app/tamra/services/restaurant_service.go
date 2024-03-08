package services

import (
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
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

func (s *RestaurantService) GetRestaurant(userId string) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepository.GetRestaurant(userId)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}
	return restaurant, nil
}

func (s *RestaurantService) UpdateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	updatedRestaurant, err := s.restaurantRepository.UpdateRestaurant(restaurant)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to update restaurant: %w", err)
	}
	return updatedRestaurant, nil
}

func (s *RestaurantService) GetLogoUploadURL(UID, uploadBucketName string) (string, string, error) {
	// UserID is the ID extracted from the JWT token created by firebase, not the user ID from the User table.
	presignedURL, storedFileURL, err := utils.GetS3PresignedURL(UID, uploadBucketName)
	if err != nil {
		return "", "", fmt.Errorf("failed to get presigned URL: %w", err)
	}
	return presignedURL, storedFileURL, nil
}
