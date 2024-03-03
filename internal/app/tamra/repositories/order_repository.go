package repositories

import (
	"Tamra/internal/pkg/models"
	"database/sql"
)

type OrderRepository interface {
	// CreateOrder creates a new order
	CreateOrder(order *models.Order) error
	// GetOrder returns an order by its ID
	GetOrder(id int) (*models.Order, error)
	// GetOrders returns a list of orders
	GetOrders() ([]*models.Order, error)
	// UpdateOrder updates an order
	UpdateOrder(order *models.Order) error
	// DeleteOrder deletes an order
	DeleteOrder(id int) error
}

type OrderRepositoryImpl struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) CreateOrder(order *models.Order) error {
	_, err := r.db.Exec("INSERT INTO orders (user_id, restaurant_id, code, description, state, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", order.UserID, order.RestaurantID, order.Code, order.Description, order.State, order.CreatedAt, order.UpdatedAt)
	return err
}

func (r *OrderRepositoryImpl) GetOrder(id int) (*models.Order, error) {
	order := &models.Order{}
	err := r.db.QueryRow("SELECT id, user_id, restaurant_id, code, description, state, created_at, updated_at FROM orders WHERE id = $1", id).Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Code, &order.Description, &order.State, &order.CreatedAt, &order.UpdatedAt)
	return order, err
}

func (r *OrderRepositoryImpl) GetOrders() ([]*models.Order, error) {
	rows, err := r.db.Query("SELECT id, user_id, restaurant_id, code, description, state, created_at, updated_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*models.Order{}
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.RestaurantID, &order.Code, &order.Description, &order.State, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepositoryImpl) UpdateOrder(order *models.Order) error {
	_, err := r.db.Exec("UPDATE orders SET user_id = $1, restaurant_id = $2, code = $3, description = $4, state = $5, updated_at = $6 WHERE id = $7", order.UserID, order.RestaurantID, order.Code, order.Description, order.State, order.UpdatedAt, order.ID)
	return err
}

func (r *OrderRepositoryImpl) DeleteOrder(id int) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id = $1", id)
	return err
}
