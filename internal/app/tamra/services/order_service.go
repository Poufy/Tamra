package services

import (
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type OrderService interface {
	// CreateOrder creates a new order
	CreateOrder(order *models.Order) (*models.Order, error)
	// GetUserOrders returns a list of orders
	GetUserOrders(userID string) ([]*models.Order, error)
	// GetRestaurantOrders returns a list of orders
	GetRestaurantOrders(restaurantID string) ([]*models.Order, error)
	// UpdateOrder updates an order
	UpdateOrder(order *models.Order) (*models.Order, error)
	// DeleteOrder deletes an order
	DeleteOrder(id int) error
	// AcceptOrder accepts an order
	AcceptOrder(id int, fbUID string) error
	// RejectOrder rejects an order
	RejectOrder(id int, fbUID string) error
	// ReassignOrder reassigns an order
	ReassignOrder(id int, fbUID string) error
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

		return nil, fmt.Errorf("failed to get user to receive order: %w", err)
	}

	order.UserID = user.ID
	// Create the order and update the last_order_received of the user
	order, err = s.orderRepository.CreateOrder(order)
	s.logger.Infof("Order created: %v", order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Notify the user that a new order has been created.
	// Ideally this would be done in a seperate service that handles notifications
	err = NotifyUser(user.FCMToken, order)
	s.logger.Info("User notified")
	if err != nil {

		return nil, fmt.Errorf("failed to notify user: %w", err)
	}

	// Update the last_order_received date of the user in the database
	user.LastOrderReceived = order.CreatedAt
	_, err = s.userRepository.UpdateUser(user)
	if err != nil {

		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return order, nil
}

func (s *OrderServiceImpl) GetUserOrders(userID string) ([]*models.Order, error) {
	orders, err := s.orderRepository.GetUserOrders(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	// Iterate over the orders and if the time since the order was created is more than 15 minutes, we update the order state to "REJECTED"
	for _, order := range orders {
		if order.CreatedAt.Add(15 * time.Minute).Before(time.Now()) {
			s.logger.Infof("Order %d is more than 15 minutes old. Rejecting it", order.ID)
			// Update the order state to "REJECTED". We could have used UpdateUserOrderState as well. It doesn't matter.
			err = s.orderRepository.UpdateRestaurantOrderState(order.ID, order.RestaurantID, "REJECTED")
			if err != nil {
				return nil, fmt.Errorf("failed to reject order: %w", err)
			}
		}
	}

	// Get the updated list of orders
	orders, err = s.orderRepository.GetUserOrders(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	return orders, nil
}

func (s *OrderServiceImpl) GetRestaurantOrders(restaurantID string) ([]*models.Order, error) {
	orders, err := s.orderRepository.GetRestaurantOrders(restaurantID)
	if err != nil {

		return nil, fmt.Errorf("failed to get restaurant orders: %w", err)
	}
	return orders, nil
}

func (s *OrderServiceImpl) UpdateOrder(order *models.Order) (*models.Order, error) {
	updatedOrder, err := s.orderRepository.UpdateOrder(order)
	if err != nil {

		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return updatedOrder, nil
}

func (s *OrderServiceImpl) AcceptOrder(id int, fbUID string) error {
	err := s.orderRepository.UpdateUserOrderState(id, fbUID, "ACCEPTED")
	if err != nil {

		return fmt.Errorf("failed to accept order: %w", err)
	}

	return nil
}

func (s *OrderServiceImpl) RejectOrder(id int, fbUID string) error {
	err := s.orderRepository.UpdateUserOrderState(id, fbUID, "REJECTED")
	if err != nil {

		return fmt.Errorf("failed to reject order: %w", err)
	}

	return nil
}

func (s *OrderServiceImpl) ReassignOrder(id int, fbUID string) error {
	// Update the order state to "REJECTED"
	err := s.orderRepository.UpdateRestaurantOrderState(id, fbUID, "REJECTED")
	if err != nil {

		return fmt.Errorf("failed to reassign order: %w", err)
	}

	// Get the order
	order, err := s.orderRepository.GetOrder(id, fbUID)
	if err != nil {

		return fmt.Errorf("failed to get order: %w", err)
	}

	// Create the new order
	_, err = s.CreateOrder(order)
	if err != nil {

		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (s *OrderServiceImpl) DeleteOrder(id int) error {
	err := s.orderRepository.DeleteOrder(id)
	if err != nil {

		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
