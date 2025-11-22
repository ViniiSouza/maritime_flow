package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetContextAndExecJSONWithErrorResponse(c *gin.Context, err error) {
	var httpStatus int
	switch {
	case errors.Is(err, ErrInvalidInput):
		httpStatus = http.StatusBadRequest
	default:
		httpStatus = http.StatusInternalServerError
	}
	
	response := ErrorResponse{
		Code: httpStatus,
		Message: err.Error(),
	}

	c.JSON(httpStatus, response)
}
