package repositories

import (
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"database/sql"

	"github.com/sirupsen/logrus"
)

//? Define generic error messages as the errors shouldn't be tied to the repository implementation
// const (
// 	ErrCreateUser  = "failed to create user"
// 	ErrGetUser     = "failed to get user"
// 	ErrGetUsers    = "failed to get users"
// 	ErrUpdateUser  = "failed to update user"
// 	ErrDeleteUser  = "failed to delete user"
// )

type UserRepository interface {
	// CreateUser creates a new user
	CreateUser(user *models.User) (*models.User, error)
	// GetUser returns a user by its ID
	GetUser(userId string) (*models.User, error)
	// UpdateUser updates a user
	UpdateUser(user *models.User) (*models.User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	const query = "INSERT INTO users (location, is_active, phone, radius, fcm_token, user_id ,last_order_received, created_at, updated_at) VALUES (ST_SetSRID(ST_MakePoint($1, $2), 4326), $3, $4, $5, $6, $7, CLOCK_TIMESTAMP(), CLOCK_TIMESTAMP(), CLOCK_TIMESTAMP()) RETURNING id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, user_id, last_order_received, created_at, updated_at"
	err := r.db.QueryRow(query, user.Longitude, user.Latitude, user.IsActive, user.Phone, user.Radius, user.FCMToken, user.UserID).Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.UserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepositoryImpl) GetUser(userId string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, user_id, last_order_received, created_at, updated_at FROM users WHERE user_id = $1", userId).Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.UserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
	// Return a custom error if the user is not found so that the service or handler can handle it.
	// In this case we want to return a 404 status code
	if err == sql.ErrNoRows {
		return nil, utils.ErrNotFound
	}

	logrus.Info("User: ", user)
	logrus.Info("User id: ", user.UserID)

	return user, err
}

func (r *UserRepositoryImpl) GetUsers() ([]*models.User, error) {
	rows, err := r.db.Query("SELECT id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, user_id, last_order_received, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.UserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, utils.ErrNotFound
	}

	return users, nil
}

func (r *UserRepositoryImpl) UpdateUser(user *models.User) (*models.User, error) {
	const query = "UPDATE users SET location = ST_SetSRID(ST_MakePoint($1, $2), 4326), is_active = $3, phone = $4, radius = $5, fcm_token = $6, last_order_received = $7, updated_at = CLOCK_TIMESTAMP() WHERE user_id = $8 RETURNING id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, user_id, last_order_received, created_at, updated_at"
	err := r.db.QueryRow(query, user.Longitude, user.Latitude, user.IsActive, user.Phone, user.Radius, user.FCMToken, user.LastOrderReceived, user.UserID).Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.UserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepositoryImpl) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
