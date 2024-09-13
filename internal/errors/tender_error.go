package errors

import (
	"fmt"
	"net/http"
)

func TenderNotFound(id string) *AppError {
	return &AppError{
		Message: "Tender not found: " + id,
		Code:    http.StatusNotFound,
	}
}

func TenderNotPublished(id string) *AppError {
	return &AppError{
		Message: "Not enough permissions, tender not published: " + id,
		Code:    http.StatusForbidden,
	}
}

func TenderRollbackNotFound(id string, version uint) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Tender rollback not found: %s, version: %d", id, version),
		Code:    http.StatusNotFound,
	}
}

var FailedToCreateTender = &AppError{
	Message: "Failed to create tender",
	Code:    http.StatusUnprocessableEntity,
}

var FailedToSaveTenderRollback = &AppError{
	Message: "Failed to save tender rollback",
	Code:    http.StatusUnprocessableEntity,
}

func FailedToUpdateTender(id string) *AppError {
	return &AppError{
		Message: "Failed to update tender: " + id,
		Code:    http.StatusBadRequest,
	}
}
