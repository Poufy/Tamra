package routes

import (
	"Tamra/internal/app/tamra/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type OrderRouter struct {
	orderHandler             *handlers.OrderHandler
	restaurantAuthMiddleware func(http.Handler) http.Handler
	userAuthMiddleware       func(http.Handler) http.Handler
	logger                   logrus.FieldLogger
}

func NewOrderRouter(orderHandler *handlers.OrderHandler, restaurantAuthMiddleware func(http.Handler) http.Handler, userAuthMiddleware func(http.Handler) http.Handler, logger logrus.FieldLogger) *OrderRouter {
	return &OrderRouter{orderHandler: orderHandler, userAuthMiddleware: userAuthMiddleware, restaurantAuthMiddleware: restaurantAuthMiddleware, logger: logger}
}

func (router *OrderRouter) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.With(router.restaurantAuthMiddleware).Group(func(r chi.Router) {
		r.Post("/", router.orderHandler.CreateOrder)
		r.Get("/restaurant", router.orderHandler.GetRestaurantOrders)
		r.Post("/{id}/reassign", router.orderHandler.ReassignOrder)
		r.Patch("/{orderID}/fulfill", router.orderHandler.FulfillOrder)
		r.Patch("/{orderID}/cancel", router.orderHandler.CancelOrder)
	})

	r.With(router.userAuthMiddleware).Group(func(r chi.Router) {
		r.Get("/user", router.orderHandler.GetUserOrders)
		r.Patch("/{id}/accept", router.orderHandler.AcceptOrder)
		r.Patch("/{id}/reject", router.orderHandler.RejectOrder)
	})
	// r.Patch("/{id}", router.orderHandler.UpdateOrder)

	return r
}
