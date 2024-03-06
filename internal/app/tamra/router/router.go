package router

import (
	"Tamra/internal/app/tamra/handlers"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *handlers.UserHandler, restaurantHandler *handlers.RestaurantHandler, orderHandler *handlers.OrderHandler) chi.Router {
	r := chi.NewRouter()
	r.Mount("/users", userRoutes(userHandler))
	r.Mount("/restaurants", restaurantRoutes(restaurantHandler))
	r.Mount("/orders", orderRoutes(orderHandler))
	r.Mount("/docs", docsServeRoute())
	r.Mount("/swagger", swaggerRoute())

	return r
}

func userRoutes(userHandler *handlers.UserHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", userHandler.CreateUser)
	r.Get("/", userHandler.GetUsers)
	r.Get("/{id}", userHandler.GetUser)
	r.Patch("/{id}", userHandler.UpdateUser)
	r.Delete("/{id}", userHandler.DeleteUser)

	return r
}

func restaurantRoutes(restaurantHandler *handlers.RestaurantHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", restaurantHandler.CreateRestaurant)
	r.Get("/", restaurantHandler.GetRestaurants)
	r.Get("/{id}", restaurantHandler.GetRestaurant)
	r.Patch("/{id}", restaurantHandler.UpdateRestaurant)
	r.Delete("/{id}", restaurantHandler.DeleteRestaurant)
	return r
}

func orderRoutes(orderHandler *handlers.OrderHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", orderHandler.CreateOrder)
	r.Get("/", orderHandler.GetOrders)
	r.Get("/{id}", orderHandler.GetOrder)
	r.Patch("/{id}", orderHandler.UpdateOrder)
	r.Delete("/{id}", orderHandler.DeleteOrder)
	return r
}

func docsServeRoute() chi.Router {
	r := chi.NewRouter()
	// Currently we only have the swagger documentation as static content
	// Use the relative path to the docs directory since it is not in the same directory as the main.go file
	fileServer := http.FileServer(http.Dir("../../docs"))

	// Strip the /api/v1/docs prefix from the URL before serving the files, since the resulting URL should be the file name
	r.Handle("/*", http.StripPrefix("/api/v1/docs", fileServer))
	return r
}

func swaggerRoute() chi.Router {
	r := chi.NewRouter()
	r.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/docs/swagger.json"),
	))
	return r
}
