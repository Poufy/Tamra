package repositories

import (
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"database/sql"
)

type OrderRepository interface {
	// CreateOrder creates a new order
	CreateOrder(order *models.Order) (*models.Order, error)
	// GetOrder returns an order by its ID
	GetOrder(id int) (*models.Order, error)
	// GetOrders returns a list of orders
	GetOrders() ([]*models.Order, error)
	// UpdateOrder updates an order
	UpdateOrder(order *models.Order) (*models.Order, error)
	// DeleteOrder deletes an order
	DeleteOrder(id int) error
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

func (r *OrderRepositoryImpl) GetOrder(id int) (*models.Order, error) {
	order := &models.Order{}
	err := r.db.QueryRow("SELECT id, user_id, restaurant_id, code, state, description, created_at, updated_at FROM orders WHERE id = $1", id).Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Code, &order.State, &order.Description, &order.CreatedAt, &order.UpdatedAt)
	return order, err
}

func (r *OrderRepositoryImpl) GetOrders() ([]*models.Order, error) {
	rows, err := r.db.Query("SELECT id, user_id, restaurant_id, code, state, description, created_at, updated_at FROM orders")
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

	if len(orders) == 0 {
		return nil, utils.ErrNotFound
	}

	return orders, nil
}

func (r *OrderRepositoryImpl) UpdateOrder(order *models.Order) (*models.Order, error) {
	const query = "UPDATE orders SET user_id=$1, restaurant_id=$2, code=$3, state=$4, description=$5, updated_at=CLOCK_TIMESTAMP() WHERE id=$6 RETURNING id, user_id, restaurant_id, code, state, description, created_at, updated_at"
	err := r.db.QueryRow(query, order.UserID, order.RestaurantID, order.Code, order.State, order.Description, order.ID).Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Code, &order.State, &order.Description, &order.CreatedAt, &order.UpdatedAt)
	return order, err
}

func (r *OrderRepositoryImpl) DeleteOrder(id int) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id = $1", id)
	return err
}
