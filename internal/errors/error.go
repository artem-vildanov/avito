package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Message string
	Code    int
}

func (self AppError) Error() string {
	return fmt.Sprintf("Error message: %s, Code: %d; \n", self.Message, self.Code)
}

var InternalError = &AppError{
	Message: "Internal error",
	Code:    http.StatusInternalServerError,
}

var DatabaseError = &AppError{
	Message: "Database error",
	Code:    http.StatusInternalServerError,
}
