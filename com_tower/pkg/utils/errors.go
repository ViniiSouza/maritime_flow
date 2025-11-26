package utils

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type ErrorResponse struct {
	Code    int
	Message string
}

type UnhandledHttpError struct {
	StatusCode int
	Body       string
}

func (e *UnhandledHttpError) Error() string {
	return fmt.Sprintf("not handled http response: status: %d, body: %s", e.StatusCode, e.Body)
}

func HttpErrorNotHandled(statusCode int, body io.ReadCloser) error {
	readBody, _ := io.ReadAll(body)
	cleanReadBody := strings.ReplaceAll(string(readBody), "\n", "")
	cleanReadBody = strings.TrimSpace(cleanReadBody)

	return &UnhandledHttpError{
		StatusCode: statusCode,
		Body:       cleanReadBody,
	}
}

var (
	ErrInvalidInput         = errors.New("invalid request body")
	ErrInvalidUUID          = errors.New("invalid or bad formated uuid")
	ErrLeaderUnreachable    = errors.New("failed to communicate with leader")
	ErrStructureUnreachable = errors.New("failed to communicate with structure")
)
