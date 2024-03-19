package repositories

import (
	"Tamra/internal/pkg/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestaurantRepository_CreateRestaurant(t *testing.T) {
	restaurantRepo := NewRestaurantRepository(Db)

	restaurant := &models.Restaurant{
		ID:                  "dasdshfgdad",
		Longitude:           12.9715987,
		Latitude:            77.5945667,
		LogoURL:             "https://www.google.com",
		Name:                "Test Restaurant",
		PhoneNumber:         "4245123423",
		LocationDescription: "Test Location",
	}
	fmt.Println("restaurantRepo", restaurantRepo)
	createdRestaurant, err := restaurantRepo.CreateRestaurant(restaurant)

	assert.NoError(t, err)
	assert.NotNil(t, createdRestaurant)
	assert.Equal(t, restaurant.ID, createdRestaurant.ID)
	assert.Equal(t, restaurant.Longitude, createdRestaurant.Longitude)
	assert.Equal(t, restaurant.Latitude, createdRestaurant.Latitude)
	assert.Equal(t, restaurant.LogoURL, createdRestaurant.LogoURL)
	assert.Equal(t, restaurant.Name, createdRestaurant.Name)
	assert.Equal(t, restaurant.PhoneNumber, createdRestaurant.PhoneNumber)
	assert.Equal(t, restaurant.LocationDescription, createdRestaurant.LocationDescription)
}

func TestRestaurantRepository_GetRestaurant(t *testing.T) {
	restaurantRepo := NewRestaurantRepository(Db)

	restaurant := &models.Restaurant{
		ID:                  "dasdsfsagdad",
		Longitude:           12.9715987,
		Latitude:            77.5945667,
		LogoURL:             "https://www.googles.com",
		Name:                "Test Refsadstaurant",
		PhoneNumber:         "4245123321423",
		LocationDescription: "Test Location",
	}

	createdRestaurant, err := restaurantRepo.CreateRestaurant(restaurant)

	assert.NoError(t, err)
	assert.NotNil(t, createdRestaurant)

	retrievedRestaurant, err := restaurantRepo.GetRestaurant(restaurant.ID)

	assert.NoError(t, err)
	assert.NotNil(t, retrievedRestaurant)
	assert.Equal(t, restaurant.ID, retrievedRestaurant.ID)
	assert.Equal(t, restaurant.Longitude, retrievedRestaurant.Longitude)
	assert.Equal(t, restaurant.Latitude, retrievedRestaurant.Latitude)
	assert.Equal(t, restaurant.LogoURL, retrievedRestaurant.LogoURL)
	assert.Equal(t, restaurant.Name, retrievedRestaurant.Name)
	assert.Equal(t, restaurant.PhoneNumber, retrievedRestaurant.PhoneNumber)
	assert.Equal(t, restaurant.LocationDescription, retrievedRestaurant.LocationDescription)
}

func TestRestaurantRepository_GetRestaurantByID(t *testing.T) {
	restaurantRepo := NewRestaurantRepository(Db)

	restaurant := &models.Restaurant{
		ID:                  "dasdshgdfgdad",
		Longitude:           12.9715987,
		Latitude:            77.5945667,
		LogoURL:             "https://www.googlewe.com",
		Name:                "Test Restaurantfdgf",
		PhoneNumber:         "424512346123",
		LocationDescription: "Test Location",
	}

	createdRestaurant, err := restaurantRepo.CreateRestaurant(restaurant)

	assert.NoError(t, err)
	assert.NotNil(t, createdRestaurant)

	retrievedRestaurant, err := restaurantRepo.GetRestaurantByID(restaurant.ID)

	assert.NoError(t, err)
	assert.NotNil(t, retrievedRestaurant)
	assert.Equal(t, restaurant.ID, retrievedRestaurant.ID)
	assert.Equal(t, restaurant.Longitude, retrievedRestaurant.Longitude)
	assert.Equal(t, restaurant.Latitude, retrievedRestaurant.Latitude)
	assert.Equal(t, restaurant.LogoURL, retrievedRestaurant.LogoURL)
	assert.Equal(t, restaurant.Name, retrievedRestaurant.Name)
	assert.Equal(t, restaurant.PhoneNumber, retrievedRestaurant.PhoneNumber)
	assert.Equal(t, restaurant.LocationDescription, retrievedRestaurant.LocationDescription)
}

func TestRestaurantRepository_UpdateRestaurant(t *testing.T) {
	restaurantRepo := NewRestaurantRepository(Db)

	restaurant := &models.Restaurant{
		ID:                  "dasdfdshfgdad",
		Longitude:           12.9715987,
		Latitude:            77.5945667,
		LogoURL:             "https://www.google.com",
		Name:                "Test Restaurandasdsadt",
		PhoneNumber:         "6213213213",
		LocationDescription: "Test Location",
	}

	createdRestaurant, err := restaurantRepo.CreateRestaurant(restaurant)

	assert.NoError(t, err)
	assert.NotNil(t, createdRestaurant)

	restaurant.Name = "Updated Restaurant"
	restaurant.LocationDescription = "Updated Location"
	restaurant.PhoneNumber = "4245123423"
	restaurant.LogoURL = "https://www.google.com"
	restaurant.Longitude = 14.955987
	restaurant.Latitude = 76.5945667

	updatedRestaurant, err := restaurantRepo.UpdateRestaurant(restaurant)

	assert.NoError(t, err)
	assert.NotNil(t, updatedRestaurant)
	assert.Equal(t, restaurant.ID, updatedRestaurant.ID)
	assert.Equal(t, restaurant.Longitude, updatedRestaurant.Longitude)
	assert.Equal(t, restaurant.Latitude, updatedRestaurant.Latitude)
	assert.Equal(t, restaurant.LogoURL, updatedRestaurant.LogoURL)
	assert.Equal(t, restaurant.Name, updatedRestaurant.Name)
	assert.Equal(t, restaurant.PhoneNumber, updatedRestaurant.PhoneNumber)
	assert.Equal(t, restaurant.LocationDescription, updatedRestaurant.LocationDescription)
}
