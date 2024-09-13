package errors

import "net/http"

var FailedToCreateReview = &AppError{
	Message: "Failed to create review",
	Code:    http.StatusUnprocessableEntity,
}
