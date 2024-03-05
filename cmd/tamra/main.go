package main

import (
	"Tamra/internal/app/tamra/handlers"
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/app/tamra/services"
	"Tamra/internal/pkg/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Read the configuration
	config := utils.GetConfig()

	// Connect to the database
	fmt.Println("Connecting to the database", config.DBConn)
	db, err := utils.NewDB(config.DBConn)

	if err != nil {
		panic(err)
	}

	// Get the validator
	validator := utils.NewValidator()

	// Get the logger
	logger := utils.NewLogger(config.LogLevel)

	userRepository := repositories.NewUserRepository(db)
	restaurantRepository := repositories.NewRestaurantRepository(db)
	orderRepository := repositories.NewOrderRepository(db)

	userService := services.NewUserService(userRepository)
	restaurantService := services.NewRestaurantService(restaurantRepository)
	orderService := services.NewOrderService(orderRepository)

	userHandler := handlers.NewUserHandler(userService, validator, logger)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, validator, logger)
	orderHandler := handlers.NewOrderHandler(orderService, validator, logger)

	// Use chi as the router
	r := chi.NewRouter()

	// Define the routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetUsers)
		r.Get("/{id}", userHandler.GetUser)
		r.Patch("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	r.Route("/restaurants", func(r chi.Router) {
		r.Post("/", restaurantHandler.CreateRestaurant)
		r.Get("/", restaurantHandler.GetRestaurants)
		r.Get("/{id}", restaurantHandler.GetRestaurant)
		r.Patch("/{id}", restaurantHandler.UpdateRestaurant)
		r.Delete("/{id}", restaurantHandler.DeleteRestaurant)
	})

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.CreateOrder)
		r.Get("/", orderHandler.GetOrders)
		r.Get("/{id}", orderHandler.GetOrder)
		r.Patch("/{id}", orderHandler.UpdateOrder)
		r.Delete("/{id}", orderHandler.DeleteOrder)
	})

	// Start the server with the port from the configuration and cast the port to a string
	strPort := ":" + strconv.Itoa(config.Port)

	http.ListenAndServe(strPort, r)
}
