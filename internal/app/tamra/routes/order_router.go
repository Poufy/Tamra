package routes

import (
	"Tamra/internal/app/tamra/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type OrderRouter struct {
	orderHandler  *handlers.OrderHandler
	authMiddlware func(http.Handler) http.Handler
	logger        logrus.FieldLogger
}

func NewOrderRouter(orderHandler *handlers.OrderHandler, authMiddleware func(http.Handler) http.Handler, logger logrus.FieldLogger) *OrderRouter {
	return &OrderRouter{orderHandler: orderHandler, authMiddlware: authMiddleware, logger: logger}
}

func (router *OrderRouter) GetRouter() chi.Router {
	r := chi.NewRouter()
	// Use the authMiddleware for all routes in the order route
	// Middleware checks if the token is valid and if it is, it will call the next handler in the chain
	// It will also append the UUID of the user to the request context so we can use it in the handler
	r.Use(router.authMiddlware)
	// Restaurant only
	r.Post("/", router.orderHandler.CreateOrder)
	// Should return all orders for the restaurant
	r.Get("/", router.orderHandler.GetOrders)
	r.Get("/{id}", router.orderHandler.GetOrder)
	r.Patch("/{id}", router.orderHandler.UpdateOrder)
	r.Patch("/{id}/accept", router.orderHandler.AcceptOrder)
	r.Patch("/{id}/reject", router.orderHandler.RejectOrder)
	r.Post("/{id}/reassign", router.orderHandler.ReassignOrder)
	return r
}
