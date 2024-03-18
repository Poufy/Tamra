package repositories

import (
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var userRepo UserRepository

var dbConnectionString string = "postgresql://postgres:mysecretpassword@db:5432/tamra-postgis-test?sslmode=disable"

func TestMain(m *testing.M) {

	db, err := utils.NewDB(dbConnectionString)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	userRepo = NewUserRepository(db)

	code := m.Run()

	os.Exit(code)
}

func TestUserRepository_CreateUser(t *testing.T) {
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
