package main

import (
	"Tamra/internal/app/tamra/handlers"
	"Tamra/internal/app/tamra/middleware"
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/app/tamra/routes"
	"Tamra/internal/app/tamra/services"
	"Tamra/internal/pkg/utils"

	"Tamra/internal/pkg/utils/firebase"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

//	@title						Tamra API
//	@version					1
//	@description				This is the API for the Tamra application
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@schemes					http https
//
// this is openapi2.0 so bearer token is not supported. so we use apikey and name it jwt
//
//	@securityDefinitions.apiKey	jwt
//	@in							header
//	@name						Authorization
//	@description				Bearer token
//	@tokenUrl					http://localhost:8080/api/v1/users/login
//	@produce					json
//	@consumes					json
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

	// Get the validator
	validator := utils.NewValidator()

	// Get the logger
	logger := utils.NewLogger(config.LogLevel)

	userRepository := repositories.NewUserRepository(db)
	restaurantRepository := repositories.NewRestaurantRepository(db)
	orderRepository := repositories.NewOrderRepository(db)

	userService := services.NewUserService(userRepository, logger)
	restaurantService := services.NewRestaurantService(restaurantRepository, logger)
	orderService := services.NewOrderService(orderRepository, userRepository, logger)

	userHandler := handlers.NewUserHandler(userService, validator, logger)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, validator, logger, config)
	orderHandler := handlers.NewOrderHandler(orderService, validator, logger)

	userAuthMiddleware := middleware.UserAuthMiddleware(firebaseAuth, logger)
	restaurantAuthMiddleware := middleware.RestaurantAuthMiddleware(firebaseAuth, logger)

	logger.Info("Starting the server")
	userRouter := routes.NewUserRouter(userHandler, userAuthMiddleware, logger)
	restaurantRouter := routes.NewRestaurantRouter(restaurantHandler, restaurantAuthMiddleware, logger)
	orderRouter := routes.NewOrderRouter(orderHandler, restaurantAuthMiddleware, userAuthMiddleware, logger)
	docsRouter := routes.NewDocsRouter(logger)

	logger.Info("Creating a new chi router")
	// Create a new chi router
	r := chi.NewRouter()

	r.Mount("/users", userRouter.GetRouter())
	r.Mount("/restaurants", restaurantRouter.GetRouter())
	r.Mount("/orders", orderRouter.GetRouter())
	r.Mount("/docs", docsRouter.GetRouter())

	logger.Info("Mounting the subrouter to the parent router")
	// Mount the subrouter to the parent router
	// We create a new router for the versioned API and mount the subrouter to it
	// if we were to use the same router r like r.Mount("/api/v1", r), routers like "/users" would still be accessible
	versionedRouter := chi.NewRouter()
	versionedRouter.Mount("/api/v1", r)

	logger.Info("Starting the server")
	// Start the server with the port from the configuration and cast the port to a string
	strPort := ":" + strconv.Itoa(config.Port)

	http.ListenAndServe(strPort, versionedRouter)
}
