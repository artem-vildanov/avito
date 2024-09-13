package errors

import (
	"fmt"
	"net/http"
)

func InvalidRequestBody() *AppError {
	return &AppError{
		Message: "Invalid request body",
		Code:    http.StatusBadRequest,
	}
}

func RequiredRequestParamNotProvided(requiredParam string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Required request parameter not provided: %s", requiredParam),
		Code:    http.StatusBadRequest,
	}
}

func InvalidRequestParam(invalidParam string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Invalid request parameter: %s", invalidParam),
		Code:    http.StatusBadRequest,
	}
}
