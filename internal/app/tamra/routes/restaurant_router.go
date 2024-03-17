package routes

import (
	"Tamra/internal/app/tamra/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type RestaurantRouter struct {
	restaurantHandler       *handlers.RestaurantHandler
	restaurantAuthMiddlware func(http.Handler) http.Handler
	userAuthMiddleware      func(http.Handler) http.Handler
	logger                  logrus.FieldLogger
}

func NewRestaurantRouter(restaurantHandler *handlers.RestaurantHandler, restaurantAuthMiddlware func(http.Handler) http.Handler, userAuthMiddleware func(http.Handler) http.Handler, logger logrus.FieldLogger) *RestaurantRouter {
	return &RestaurantRouter{restaurantHandler: restaurantHandler, restaurantAuthMiddlware: restaurantAuthMiddlware, userAuthMiddleware: userAuthMiddleware, logger: logger}
}

func (router *RestaurantRouter) GetRouter() chi.Router {
	r := chi.NewRouter()
	// Use the restaurantAuthMiddlware for all routes in the restaurant route
	// Middleware checks if the token is valid and if it is, it will call the next handler in the chain
	// It will also append the UUID of the user to the request context so we can use it in the handler
	r.With(router.restaurantAuthMiddlware).Group(func(r chi.Router) {
		r.Post("/", router.restaurantHandler.CreateRestaurant)
		r.Get("/me", router.restaurantHandler.GetRestaurant)
		r.Get("/logo/uploadurl", router.restaurantHandler.GetLogoUploadURL)
		r.Patch("/me", router.restaurantHandler.UpdateRestaurant)
	})

	r.With(router.userAuthMiddleware).Group(func(r chi.Router) {
		// Users will call this route to get restaurant details of the restaurant that sent them the order
		r.Get("/{restaurantID}", router.restaurantHandler.GetRestaurantByID)
	})

	return r
}
