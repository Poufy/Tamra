package utils

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
	ErrOrderNotAccepted = errors.New("order not accepted")
)
