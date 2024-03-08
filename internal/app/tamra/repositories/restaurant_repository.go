package repositories

import (
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"database/sql"
)

type RestaurantRepository interface {
	// CreateRestaurant creates a new restaurant
	CreateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error)
	// GetRestaurant returns a restaurant by its ID
	GetRestaurant(userId string) (*models.Restaurant, error)
	// UpdateRestaurant updates a restaurant
	UpdateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error)
}

type RestaurantRepositoryImpl struct {
	db *sql.DB
}

func NewRestaurantRepository(db *sql.DB) RestaurantRepository {
	return &RestaurantRepositoryImpl{db: db}
}

func (r *RestaurantRepositoryImpl) CreateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	const query = "INSERT INTO restaurants (name, location, logo_url, user_id, created_at, updated_at) VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4, $5, CLOCK_TIMESTAMP(), CLOCK_TIMESTAMP()) RETURNING id, name, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, logo_url, user_id, created_at, updated_at"
	err := r.db.QueryRow(query, restaurant.Name, restaurant.Longitude, restaurant.Latitude, restaurant.LogoURL, restaurant.UserID).Scan(&restaurant.ID, &restaurant.Name, &restaurant.Longitude, &restaurant.Latitude, &restaurant.LogoURL, &restaurant.UserID, &restaurant.CreatedAt, &restaurant.UpdatedAt)
	return restaurant, err
}

func (r *RestaurantRepositoryImpl) GetRestaurant(userId string) (*models.Restaurant, error) {
	restaurant := &models.Restaurant{}
	err := r.db.QueryRow("SELECT id, name, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, logo_url, user_id, created_at, updated_at FROM restaurants WHERE user_id = $1", userId).Scan(&restaurant.ID, &restaurant.Name, &restaurant.Longitude, &restaurant.Latitude, &restaurant.LogoURL, &restaurant.UserID, &restaurant.CreatedAt, &restaurant.UpdatedAt)
	// Return a custom error if the restaurant is not found so that the service or handler can handle it.
	// In this case we want to return a 404 status code
	if err == sql.ErrNoRows {
		return nil, utils.ErrNotFound
	}

	return restaurant, err
}

func (r *RestaurantRepositoryImpl) UpdateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	const query = "UPDATE restaurants SET name = $1, location = ST_SetSRID(ST_MakePoint($2, $3), 4326), logo_url = $4, updated_at = CLOCK_TIMESTAMP() WHERE user_id = $5 RETURNING id, name, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, logo_url, user_id, created_at, updated_at"
	err := r.db.QueryRow(query, restaurant.Name, restaurant.Longitude, restaurant.Latitude, restaurant.LogoURL, restaurant.UserID).Scan(&restaurant.ID, &restaurant.Name, &restaurant.Longitude, &restaurant.Latitude, &restaurant.LogoURL, &restaurant.UserID, &restaurant.CreatedAt, &restaurant.UpdatedAt)
	return restaurant, err
}
