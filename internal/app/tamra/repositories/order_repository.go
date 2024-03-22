package repositories

import (
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"database/sql"
	"fmt"
)

type OrderRepository interface {
	// CreateOrder creates a new order
	CreateOrder(order *models.Order) (*models.Order, error)
	// GetOrder returns an order by its ID
	GetOrder(id int, fbUID string) (*models.Order, error)
	// GetUserOrders returns a list of orders
	GetUserOrders(userID string) ([]*models.Order, error)
	// GetRestaurantOrders returns a list of orders
	GetRestaurantOrders(restaurantID string) ([]*models.Order, error)
	// UpdateOrder updates an order
	UpdateOrder(order *models.Order) (*models.Order, error)
	// UpdateUserOrderState updates the state of an order that belongs to a user
	UpdateUserOrderState(id int, fbUID string, state string) error
	// UpdateRestaurantOrderState updates the state of an order that belongs to a restaurant
	UpdateRestaurantOrderState(id int, fbUID string, state string) error
	// DeleteOrder deletes an order
	DeleteOrder(id int) error
	// Check if the restaurant is the owner of the order
	IsRestaurantOwnerOfOrder(id int, fbUID string) (bool, error)
}

type OrderRepositoryImpl struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) CreateOrder(order *models.Order) (*models.Order, error) {
	// The state is set to "PENDING" by default. That's why it's not included in the query
	const query = "INSERT INTO orders (user_id, restaurant_id, code, description, created_at, updated_at) VALUES ($1, $2, $3, $4, CLOCK_TIMESTAMP(), CLOCK_TIMESTAMP()) RETURNING id, user_id, restaurant_id, code, state, description, created_at, updated_at"
	err := r.db.QueryRow(query, order.UserID, order.RestaurantID, order.Code, order.Description).Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Code, &order.State, &order.Description, &order.CreatedAt, &order.UpdatedAt)
	return order, err
}

func (r *OrderRepositoryImpl) GetOrder(id int, fbUID string) (*models.Order, error) {
	order := &models.Order{}
	var userID sql.NullString // We use sql.NullString to handle the case where the user_id is null

	err := r.db.QueryRow("SELECT id, user_id, restaurant_id, code, state, description, created_at, updated_at FROM orders WHERE id = $1 AND restaurant_id = $2", id, fbUID).Scan(&order.ID, &userID, &order.RestaurantID, &order.Code, &order.State, &order.Description, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// If the user_id is not null, we set it in the order, otherwise we leave it as an empty string
	if userID.Valid {
		order.UserID = userID.String
	}

	return order, err
}

func (r *OrderRepositoryImpl) GetUserOrders(userID string) ([]*models.Order, error) {
	rows, err := r.db.Query("SELECT id, user_id, restaurant_id, code, state, description, created_at, updated_at FROM orders WHERE user_id = $1 AND (state = 'PENDING' OR state = 'ACCEPTED')", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*models.Order{}
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Code, &order.State, &order.Description, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepositoryImpl) GetRestaurantOrders(restaurantID string) ([]*models.Order, error) {
	rows, err := r.db.Query("SELECT id, user_id, restaurant_id, code, state, description, created_at, updated_at FROM orders WHERE restaurant_id = $1", restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*models.Order{}
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Code, &order.State, &order.Description, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepositoryImpl) UpdateOrder(order *models.Order) (*models.Order, error) {
	const query = "UPDATE orders SET user_id=$1, restaurant_id=$2, code=$3, state=$4, description=$5, updated_at=CLOCK_TIMESTAMP() WHERE id=$6 RETURNING id, user_id, restaurant_id, code, state, description, created_at, updated_at"
	err := r.db.QueryRow(query, order.UserID, order.RestaurantID, order.Code, order.State, order.Description, order.ID).Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Code, &order.State, &order.Description, &order.CreatedAt, &order.UpdatedAt)
	return order, err
}

// Updates the state of an order that belongs to a user
func (r *OrderRepositoryImpl) UpdateUserOrderState(id int, fbUID string, state string) error {
	_, err := r.db.Exec("UPDATE orders SET state = $1 WHERE id = $2 AND user_id = $3", state, id, fbUID)
	return err
}

// Updates the state of an order that belongs to a restaurant
func (r *OrderRepositoryImpl) UpdateRestaurantOrderState(id int, fbUID string, state string) error {
	_, err := r.db.Exec("UPDATE orders SET state = $1 WHERE id = $2 AND restaurant_id = $3", state, id, fbUID)
	return err
}

// Fulfill the order if the order is accepted
func (r *OrderRepositoryImpl) FulfillOrder(id int, fbUID string) error {
	var currentState string
	err := r.db.QueryRow("SELECT state FROM orders WHERE id = $1 AND restaurant_id = $2", id, fbUID).Scan(&currentState)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no order found with the given id and restaurant_id")
		}
		return err
	}

	if currentState != "ACCEPTED" {
		return utils.ErrOrderNotAccepted
	}

	_, err = r.db.Exec("UPDATE orders SET state = 'FULFILLED' WHERE id = $1 AND restaurant_id = $2", id, fbUID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryImpl) IsRestaurantOwnerOfOrder(id int, fbUID string) (bool, error) {
	var exists bool
	vertificationQuery := `
			SELECT EXISTS (
				SELECT 1
				FROM orders o
				JOIN restaurants r ON o.restaurant_id = r.id
				WHERE o.id = $1 AND r.fb_user_id = $2
		)`
	err := r.db.QueryRow(vertificationQuery, id, fbUID).Scan(&exists)
	return exists, err
}

func (r *OrderRepositoryImpl) DeleteOrder(id int) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id = $1", id)
	return err
}
