package services

import (
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/pkg/models"
	"fmt"
)

type OrderService struct {
	orderRepository repositories.OrderRepository
}

func NewOrderService(orderRepository repositories.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

func (s *OrderService) CreateOrder(order *models.Order) (*models.Order, error) {
	createdOrder, err := s.orderRepository.CreateOrder(order)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return createdOrder, nil
}

func (s *OrderService) GetOrder(id int) (*models.Order, error) {
	order, err := s.orderRepository.GetOrder(id)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (s *OrderService) GetOrders() ([]*models.Order, error) {
	orders, err := s.orderRepository.GetOrders()
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (s *OrderService) UpdateOrder(order *models.Order) (*models.Order, error) {
	updatedOrder, err := s.orderRepository.UpdateOrder(order)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return updatedOrder, nil
}

func (s *OrderService) DeleteOrder(id int) error {
	err := s.orderRepository.DeleteOrder(id)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
