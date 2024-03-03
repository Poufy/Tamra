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

	// Create the user repository
	userRepository := repositories.NewUserRepository(db)

	// Create the user service
	userService := services.NewUserService(userRepository)

	// Create the user handler
	userHandler := handlers.NewUserHandler(userService, validator)

	// Use chi as the router
	r := chi.NewRouter()

	// Define the routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetUsers)
		r.Get("/{id}", userHandler.GetUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})

	// Start the server with the port from the configuration and cast the port to a string
	strPort := ":" + strconv.Itoa(config.Port)

	http.ListenAndServe(strPort, r)
}
