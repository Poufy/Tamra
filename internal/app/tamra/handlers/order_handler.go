package handlers

import (
	"Tamra/internal/app/tamra/services"
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	orderService services.OrderService
	validator    Validator
	logger       logrus.FieldLogger
}

func NewOrderHandler(orderService services.OrderService, validator Validator, logger logrus.FieldLogger) *OrderHandler {
	return &OrderHandler{orderService: orderService, validator: validator, logger: logger}
}

// CreateOrder godoc
//
//	@Summary		Create a new order
//	@Description	Create a new order with the given request body
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.CreateOrderRequest	true	"Create Order Request"
//	@Security		jwt
//	@Success		201	{object}	models.Order	"Created Order"
//	@Failure		400	{string}	string			"Invalid request body"
//	@Failure		500	{string}	string			"Failed to create order"
//	@Router			/orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	createOrderRequest := &models.CreateOrderRequest{}
	err := json.NewDecoder(r.Body).Decode(createOrderRequest)
	if err != nil {
		h.logger.WithError(err).Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	err = h.validator.Struct(createOrderRequest)
	if err != nil {
		h.logger.WithError(err).Error("invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	// Here we would map the CreateOrderRequest to an Order and pass it to the service
	// The reason why we map to an Order is because the service should not know about the request/response models
	// It should be loosely coupled and only know about the domain models
	order := utils.MapCreateOrderRequestToOrder(createOrderRequest)

	createdOrder, err := h.orderService.CreateOrder(order)
	if err != nil {
		h.logger.WithError(err).Error("failed to create order")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to create order")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}

// GetOrders godoc
//
//	@Summary		Get all orders
//	@Description	Get a list of all orders
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Security		jwt
//	@Success		200	{array}		models.Order	"List of Orders"
//	@Failure		404	{string}	string			"order not found"
//	@Failure		500	{string}	string			"failed to get orders"
//	@Router			/orders [get]
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetOrders()
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			h.logger.WithError(err).Error("order not found")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "order not found")
			return
		}
		h.logger.WithError(err).Error("failed to get orders")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get orders")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// GetOrder godoc
//
//	@Summary		Get a order
//	@Description	Get a order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Order ID"
//	@Security		jwt
//	@Success		200	{object}	models.Order	"Order"
//	@Failure		400	{string}	string			"invalid order ID"
//	@Failure		404	{string}	string			"order not found"
//	@Failure		500	{string}	string			"failed to get order"
//	@Router			/orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.WithError(err).Error("failed to parse id")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid id")
		return
	}

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			h.logger.WithError(err).Error("order not found")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "order not found")
			return
		}
		h.logger.WithError(err).Error("failed to get order")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get order")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

// UpdateOrder godoc
//
//	@Summary		Update a order
//	@Description	Update a order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int							true	"Order ID"
//	@Param			order	body	models.UpdateOrderRequest	true	"Order data to be updated"
//	@Security		jwt
//	@Success		200	{object}	models.Order	"Updated Order"
//	@Failure		400	{string}	string			"invalid order ID"
//	@Failure		400	{string}	string			"invalid request body"
//	@Failure		500	{string}	string			"failed to decode request body"
//	@Failure		500	{string}	string			"failed to update order"
//	@Router			/orders/{id} [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.WithError(err).Error("failed to parse id")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid id")
		return
	}

	updateOrderRequest := &models.UpdateOrderRequest{}
	err = json.NewDecoder(r.Body).Decode(updateOrderRequest)
	if err != nil {
		h.logger.WithError(err).Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	err = h.validator.Struct(updateOrderRequest)
	if err != nil {
		h.logger.WithError(err).Error("invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	// Here we would map the UpdateOrderRequest to an Order and pass it to the service
	// The reason why we map to an Order is because the service should not know about the request/response models
	// It should be loosely coupled and only know about the domain models
	order := utils.MapUpdateOrderRequestToOrder(updateOrderRequest)
	order.ID = id

	updatedOrder, err := h.orderService.UpdateOrder(order)
	if err != nil {
		h.logger.WithError(err).Error("failed to update order")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to update order")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedOrder)
}

// DeleteOrder godoc
//
//	@Summary		Delete a order
//	@Description	Delete a order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Order ID"
//	@Security		jwt
//	@Success		204	{string}	string	"No Content"
//	@Failure		400	{string}	string	"invalid order ID"
//	@Failure		500	{string}	string	"failed to delete order"
//	@Router			/orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.WithError(err).Error("failed to parse id")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid id")
		return
	}

	err = h.orderService.DeleteOrder(id)
	if err != nil {
		h.logger.WithError(err).Error("failed to delete order")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to delete order")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AcceptOrder godoc
//
//	@Summary		Accept a order
//	@Description	Accept a order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Order ID"
//	@Security		jwt
//	@Success		200	{string}	string	"OK"
//	@Failure		400	{string}	string	"invalid order ID"
//	@Failure		500	{string}	string	"failed to accept order"
//	@Router			/orders/{id}/accept [patch]
func (h *OrderHandler) AcceptOrder(w http.ResponseWriter, r *http.Request) {
	// We must extract the ID from the URL and the user ID from the request context to make sure the user is the owner of the order
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.WithError(err).Error("failed to parse id")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid id")
		return
	}

	// Here we would get the user ID from the request context
	firebaseUserID := r.Context().Value("UID").(string)

	err = h.orderService.AcceptOrder(id, firebaseUserID)
	if err != nil {
		h.logger.WithError(err).Error("failed to accept order")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to accept order")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// RejectOrder godoc
//
//	@Summary		Reject a order
//	@Description	Reject a order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Order ID"
//	@Security		jwt
//	@Success		200	{string}	string	"OK"
//	@Failure		400	{string}	string	"invalid order ID"
//	@Failure		500	{string}	string	"failed to reject order"
//	@Router			/orders/{id}/reject [patch]
func (h *OrderHandler) RejectOrder(w http.ResponseWriter, r *http.Request) {
	// We must extract the ID from the URL and the user ID from the request context to make sure the user is the owner of the order
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.WithError(err).Error("failed to parse id")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid id")
		return
	}

	// Here we would get the user ID from the request context
	firebaseUserID := r.Context().Value("UID").(string)

	err = h.orderService.RejectOrder(id, firebaseUserID)
	if err != nil {
		h.logger.WithError(err).Error("failed to reject order")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to reject order")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// ReassignOrder godoc
//
//	@Summary		Reassign a order
//	@Description	Reassign a order
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Order ID"
//	@Security		jwt
//	@Success		200	{string}	string	"OK"
//	@Failure		400	{string}	string	"invalid order ID"
//	@Failure		500	{string}	string	"failed to reassign order"
//	@Router			/orders/{id}/reassign [post]
func (h *OrderHandler) ReassignOrder(w http.ResponseWriter, r *http.Request) {
	// We must extract the ID from the URL and the user ID from the request context to make sure the user is the owner of the order
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.logger.WithError(err).Error("failed to parse id")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid id")
		return
	}

	firebaseUserID := r.Context().Value("UID").(string)

	err = h.orderService.ReassignOrder(id, firebaseUserID)

	if err != nil {
		h.logger.WithError(err).Error("failed to reassign order")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to reassign order")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
