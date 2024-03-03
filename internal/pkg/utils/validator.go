package utils

import (
	"github.com/go-playground/validator/v10"
)

// NewValidator returns a new instance of the validator.
func NewValidator() *validator.Validate {
	return validator.New()
}
