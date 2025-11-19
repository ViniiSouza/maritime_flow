package utils

import "errors"

type ErrorResponse struct {
	Code    int
	Message string
}

var (
	ErrInvalidInput = errors.New("invalid request body")
	ErrInvalidUUID  = errors.New("invalid or bad formated uuid")
)
