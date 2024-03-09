package services

import (
	"Tamra/internal/app/tamra/repositories"
	"Tamra/internal/pkg/models"
	"fmt"

	"github.com/sirupsen/logrus"
)

// UserRepository is an interface for the user repository. Since it is a pointer to the implementation of the interface
// we shouldn't import it as a pointer to the interface, since when you pass an interface as a parameter, it is already a pointer to the underlying implementation.
type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUser(userID string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	GetUsers() ([]*models.User, error)
}

type UserServiceImpl struct {
	userRepository repositories.UserRepository
	logger         logrus.FieldLogger
}

func NewUserService(userRepository repositories.UserRepository, logger logrus.FieldLogger) UserService {
	return &UserServiceImpl{userRepository: userRepository, logger: logger}
}

func (s *UserServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	createdUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return createdUser, nil
}

func (s *UserServiceImpl) GetUser(userID string) (*models.User, error) {
	user, err := s.userRepository.GetUser(userID)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateUser(user *models.User) (*models.User, error) {
	updatedUser, err := s.userRepository.UpdateUser(user)
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}

func (s *UserServiceImpl) GetUsers() ([]*models.User, error) {
	users, err := s.userRepository.GetUsers()
	if err != nil {
		// Wrap the error returned by the repository and add some context
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

// func (s *UserServiceImpl) DeleteUser(id int) error {
// 	err := s.userRepository.DeleteUser(id)
// 	if err != nil {
// 		// Wrap the error returned by the repository and add some context
// 		return fmt.Errorf("failed to delete user: %w", err)
// 	}
// 	return nil
// }
