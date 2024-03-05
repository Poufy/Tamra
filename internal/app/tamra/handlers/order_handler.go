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
	orderService *services.OrderService
	validator    Validator
	logger       logrus.FieldLogger
}

func NewOrderHandler(orderService *services.OrderService, validator Validator, logger logrus.FieldLogger) *OrderHandler {
	return &OrderHandler{orderService: orderService, validator: validator, logger: logger}
}

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
