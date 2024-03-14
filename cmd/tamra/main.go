package main

import (
	"Tamra/internal/app/tamra/handlers"
	"Tamra/internal/app/tamra/middleware"
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/app/tamra/routes"
	"Tamra/internal/app/tamra/services"
	"Tamra/internal/pkg/utils"
	"net/http"
	"os"
	"strconv"

	"Tamra/internal/pkg/utils/firebase"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

//	@title						Tamra API
//	@version					1
//	@description				This is the API for the Tamra application
//	@host						tamra.gulbababaklava.com
//	@BasePath					/api/v1
//	@schemes					https
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
	// Load the environment variables from the .env file
	if os.Getenv("LAMBDA_TASK_ROOT") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
		}
	}

	// Read the configuration
	config := utils.GetConfig()

	fmt.Println("Configuration values in main", config)
	// Connect to the database
	db, err := utils.NewDB(config.DBConn)

	if err != nil {
		panic(err)
	}

	// Get the logger
	logger := utils.NewLogger(config.LogLevel)

	logger.Info("Configuration values before initializing firebase auth", config.FirebaseConfigJSON)
	firebaseApp := firebase.NewFirebaseApp(config.FirebaseConfigJSON)
	if err != nil {
		logrus.Panic("Failed to initialize firebase auth: ", err)
	}

	firebaseAuthClient, err := firebaseApp.FetchFirebaseAuthClient()
	if err != nil {
		logrus.Panic("Failed to initialize firebase auth: ", err)
	}

	firebaseMessagingClient, err := firebaseApp.FetchFirebaseMessagingClient()
	if err != nil {
		logrus.Panic("Failed to initialize firebase messaging client: ", err)
	}

	// Get the validator
	validator := utils.NewValidator()

	userRepository := repositories.NewUserRepository(db)
	restaurantRepository := repositories.NewRestaurantRepository(db)
	orderRepository := repositories.NewOrderRepository(db)

	notificationService := services.NewNotificationService(logger, firebaseMessagingClient)
	userService := services.NewUserService(userRepository, logger)
	restaurantService := services.NewRestaurantService(restaurantRepository, logger)
	orderService := services.NewOrderService(orderRepository, userRepository, notificationService, logger)

	userHandler := handlers.NewUserHandler(userService, validator, logger)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService, validator, logger, config)
	orderHandler := handlers.NewOrderHandler(orderService, validator, logger)

	userAuthMiddleware := middleware.UserAuthMiddleware(firebaseAuthClient, logger)
	restaurantAuthMiddleware := middleware.RestaurantAuthMiddleware(firebaseAuthClient, logger)

	logger.Info("Starting the server")
	userRouter := routes.NewUserRouter(userHandler, userAuthMiddleware, logger)
	restaurantRouter := routes.NewRestaurantRouter(restaurantHandler, restaurantAuthMiddleware, logger)
	orderRouter := routes.NewOrderRouter(orderHandler, restaurantAuthMiddleware, userAuthMiddleware, logger)
	docsRouter := routes.NewDocsRouter(logger)

	logger.Info("Creating a new chi router")
	// Create a new chi router
	r := chi.NewRouter()

	// Log all requests
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestIDMiddleware)

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

	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		// If we are running on AWS Lambda, we use the chiadapter to convert the chi router to a lambda handler
		chiLambda := chiadapter.New(versionedRouter)
		lambda.Start(chiLambda.Proxy)
	} else {
		// If we are not running on AWS Lambda, we start the server using the port from the configuration
		strPort := ":" + strconv.Itoa(config.Port)
		http.ListenAndServe(strPort, versionedRouter)
	}
}
