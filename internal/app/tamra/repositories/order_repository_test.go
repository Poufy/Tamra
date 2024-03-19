package repositories

import (
	"Tamra/internal/pkg/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_CreateOrder(t *testing.T) {
	orderRepo := NewOrderRepository(Db)

	order := &models.Order{
		UserID:       "user1",       // User from the seed
		RestaurantID: "restaurant1", // Restaurant from the seed
		Code:         "123",
		Description:  "Test Order",
	}

	createdOrder, err := orderRepo.CreateOrder(order)
	fmt.Printf("createdOrder: %v", createdOrder)

	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.Equal(t, order.UserID, createdOrder.UserID)
	assert.Equal(t, order.RestaurantID, createdOrder.RestaurantID)
	assert.Equal(t, order.Code, createdOrder.Code)
	assert.Equal(t, order.Description, createdOrder.Description)
	assert.Equal(t, "PENDING", createdOrder.State)
}

func TestOrderRepository_GetOrder(t *testing.T) {
	orderRepo := NewOrderRepository(Db)

	order := &models.Order{
		UserID:       "user1",       // User from the seed
		RestaurantID: "restaurant1", // Restaurant from the seed
		Code:         "12345",
		Description:  "Test Order",
	}

	createdOrder, err := orderRepo.CreateOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)

	retrievedOrder, err := orderRepo.GetOrder(order.ID, order.RestaurantID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrder)
	assert.Equal(t, order.ID, retrievedOrder.ID)
	assert.Equal(t, order.UserID, retrievedOrder.UserID)
	assert.Equal(t, order.RestaurantID, retrievedOrder.RestaurantID)
	assert.Equal(t, order.Code, retrievedOrder.Code)
	assert.Equal(t, order.Description, retrievedOrder.Description)
}

func TestOrderRepository_GetUserOrders(t *testing.T) {
	orderRepo := NewOrderRepository(Db)

	// Create a new order
	order := &models.Order{
		UserID:       "user1",       // User from the seed
		RestaurantID: "restaurant1", // Restaurant from the seed
		Code:         "312312",
		Description:  "Test Order",
	}
	createdOrder, err := orderRepo.CreateOrder(order)

	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)

	// Get the user's orders
	orders, err := orderRepo.GetUserOrders(order.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.NotEmpty(t, orders)
}

func TestOrderRepository_GetRestaurantOrders(t *testing.T) {
	orderRepo := NewOrderRepository(Db)

	// Create a new order
	order := &models.Order{
		UserID:       "user1",       // User from the seed
		RestaurantID: "restaurant1", // Restaurant from the seed
		Code:         "31231dsa2",
		Description:  "Test Order",
	}
	createdOrder, err := orderRepo.CreateOrder(order)

	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)

	// Get the restaurant's orders
	orders, err := orderRepo.GetRestaurantOrders(order.RestaurantID)

	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.NotEmpty(t, orders)
}

func TestOrderRepository_UpdateOrder(t *testing.T) {
	orderRepo := NewOrderRepository(Db)

	// Create a new order
	order := &models.Order{
		UserID:       "user1",       // User from the seed
		RestaurantID: "restaurant1", // Restaurant from the seed
		Code:         "312314521",
		Description:  "Test Order",
	}

	createdOrder, err := orderRepo.CreateOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)

	// Update the order
	createdOrder.Description = "Updated Order"
	updatedOrder, err := orderRepo.UpdateOrder(createdOrder)
	assert.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, createdOrder.ID, updatedOrder.ID)
	assert.Equal(t, createdOrder.UserID, updatedOrder.UserID)
	assert.Equal(t, createdOrder.RestaurantID, updatedOrder.RestaurantID)
	assert.Equal(t, createdOrder.Code, updatedOrder.Code)
	assert.Equal(t, createdOrder.Description, updatedOrder.Description)
	assert.Equal(t, createdOrder.State, updatedOrder.State)
}

func TestOrderRepository_UpdateUserOrderState(t *testing.T) {
	orderRepo := NewOrderRepository(Db)

	// Create a new order
	order := &models.Order{
		UserID:       "user2",       // User from the seed
		RestaurantID: "restaurant2", // Restaurant from the seed
		Code:         "61232",
		Description:  "Test Order",
	}

	createdOrder, err := orderRepo.CreateOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)

	fmt.Printf("Created Order: %v", createdOrder)
	// Update the order state
	err = orderRepo.UpdateUserOrderState(createdOrder.ID, createdOrder.UserID, "ACCEPTED")
	assert.NoError(t, err)

	// Get the updated order
	retrievedOrder, err := orderRepo.GetOrder(createdOrder.ID, createdOrder.RestaurantID)
	fmt.Printf("Retrieved Order: %v", retrievedOrder)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrder)
	assert.Equal(t, "ACCEPTED", retrievedOrder.State)
}

func TestOrderRepository_UpdateRestaurantOrderState(t *testing.T) {
	orderRepo := NewOrderRepository(Db)

	// Create a new order
	order := &models.Order{
		UserID:       "user1",       // User from the seed
		RestaurantID: "restaurant1", // Restaurant from the seed
		Code:         "612532",
		Description:  "Test Order",
	}

	createdOrder, err := orderRepo.CreateOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)

	// Update the order state
	err = orderRepo.UpdateRestaurantOrderState(createdOrder.ID, createdOrder.RestaurantID, "ACCEPTED")
	assert.NoError(t, err)

	// Get the updated order
	retrievedOrder, err := orderRepo.GetOrder(createdOrder.ID, createdOrder.RestaurantID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrder)
	assert.Equal(t, "ACCEPTED", retrievedOrder.State)
}

// func TestOrderRepository_DeleteOrder(t *testing.T) {
// 	orderRepo := NewOrderRepository(Db)

// 	// Create a new order
// 	order := &models.Order{
// 		UserID:       "user1",       // User from the seed
// 		RestaurantID: "restaurant1", // Restaurant from the seed
// 		Code:         "6125213",
// 		Description:  "Test Order",
// 	}

// 	createdOrder, err := orderRepo.CreateOrder(order)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, createdOrder)

// 	// Delete the order
// 	err = orderRepo.DeleteOrder(createdOrder.ID)
// 	assert.NoError(t, err)

// 	// Get the deleted order
// 	retrievedOrder, err := orderRepo.GetOrder(createdOrder.ID, createdOrder.RestaurantID)
// 	assert.Error(t, err)
// 	assert.Nil(t, retrievedOrder)
// }
