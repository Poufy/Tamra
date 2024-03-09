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
	// Retrieve the user that last received an order
	GetUserToReceiveOrder() (*models.User, error)
	// GetUsers returns a list of users
	GetUsers() ([]*models.User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	const query = "INSERT INTO users (location, is_active, phone, radius, fcm_token, fb_user_id ,last_order_received, created_at, updated_at) VALUES (ST_SetSRID(ST_MakePoint($1, $2), 4326), $3, $4, $5, $6, $7, CLOCK_TIMESTAMP(), CLOCK_TIMESTAMP(), CLOCK_TIMESTAMP()) RETURNING id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, fb_user_id, last_order_received, created_at, updated_at"
	err := r.db.QueryRow(query, user.Longitude, user.Latitude, user.IsActive, user.Phone, user.Radius, user.FCMToken, user.FBUserID).Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.FBUserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepositoryImpl) GetUser(userId string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow("SELECT id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, fb_user_id, last_order_received, created_at, updated_at FROM users WHERE fb_user_id = $1", userId).Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.FBUserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
	// Return a custom error if the user is not found so that the service or handler can handle it.
	// In this case we want to return a 404 status code
	if err == sql.ErrNoRows {
		return nil, utils.ErrNotFound
	}

	logrus.Info("User: ", user)
	logrus.Info("User id: ", user.FBUserID)

	return user, err
}

func (r *UserRepositoryImpl) GetUsers() ([]*models.User, error) {
	rows, err := r.db.Query("SELECT id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, fb_user_id, last_order_received, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.FBUserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
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
	const query = "UPDATE users SET location = ST_SetSRID(ST_MakePoint($1, $2), 4326), is_active = $3, phone = $4, radius = $5, fcm_token = $6, last_order_received = $7, updated_at = CLOCK_TIMESTAMP() WHERE fb_user_id = $8 RETURNING id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, fb_user_id, last_order_received, created_at, updated_at"
	err := r.db.QueryRow(query, user.Longitude, user.Latitude, user.IsActive, user.Phone, user.Radius, user.FCMToken, user.LastOrderReceived, user.FBUserID).Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.FBUserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

// GetUserToReceiveOrder retrieves the user that last received an order.
// Newly created users will be last in line to receive an order as
// last_order_recieved is set to the current time when the user is created
func (r *UserRepositoryImpl) GetUserToReceiveOrder() (*models.User, error) {
	const query = `
	SELECT id, ST_X(location::geometry) as longitude, ST_Y(location::geometry) as latitude, is_active, phone, radius, fcm_token, fb_user_id, last_order_received, created_at, updated_at 
	FROM users 
	WHERE is_active = true 
	ORDER BY last_order_received 
	ASC LIMIT 1
	`
	user := &models.User{}
	err := r.db.QueryRow(query).Scan(&user.ID, &user.Longitude, &user.Latitude, &user.IsActive, &user.Phone, &user.Radius, &user.FCMToken, &user.FBUserID, &user.LastOrderReceived, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepositoryImpl) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}