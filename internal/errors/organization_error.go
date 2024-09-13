package errors

import "net/http"

func OrganizationNotFound(id string) *AppError {
	return &AppError{
		Message: "Organization not found: " + id,
		Code:    http.StatusNotFound,
	}
}
