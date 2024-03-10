package services

import (
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"fmt"

	"github.com/sirupsen/logrus"
)

type OrderService interface {
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
	// AcceptOrder accepts an order
	AcceptOrder(id int, fbUserID string) error
	// RejectOrder rejects an order
	RejectOrder(id int, fbUserID string) error
	// ReassignOrder reassigns an order
	ReassignOrder(id int, fbUserID string) error
}

type OrderServiceImpl struct {
	orderRepository repositories.OrderRepository
	userRepository  repositories.UserRepository
	logger          logrus.FieldLogger
}

// We return an implementation of the OrderService interface. This is so that we can easily swap out the implementation or mock it in tests.
func NewOrderService(orderRepository repositories.OrderRepository, userRepository repositories.UserRepository, logger logrus.FieldLogger) OrderService {
	return &OrderServiceImpl{orderRepository: orderRepository, userRepository: userRepository, logger: logger}
}

// We first generate a 6 digit random number as the code for the order
// We then find which user to send to based on the last_order_received of the user
// we then create the order and update the last_order_received of the user
// we then notify the user that a new order has been created
// we then return the created order
func (s *OrderServiceImpl) CreateOrder(order *models.Order) (*models.Order, error) {
	// Generate a 6 digit random number as the code for the order
	order.Code = utils.GenerateCode()

	// Find which user to send to based on the last_order_received of the user
	user, err := s.userRepository.GetUserToReceiveOrder(order.RestaurantID)
	s.logger.Infof("User to receive order: %v", user)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get user to receive order: %w", err)
	}

	order.UserID = user.ID
	// Create the order and update the last_order_received of the user
	order, err = s.orderRepository.CreateOrder(order)
	s.logger.Infof("Order created: %v", order)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Notify the user that a new order has been created.
	// Ideally this would be done in a seperate service that handles notifications
	err = NotifyUser(user.FCMToken, order)
	s.logger.Info("User notified")
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to notify user: %w", err)
	}

	// Update the last_order_received date of the user in the database
	user.LastOrderReceived = order.CreatedAt
	_, err = s.userRepository.UpdateUser(user)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return order, nil
}

func (s *OrderServiceImpl) GetOrder(id int) (*models.Order, error) {
	order, err := s.orderRepository.GetOrder(id)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (s *OrderServiceImpl) GetOrders() ([]*models.Order, error) {
	orders, err := s.orderRepository.GetOrders()
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (s *OrderServiceImpl) UpdateOrder(order *models.Order) (*models.Order, error) {
	updatedOrder, err := s.orderRepository.UpdateOrder(order)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return updatedOrder, nil
}

func (s *OrderServiceImpl) AcceptOrder(id int, fbUserID string) error {
	// check if the user is the owner of the order
	isOwner, err := s.orderRepository.IsUserOwnerOfOrder(id, fbUserID)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to check if user is owner of order: %w", err)
	}

	if !isOwner {
		return fmt.Errorf("user is not the owner of the order")
	}

	err = s.orderRepository.UpdateOrderState(id, "ACCEPTED")
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to accept order: %w", err)
	}

	return nil
}

func (s *OrderServiceImpl) RejectOrder(id int, fbUserID string) error {
	// check if the user is the owner of the order
	isOwner, err := s.orderRepository.IsUserOwnerOfOrder(id, fbUserID)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to check if user is owner of order: %w", err)
	}

	if !isOwner {
		return fmt.Errorf("user is not the owner of the order")
	}

	err = s.orderRepository.UpdateOrderState(id, "REJECTED")
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to reject order: %w", err)
	}

	return nil
}

func (s *OrderServiceImpl) ReassignOrder(id int, fbUserID string) error {
	//! TODO: The updating operations should be done in a transaction
	// check if the restaurant is the owner of the order
	isOwner, err := s.orderRepository.IsRestaurantOwnerOfOrder(id, fbUserID)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to check if restaurant is owner of order: %w", err)
	}

	if !isOwner {
		return fmt.Errorf("restaurant is not the owner of the order")
	}

	// Update the order state to "REJECTED"
	err = s.orderRepository.UpdateOrderState(id, "REJECTED")
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to reassign order: %w", err)
	}

	// Get the order
	order, err := s.orderRepository.GetOrder(id)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Create the new order
	_, err = s.CreateOrder(order)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (s *OrderServiceImpl) DeleteOrder(id int) error {
	err := s.orderRepository.DeleteOrder(id)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
