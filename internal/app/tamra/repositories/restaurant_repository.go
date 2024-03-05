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
	GetRestaurant(id int) (*models.Restaurant, error)
	// GetRestaurants returns a list of restaurants
	GetRestaurants() ([]*models.Restaurant, error)
	// UpdateRestaurant updates a restaurant
	UpdateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error)
	// DeleteRestaurant deletes a restaurant
	DeleteRestaurant(id int) error
}

type RestaurantRepositoryImpl struct {
	db *sql.DB
}

func NewRestaurantRepository(db *sql.DB) RestaurantRepository {
	return &RestaurantRepositoryImpl{db: db}
}

func (r *RestaurantRepositoryImpl) CreateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	const query = "INSERT INTO restaurants (name, location, image_url, created_at, updated_at) VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4, CLOCK_TIMESTAMP(), CLOCK_TIMESTAMP()) RETURNING id, name, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, image_url, created_at, updated_at"
	err := r.db.QueryRow(query, restaurant.Name, restaurant.Longitude, restaurant.Latitude, restaurant.ImageURL).Scan(&restaurant.ID, &restaurant.Name, &restaurant.Longitude, &restaurant.Latitude, &restaurant.ImageURL, &restaurant.CreatedAt, &restaurant.UpdatedAt)
	return restaurant, err
}

func (r *RestaurantRepositoryImpl) GetRestaurant(id int) (*models.Restaurant, error) {
	restaurant := &models.Restaurant{}
	err := r.db.QueryRow("SELECT id, name, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, image_url, created_at, updated_at FROM restaurants WHERE id = $1", id).Scan(&restaurant.ID, &restaurant.Name, &restaurant.Longitude, &restaurant.Latitude, &restaurant.ImageURL, &restaurant.CreatedAt, &restaurant.UpdatedAt)
	// Return a custom error if the restaurant is not found so that the service or handler can handle it.
	// In this case we want to return a 404 status code
	if err == sql.ErrNoRows {
		return nil, utils.ErrNotFound
	}

	return restaurant, err
}

func (r *RestaurantRepositoryImpl) GetRestaurants() ([]*models.Restaurant, error) {
	rows, err := r.db.Query("SELECT id, name, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, image_url, created_at, updated_at FROM restaurants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := []*models.Restaurant{}
	for rows.Next() {
		restaurant := &models.Restaurant{}
		err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Longitude, &restaurant.Latitude, &restaurant.ImageURL, &restaurant.CreatedAt, &restaurant.UpdatedAt)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	if len(restaurants) == 0 {
		return nil, utils.ErrNotFound
	}

	return restaurants, nil
}

func (r *RestaurantRepositoryImpl) UpdateRestaurant(restaurant *models.Restaurant) (*models.Restaurant, error) {
	_, err := r.db.Exec("UPDATE restaurants SET name = $1, location = ST_SetSRID(ST_MakePoint($2, $3), 4326), image_url = $4, updated_at = CLOCK_TIMESTAMP() WHERE id = $5", restaurant.Name, restaurant.Longitude, restaurant.Latitude, restaurant.ImageURL, restaurant.ID)
	return restaurant, err
}

func (r *RestaurantRepositoryImpl) DeleteRestaurant(id int) error {
	_, err := r.db.Exec("DELETE FROM restaurants WHERE id = $1", id)
	return err
}
