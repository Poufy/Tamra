package repositories

import (
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateUser(t *testing.T) {
	userRepo := NewUserRepository(Db)

	// Set the id, phone, and FCM tokens to a random string to avoid conflicts with other tests
	user := &models.User{
		ID:        "dasdsad",
		Longitude: 12.9715987,
		Latitude:  77.5945667,
		IsActive:  true,
		Phone:     "42423423",
		Radius:    1000,
		FCMToken:  "14213215",
	}

	createdUser, err := userRepo.CreateUser(user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, user.ID, createdUser.ID)
	assert.Equal(t, user.Longitude, createdUser.Longitude)
	assert.Equal(t, user.Latitude, createdUser.Latitude)
	assert.Equal(t, user.IsActive, createdUser.IsActive)
	assert.Equal(t, user.Phone, createdUser.Phone)
	assert.Equal(t, user.Radius, createdUser.Radius)
	assert.Equal(t, user.FCMToken, createdUser.FCMToken)
}

func TestUserRepository_GetUser(t *testing.T) {
	userRepo := NewUserRepository(Db)

	user := &models.User{
		ID:        "dasdshfgdad",
		Longitude: 12.9715987,
		Latitude:  77.5945667,
		IsActive:  true,
		Phone:     "4245123423",
		Radius:    1000,
		FCMToken:  "tewrewr",
	}

	createdUser, err := userRepo.CreateUser(user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	retrievedUser, err := userRepo.GetUser(user.ID)

	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Longitude, retrievedUser.Longitude)
	assert.Equal(t, user.Latitude, retrievedUser.Latitude)
	assert.Equal(t, user.IsActive, retrievedUser.IsActive)
	assert.Equal(t, user.Phone, retrievedUser.Phone)
	assert.Equal(t, user.Radius, retrievedUser.Radius)
	assert.Equal(t, user.FCMToken, retrievedUser.FCMToken)
}

func TestUserRepository_GetUsers(t *testing.T) {
	userRepo := NewUserRepository(Db)

	user := &models.User{
		ID:        "dasdshsfad",
		Longitude: 12.9715987,
		Latitude:  77.5945667,
		IsActive:  true,
		Phone:     "42423431223",
		Radius:    1000,
		FCMToken:  "14213215213s5",
	}

	createdUser, err := userRepo.CreateUser(user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	users, err := userRepo.GetUsers()

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.NotEmpty(t, users)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	userRepo := NewUserRepository(Db)

	user := &models.User{
		ID:        "gsdfds",
		Longitude: 12.9715987,
		Latitude:  77.5945667,
		IsActive:  true,
		Phone:     "4242376423",
		Radius:    1000,
		FCMToken:  "14213564215",
	}

	createdUser, err := userRepo.CreateUser(user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	createdUser.IsActive = false
	updatedUser, err := userRepo.UpdateUser(createdUser)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, createdUser.ID, updatedUser.ID)
	assert.Equal(t, createdUser.Longitude, updatedUser.Longitude)
	assert.Equal(t, createdUser.Latitude, updatedUser.Latitude)
	assert.Equal(t, createdUser.IsActive, updatedUser.IsActive)
	assert.Equal(t, createdUser.Phone, updatedUser.Phone)
	assert.Equal(t, createdUser.Radius, updatedUser.Radius)
	assert.Equal(t, createdUser.FCMToken, updatedUser.FCMToken)
}

func TestUserRepository_GetUserToReceiveOrder(t *testing.T) {
	userRepo := NewUserRepository(Db)
	restaurantRepo := NewRestaurantRepository(Db)

	user := &models.User{
		ID:        "gfdgdsa",
		Longitude: 19.4213234,
		Latitude:  77.5945667,
		IsActive:  true,
		Phone:     "4246123423",
		Radius:    100,
		FCMToken:  "gdfadsadcwqeqdsdf",
	}

	restaurantInReach := &models.Restaurant{
		ID:                  "uyfrghgf",
		Longitude:           19.4213234,
		Latitude:            77.5945667,
		LogoURL:             "https://www.google.com",
		Name:                "Test Restaurant1",
		PhoneNumber:         "427536423423",
		LocationDescription: "Test Location",
	}

	restaurantNotInReach := &models.Restaurant{
		ID:                  "hgfhgfh",
		Longitude:           40.9715987,
		Latitude:            50.5945667,
		LogoURL:             "https://www.google.com",
		Name:                "Test Restaurant2",
		PhoneNumber:         "3216531",
		LocationDescription: "Test Location",
	}

	createdUser, err := userRepo.CreateUser(user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	createdRestaurantInReach, err := restaurantRepo.CreateRestaurant(restaurantInReach)

	assert.NoError(t, err)
	assert.NotNil(t, createdRestaurantInReach)

	createdRestaurantNotInReach, err := restaurantRepo.CreateRestaurant(restaurantNotInReach)

	assert.NoError(t, err)
	assert.NotNil(t, createdRestaurantNotInReach)

	// Get the user that should receive the order
	userToReceiveOrder, err := userRepo.GetUserToReceiveOrder(restaurantInReach.ID)

	assert.NoError(t, err)
	assert.NotNil(t, userToReceiveOrder)
	assert.Equal(t, user.ID, userToReceiveOrder.ID)
	assert.Equal(t, user.Longitude, userToReceiveOrder.Longitude)
	assert.Equal(t, user.Latitude, userToReceiveOrder.Latitude)
	assert.Equal(t, user.IsActive, userToReceiveOrder.IsActive)
	assert.Equal(t, user.Phone, userToReceiveOrder.Phone)
	assert.Equal(t, user.Radius, userToReceiveOrder.Radius)
	assert.Equal(t, user.FCMToken, userToReceiveOrder.FCMToken)

	// Case where no user is in reach and we should get an error of type ErrNotFound
	_, err = userRepo.GetUserToReceiveOrder(restaurantNotInReach.ID)

	assert.Error(t, err)
	assert.Equal(t, err, utils.ErrNotFound)
}
