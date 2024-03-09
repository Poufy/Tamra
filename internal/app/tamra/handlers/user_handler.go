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

// ? Is this okay? Define an interface for the validator so we can pass the validator as a parameter
// ? to the handler without having to import the validator package.
type Validator interface {
	Struct(s interface{}) error
}

type UserHandler struct {
	userService services.UserService
	validator   Validator
	logger      logrus.FieldLogger
}

func NewUserHandler(userService services.UserService, validator Validator, logger logrus.FieldLogger) *UserHandler {
	return &UserHandler{userService: userService, validator: validator, logger: logger}
}

// CreateUser godoc
//
//	@Summary		Create a user
//	@Description	Create a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.CreateUserRequest	true	"Create User Request"
//	@Security		jwt
//	@Success		201	{object}	models.UserResponse	"Created User"
//	@Failure		400	{string}	string				"Invalid request body"
//	@Failure		500	{string}	string				"Failed to create user"
//	@Router			/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	createUserRequest := &models.CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(createUserRequest)
	if err != nil {
		h.logger.WithError(err).Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	err = h.validator.Struct(createUserRequest)
	if err != nil {
		h.logger.WithError(err).Error("invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	// Here we would map the CreateUserRequest to a User and pass it to the service
	// The reason why we map to a user is because the service should not know about the request/response models
	// It should be loosely coupled and only know about the domain models
	user := utils.MapCreateUserRequestToUser(createUserRequest)

	firebaseUserID, ok := r.Context().Value("UID").(string)
	if !ok {
		h.logger.Error("failed to get user ID from request context")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	user.FBUserID = firebaseUserID
	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		h.logger.WithError(err).Error("failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to create user")
		return
	}

	h.logger.Infof("user created: %v", createdUser)

	userResponse := utils.MapUserToUserResponse(createdUser)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Get a user by the user ID
//	@Tags			users
//	@Produce		json
//	@Security		jwt
//	@Success		200	{object}	models.UserResponse
//	@Failure		404	{string}	string	"user not found"
//	@Failure		500	{string}	string	"failed to get user"
//	@Router			/users/me [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("UID").(string)
	if !ok {
		h.logger.Error("failed to get user ID from request context")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		// We use errors.Is instead of checking with == because the error might be wrapped and we want to check the underlying error type.
		if errors.Is(err, utils.ErrNotFound) {
			h.logger.WithError(err).Error("user not found")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "user not found")
			return
		}

		h.logger.WithError(err).Error("failed to get user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user")
		return
	}

	// Map the user to a UserResponse
	userResponse := utils.MapUserToUserResponse(user)

	json.NewEncoder(w).Encode(userResponse)
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body	models.UpdateUserRequest	true	"Update User Request"
//	@Security		jwt
//	@Success		200	{object}	models.UserResponse	"Updated User"
//	@Failure		400	{string}	string				"Invalid request body"
//	@Failure		500	{string}	string				"Failed to update user"
//	@Router			/users/me [patch]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	updateUserRequest := &models.UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(updateUserRequest)
	if err != nil {
		h.logger.WithError(err).Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	err = h.validator.Struct(updateUserRequest)
	if err != nil {
		h.logger.WithError(err).Error("invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invalid request body")
		return
	}

	user := utils.MapUpdateUserRequestToUser(updateUserRequest)

	firebaseUserID, ok := r.Context().Value("UID").(string)

	if !ok {
		h.logger.Error("failed to get user ID from request context")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get user ID from request context")
		return
	}

	// Set the ID of the user to the ID from the URL
	user.FBUserID = firebaseUserID

	updatedUser, err := h.userService.UpdateUser(user)
	if err != nil {
		h.logger.WithError(err).Error("failed to update user")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to update user")
		return
	}

	h.logger.Infof("user updated: %v", updatedUser)

	userResponse := utils.MapUserToUserResponse(updatedUser)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse)
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Get all users
//	@Tags			users
//	@Produce		json
//	@Security		jwt
//	@Success		200	{array}		models.UserResponse
//	@Failure		404	{string}	string	"users not found"
//	@Failure		500	{string}	string	"failed to get users"
//	@Router			/users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetUsers()
	if err != nil {
		// We use errors.Is instead of checking with == because the error might be wrapped and we want to check the underlying error type.
		if errors.Is(err, utils.ErrNotFound) {
			h.logger.WithError(err).Error("users not found")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "users not found")
			return
		}

		h.logger.WithError(err).Error("failed to get users")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to get users")
	}

	h.logger.Info("users retrieved.")

	userResponses := utils.MapUsersToUserResponses(users)

	json.NewEncoder(w).Encode(userResponses)
}

// func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		h.logger.WithError(err).Error("invalid user ID")
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, "invalid user ID")
// 		return
// 	}

// 	err = h.userService.DeleteUser(id)
// 	if err != nil {
// 		h.logger.WithError(err).Error("failed to delete user")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, "failed to delete user")
// 		return
// 	}

// 	h.logger.Infof("user deleted: %v", id)
// 	w.WriteHeader(http.StatusOK)
// }
