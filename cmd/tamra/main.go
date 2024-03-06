package main

import (
	"Tamra/internal/app/tamra/handlers"
	"Tamra/internal/app/tamra/middleware"
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/app/tamra/router"
	"Tamra/internal/app/tamra/services"
	"Tamra/internal/pkg/utils"
	"Tamra/internal/pkg/utils/firebase"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
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

	fmt.Println("Configuration values in main", config)
	// Connect to the database
	db, err := utils.NewDB(config.DBConn)

	if err != nil {
		panic(err)
	}

	fmt.Println("Configuration values before initializing firebase auth", config.FirebaseConfig)
	firebaseAuth, err := firebase.NewFirebaseAuth(config.FirebaseConfig)

	if err != nil {
		fmt.Println("Error initializing firebase auth", err)
		panic(err)
	}

	logrus.Infof("Firebase auth: %v", firebaseAuth)

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

	authMiddleware := middleware.AuthMiddleware(firebaseAuth, logger)

	// Create a parent router
	r := chi.NewRouter()

	// Create a new router instance and inject the handlers, middleware and logger
	routerInstance := router.NewRouter(userHandler, restaurantHandler, orderHandler, authMiddleware, logger)

	tamraRouter := routerInstance.TamraRouter()

	// Mount the subrouter to the parent router
	r.Mount("/api/v1", tamraRouter)

	// Start the server with the port from the configuration and cast the port to a string
	strPort := ":" + strconv.Itoa(config.Port)

	http.ListenAndServe(strPort, r)
}
