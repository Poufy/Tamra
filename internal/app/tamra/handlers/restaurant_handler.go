package handlers

import (
	"Tamra/internal/app/tamra/services"
	"Tamra/internal/pkg/models"
	"Tamra/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type RestaurantHandler struct {
	restaurantService services.RestaurantService
	validator         Validator
	logger            logrus.FieldLogger
	config            utils.Config
}

func NewRestaurantHandler(restaurantService services.RestaurantService, validator Validator, logger logrus.FieldLogger, config utils.Config) *RestaurantHandler {
	return &RestaurantHandler{restaurantService: restaurantService, validator: validator, logger: logger, config: config}
}

// CreateRestaurant godoc
//
//	@Summary		Create a new restaurant
//	@Description	Create a new restaurant with the given request body
//	@Tags			restaurants
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.CreateRestaurantRequest	true	"Create Restaurant Request"
//	@Security		jwt
//	@Success		201	{object}	models.Restaurant	"Created Restaurant"
//	@Failure		400	{string}	string				"Invalid request body"
//	@Failure		500	{string}	string				"Failed to create restaurant"
//	@Router			/restaurants [post]
func (h *RestaurantHandler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof("Request ID %s: Received request to create restaurant.", r.Context().Value(chimiddleware.RequestIDKey))
	createRestaurantRequest := &models.CreateRestaurantRequest{}
	err := json.NewDecoder(r.Body).Decode(createRestaurantRequest)
	if err != nil {
		h.logger.WithError(err).Errorf("Request ID %s: Failed to decode request body", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	err = h.validator.Struct(createRestaurantRequest)
	if err != nil {
		h.logger.WithError(err).Errorf("Request ID %s: Invalid request body", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	// Here we would map the CreateRestaurantRequest to a Restaurant and pass it to the service
	// The reason why we map to a Restaurant is because the service should not know about the request/response models
	// It should be loosely coupled and only know about the domain models
	restaurant := utils.MapCreateRestaurantRequestToRestaurant(createRestaurantRequest)
	// Extract the user ID from the request context
	firebaseUserID, ok := r.Context().Value("UID").(string)
	if !ok {
		h.logger.Errorf("Request ID %s: Failed to get user ID from request context", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	restaurant.ID = firebaseUserID

	createdRestaurant, err := h.restaurantService.CreateRestaurant(restaurant)
	if err != nil {
		h.logger.WithError(err).Errorf("Request ID %s: Failed to create restaurant", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to create restaurant")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdRestaurant)
	h.logger.Infof("Request ID %s: Finished processing request to create restaurant.", r.Context().Value(chimiddleware.RequestIDKey))
}

// GetRestaurant godoc
//
//	@Summary		Get a restaurant
//	@Description	Get a restaurant by the user ID
//	@Tags			restaurants
//	@Produce		json
//	@Security		jwt
//	@Success		200	{object}	models.Restaurant	"Restaurant"
//	@Failure		404	{string}	string				"Restaurant not found"
//	@Failure		500	{string}	string				"Failed to get restaurant"
//	@Router			/restaurants/me [get]
func (h *RestaurantHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof("Request ID %s: Received request to get restaurant.", r.Context().Value(chimiddleware.RequestIDKey))
	// Extract the user ID from the request context
	fbUID := r.Context().Value("UID").(string)

	if fbUID == "" {
		h.logger.Errorf("Request ID %s: Failed to get user ID from request context", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	restaurant, err := h.restaurantService.GetRestaurant(fbUID)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			h.logger.WithError(err).Errorf("Request ID %s: Restaurant not found", r.Context().Value(chimiddleware.RequestIDKey))
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "restaurant not found")
			return
		}
		h.logger.WithError(err).Errorf("Request ID %s: Failed to get restaurant", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get restaurant")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurant)
	h.logger.Infof("Request ID %s: Finished processing request to get restaurant.", r.Context().Value(chimiddleware.RequestIDKey))
}

// UpdateRestaurant godoc
//
//	@Summary		Update a restaurant
//	@Description	Update a restaurant with the given request body
//	@Tags			restaurants
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.UpdateRestaurantRequest	true	"Update Restaurant Request"
//	@Security		jwt
//	@Success		200	{object}	models.Restaurant	"Updated Restaurant"
//	@Failure		400	{string}	string				"Invalid request body"
//	@Failure		500	{string}	string				"Failed to update restaurant"
//	@Router			/restaurants/me [patch]
func (h *RestaurantHandler) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof("Request ID %s: Received request to update restaurant.", r.Context().Value(chimiddleware.RequestIDKey))

	updateRestaurantRequest := &models.UpdateRestaurantRequest{}
	err := json.NewDecoder(r.Body).Decode(updateRestaurantRequest)
	if err != nil {
		h.logger.WithError(err).Errorf("Request ID %s: Failed to decode request body", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	err = h.validator.Struct(updateRestaurantRequest)
	if err != nil {
		h.logger.WithError(err).Errorf("Request ID %s: Invalid request body", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	restaurant := utils.MapUpdateRestaurantRequestToRestaurant(updateRestaurantRequest)

	firebaseUserID, ok := r.Context().Value("UserID").(string)

	if !ok {
		h.logger.Errorf("Request ID %s: Failed to get user ID from request context", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	restaurant.ID = firebaseUserID

	updatedRestaurant, err := h.restaurantService.UpdateRestaurant(restaurant)
	if err != nil {
		h.logger.WithError(err).Errorf("Request ID %s: Failed to update restaurant", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to update restaurant")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedRestaurant)
	h.logger.Infof("Request ID %s: Finished processing request to update restaurant.", r.Context().Value(chimiddleware.RequestIDKey))
}

// GetLogoUploadURL godoc
//
//	@Summary		Get a signed URL to upload a restaurant logo
//	@Description	Get a signed URL to upload a restaurant logo to the S3 bucket
//	@Tags			restaurants
//	@Produce		json
//	@Security		jwt
//	@Success		200	{object}	models.RestaurantLogoUploadResponse	"Presigned URL"
//	@Failure		500	{string}	string								"Failed to get upload URL"
//	@Router			/restaurants/logo/uploadurl [get]
func (h *RestaurantHandler) GetLogoUploadURL(w http.ResponseWriter, r *http.Request) {
	h.logger.Infof("Request ID %s: Received request to get logo upload URL.", r.Context().Value(chimiddleware.RequestIDKey))
	// Extract the user ID from the request context
	UID := r.Context().Value("UID").(string)

	if UID == "" {
		h.logger.Errorf("Request ID %s: Failed to get user ID from request context", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	presignedURL, storedFileURL, err := h.restaurantService.GetLogoUploadURL(UID, h.config.RestaurantLogosBucket)
	if err != nil {
		h.logger.WithError(err).Errorf("Request ID %s: Failed to get upload URL", r.Context().Value(chimiddleware.RequestIDKey))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get upload URL")
		return
	}

	// Make sure the PresignedURL and StoredFileURL are correctly URL encoded
	presignedURLResposne := &models.RestaurantLogoUploadResponse{
		PresignedURL:  presignedURL,
		StoredFileURL: storedFileURL,
		Description:   "The stored_file_url is the URL you have to save in the restaurant table in the logo_url column",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presignedURLResposne)
	h.logger.Infof("Request ID %s: Finished processing request to get logo upload URL.", r.Context().Value(chimiddleware.RequestIDKey))
}
