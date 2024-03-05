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

type RestaurantHandler struct {
	restaurantService *services.RestaurantService
	validator         Validator
	logger            logrus.FieldLogger
}

func NewRestaurantHandler(restaurantService *services.RestaurantService, validator Validator, logger logrus.FieldLogger) *RestaurantHandler {
	return &RestaurantHandler{restaurantService: restaurantService, validator: validator, logger: logger}
}

func (h *RestaurantHandler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	createRestaurantRequest := &models.CreateRestaurantRequest{}
	err := json.NewDecoder(r.Body).Decode(createRestaurantRequest)
	if err != nil {
		h.logger.WithError(err).Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	err = h.validator.Struct(createRestaurantRequest)
	if err != nil {
		h.logger.WithError(err).Error("invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Here we would map the CreateRestaurantRequest to a Restaurant and pass it to the service
	// The reason why we map to a Restaurant is because the service should not know about the request/response models
	// It should be loosely coupled and only know about the domain models
	restaurant := utils.MapCreateRestaurantRequestToRestaurant(createRestaurantRequest)

	createdRestaurant, err := h.restaurantService.CreateRestaurant(restaurant)
	if err != nil {
		h.logger.WithError(err).Error("failed to create restaurant")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to create restaurant")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdRestaurant)
}

func (h *RestaurantHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.restaurantService.GetRestaurants()
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			h.logger.WithError(err).Error("restaurants not found")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "restaurants not found")
			return
		}
		h.logger.WithError(err).Error("failed to get restaurants")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get restaurants")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurants)
}

func (h *RestaurantHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.logger.WithError(err).Error("invalid restaurant ID")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid restaurant ID")
		return
	}

	restaurant, err := h.restaurantService.GetRestaurant(id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			h.logger.WithError(err).Error("restaurant not found")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "restaurant not found")
			return
		}
		h.logger.WithError(err).Error("failed to get restaurant")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get restaurant")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurant)
}

func (h *RestaurantHandler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.logger.WithError(err).Error("invalid restaurant ID")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid restaurant ID")
		return
	}

	updateRestaurantRequest := &models.UpdateRestaurantRequest{}
	err = json.NewDecoder(r.Body).Decode(updateRestaurantRequest)
	if err != nil {
		h.logger.WithError(err).Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	err = h.validator.Struct(updateRestaurantRequest)
	if err != nil {
		h.logger.WithError(err).Error("invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	restaurant := utils.MapUpdateRestaurantRequestToRestaurant(updateRestaurantRequest)
	restaurant.ID = id

	updatedRestaurant, err := h.restaurantService.UpdateRestaurant(restaurant)
	if err != nil {
		h.logger.WithError(err).Error("failed to update restaurant")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to update restaurant")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedRestaurant)
}

func (h *RestaurantHandler) DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.logger.WithError(err).Error("invalid restaurant ID")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid restaurant ID")
		return
	}

	err = h.restaurantService.DeleteRestaurant(id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			h.logger.WithError(err).Error("restaurant not found")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "restaurant not found")
			return
		}
		h.logger.WithError(err).Error("failed to delete restaurant")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to delete restaurant")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
