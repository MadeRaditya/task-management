package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MadeRaditya/task-management/internal/apperror"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

func respondError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, ErrorResponse{
			Error: ErrorBody{Code: appErr.Code, Message: appErr.Message},
		})
		return
	}

	log.Printf("unexpected error: %v", err)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{Code: "INTERNAL_ERROR", Message: "Internal server error"},
	})
}

func apperrorBadRequest() *apperror.AppError {
	return apperror.BadRequest("INVALID_ID", "Invalid task id")
}

func badRequestFromBinding(err error) *apperror.AppError {
	return apperror.BadRequest("INVALID_BODY", err.Error())
}
