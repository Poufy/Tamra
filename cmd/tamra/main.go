package main

import (
	"Tamra/internal/app/tamra/handlers"
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/app/tamra/router"
	"Tamra/internal/app/tamra/services"
	"Tamra/internal/pkg/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

// @title Tamra API
// @version 1
// @description This is the API for the Tamra application
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @produce json
// @consumes json
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

	// Create a parent router
	r := chi.NewRouter()

	// Create a new router and pass the handlers to it
	apiRouter := router.NewRouter(userHandler, restaurantHandler, orderHandler)

	// Mount the subrouter to the parent router
	r.Mount("/api/v1", apiRouter)

	// Start the server with the port from the configuration and cast the port to a string
	strPort := ":" + strconv.Itoa(config.Port)

	http.ListenAndServe(strPort, r)
}
