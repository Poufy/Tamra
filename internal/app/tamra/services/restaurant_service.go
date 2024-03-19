package services

import (
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"fmt"

	"github.com/sirupsen/logrus"
)

type RestaurantService interface {
	CreateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error)
	GetRestaurant(fbUID string) (*models.Restaurant, error)
	GetRestaurantByID(restaurantID string) (*models.Restaurant, error)
	UpdateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error)
	GetLogoUploadURL(UID, uploadBucketName string) (string, string, error)
	DeleteRestaurant(restaurantID string) error
}

type RestaurantServiceImpl struct {
	restaurantRepository repositories.RestaurantRepository
	logger               logrus.FieldLogger
}

func NewRestaurantService(restaurantRepository repositories.RestaurantRepository, logger logrus.FieldLogger) RestaurantService {
	return &RestaurantServiceImpl{restaurantRepository: restaurantRepository, logger: logger}
}

func (s *RestaurantServiceImpl) CreateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	createdRestaurant, err := s.restaurantRepository.CreateRestaurant(restaurant)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to create restaurant: %w", err)
	}
	return createdRestaurant, nil
}

func (s *RestaurantServiceImpl) GetRestaurant(fbUID string) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepository.GetRestaurant(fbUID)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}
	return restaurant, nil
}

func (s *RestaurantServiceImpl) GetRestaurantByID(restaurantID string) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepository.GetRestaurantByID(restaurantID)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get restaurant by ID: %w", err)
	}
	return restaurant, nil
}

func (s *RestaurantServiceImpl) UpdateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	updatedRestaurant, err := s.restaurantRepository.UpdateRestaurant(restaurant)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to update restaurant: %w", err)
	}
	return updatedRestaurant, nil
}

func (s *RestaurantServiceImpl) GetLogoUploadURL(UID, uploadBucketName string) (string, string, error) {
	// UserID is the ID extracted from the JWT token created by firebase, not the user ID from the User table.
	presignedURL, storedFileURL, err := utils.GetS3PresignedURL(UID, uploadBucketName)
	if err != nil {
		return "", "", fmt.Errorf("failed to get presigned URL: %w", err)
	}
	return presignedURL, storedFileURL, nil
}

func (s *RestaurantServiceImpl) DeleteRestaurant(restaurantID string) error {
	err := s.restaurantRepository.DeleteRestaurant(restaurantID)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to delete restaurant: %w", err)
	}
	return nil
}
