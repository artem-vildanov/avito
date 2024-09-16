package errors

import (
	"fmt"
	"net/http"
)

func BidNotFound(id string) *AppError {
	return &AppError{
		Message: "Bid not found: " + id,
		Code:    http.StatusNotFound,
	}
}

func BidNotPublished(id string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Bid not published: %s, not enough permissions", id),
	}
}

func BidRollbackNotFound(id string, version uint) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Bid rollback not found: id: %s, version: %d", id, version),
		Code:    http.StatusNotFound,
	}
}

var FailedToCreateBid = &AppError{
	Message: "Failed to create bid",
	Code:    http.StatusUnprocessableEntity,
}

var FailedToSaveBidRollback = &AppError{
	Message: "Failed to save bid rollback",
	Code:    http.StatusUnprocessableEntity,
}

func FailedToUpdateBid(id string) *AppError {
	return &AppError{
		Message: "Failed to update bid: " + id,
		Code:    http.StatusUnprocessableEntity,
	}
}
