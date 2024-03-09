package routes

import (
	"Tamra/internal/app/tamra/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type RestaurantRouter struct {
	restaurantHandler *handlers.RestaurantHandler
	authMiddlware     func(http.Handler) http.Handler
	logger            logrus.FieldLogger
}

func NewRestaurantRouter(restaurantHandler *handlers.RestaurantHandler, authMiddleware func(http.Handler) http.Handler, logger logrus.FieldLogger) *RestaurantRouter {
	return &RestaurantRouter{restaurantHandler: restaurantHandler, authMiddlware: authMiddleware, logger: logger}
}

func (router *RestaurantRouter) GetRouter() chi.Router {
	r := chi.NewRouter()
	// Use the authMiddleware for all routes in the restaurant route
	// Middleware checks if the token is valid and if it is, it will call the next handler in the chain
	// It will also append the UUID of the user to the request context so we can use it in the handler
	r.Use(router.authMiddlware)
	r.Post("/", router.restaurantHandler.CreateRestaurant)
	r.Get("/me", router.restaurantHandler.GetRestaurant)
	r.Get("/logo/uploadurl", router.restaurantHandler.GetLogoUploadURL)
	r.Patch("/me", router.restaurantHandler.UpdateRestaurant)
	return r
}