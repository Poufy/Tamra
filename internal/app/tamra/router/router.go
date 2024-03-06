package router

import (
	"Tamra/internal/app/tamra/handlers"
	"net/http"

	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	userHandler       *handlers.UserHandler
	restaurantHandler *handlers.RestaurantHandler
	orderHandler      *handlers.OrderHandler
	authMiddlware     func(http.Handler) http.Handler
	logger            logrus.FieldLogger
}

func NewRouter(userHandler *handlers.UserHandler, restaurantHandler *handlers.RestaurantHandler, orderHandler *handlers.OrderHandler, authMiddleware func(http.Handler) http.Handler, logger logrus.FieldLogger) *Router {
	return &Router{userHandler: userHandler, restaurantHandler: restaurantHandler, orderHandler: orderHandler, authMiddlware: authMiddleware, logger: logger}
}

//? Is this an overkill? Should we just return a chi.Router instead of a chi.Router wrapped in a Router struct?
//? Should we just pass the handlers to TamraRouter instead of creating a new Router struct?
func (router *Router) TamraRouter() chi.Router {
	r := chi.NewRouter()
	r.Mount("/users", userRoutes(router.userHandler))
	r.Mount("/restaurants", restaurantRoutes(router.restaurantHandler, router.authMiddlware))
	r.Mount("/orders", orderRoutes(router.orderHandler))
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

	return r
}

func restaurantRoutes(restaurantHandler *handlers.RestaurantHandler, authMiddleware func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	// Use the authMiddleware for all routes in the restaurant route
	// Middleware checks if the token is valid and if it is, it will call the next handler in the chain
	// It will also append the UUID of the user to the request context so we can use it in the handler
	r.Use(authMiddleware)
	r.Post("/", restaurantHandler.CreateRestaurant)
	r.Get("/", restaurantHandler.GetRestaurants)
	r.Get("/{id}", restaurantHandler.GetRestaurant)
	r.Patch("/{id}", restaurantHandler.UpdateRestaurant)
	return r
}

func orderRoutes(orderHandler *handlers.OrderHandler) chi.Router {
	r := chi.NewRouter()
	r.Post("/", orderHandler.CreateOrder)
	r.Get("/", orderHandler.GetOrders)
	r.Get("/{id}", orderHandler.GetOrder)
	r.Patch("/{id}", orderHandler.UpdateOrder)
	r.Post("/{id}/reassign", orderHandler.ReassignOrder)
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
