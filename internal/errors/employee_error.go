package errors

import (
	"fmt"
	"net/http"
)

func EmployeeNotFoundById(id string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Employee with id: %s not found", id),
		Code:    http.StatusUnauthorized,
	}
}

func EmployeeNotFoundByUsername(username string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Employee with username: %s not found", username),
		Code:    http.StatusUnauthorized,
	}
}

func NotEnoughPermissions(username string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Not enough permissions for user: %s to complete action", username),
		Code:    http.StatusForbidden,
	}
}

func AlreadyApprovedBid(bidId, employeeId string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Bid %s already approved by %s", bidId, employeeId),
		Code:    http.StatusBadRequest,
	}
}
