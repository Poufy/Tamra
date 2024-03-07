package handlers

import (
	"Tamra/internal/app/tamra/services"
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

// CreateRestaurant godoc
// @Summary Create a new restaurant
// @Description Create a new restaurant with the given request body
// @Tags restaurants
// @Accept json
// @Produce json
// @Param request body models.CreateRestaurantRequest true "Create Restaurant Request"
// @Security jwt
// @Success 201 {object} models.Restaurant "Created Restaurant"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create restaurant"
// @Router /restaurants [post]
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
		fmt.Fprint(w, "invalid request body")
		return
	}

	// Here we would map the CreateRestaurantRequest to a Restaurant and pass it to the service
	// The reason why we map to a Restaurant is because the service should not know about the request/response models
	// It should be loosely coupled and only know about the domain models
	restaurant := utils.MapCreateRestaurantRequestToRestaurant(createRestaurantRequest)
	// Extract the user ID from the request context
	userID, ok := r.Context().Value("UserID").(string)
	if !ok {
		h.logger.Error("Context", r.Context())
		h.logger.Error("failed to get user ID from request context")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	restaurant.UserID = userID

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

// GetRestaurant godoc
// @Summary Get a restaurant
// @Description Get a restaurant by the user ID
// @Tags restaurants
// @Produce json
// @Security jwt
// @Success 200 {object} models.Restaurant "Restaurant"
// @Failure 404 {string} string "Restaurant not found"
// @Failure 500 {string} string "Failed to get restaurant"
// @Router /restaurants/me [get]
func (h *RestaurantHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the request context
	userID := r.Context().Value("UserID").(string)

	if userID == "" {
		h.logger.Error("failed to get user ID from request context")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	restaurant, err := h.restaurantService.GetRestaurant(userID)
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

// UpdateRestaurant godoc
// @Summary Update a restaurant
// @Description Update a restaurant with the given request body
// @Tags restaurants
// @Accept json
// @Produce json
// @Param request body models.UpdateRestaurantRequest true "Update Restaurant Request"
// @Security jwt
// @Success 200 {object} models.Restaurant "Updated Restaurant"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to update restaurant"
// @Router /restaurants/me [patch]
func (h *RestaurantHandler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the request context
	h.logger.Infof("Context: %+v", r.Context())
	userId, ok := r.Context().Value("UserID").(string)

	if !ok {
		h.logger.Error("Context", r.Context())
		h.logger.Error("failed to get user ID from request context")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	updateRestaurantRequest := &models.UpdateRestaurantRequest{}
	err := json.NewDecoder(r.Body).Decode(updateRestaurantRequest)
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
		fmt.Fprint(w, "invalid request body")
		return
	}

	restaurant := utils.MapUpdateRestaurantRequestToRestaurant(updateRestaurantRequest)
	restaurant.UserID = userId

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
