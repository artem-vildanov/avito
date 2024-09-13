package handlers

import (
	"avito/internal/db"
	"avito/internal/errors"
)

type ReviewHandler struct {
	reviewRepos *db.ReviewRepository
}

func NewReviewHandler(storage *db.PostgresStorage) *ReviewHandler {
	return &ReviewHandler{
		reviewRepos: db.NewReviewRepository(storage),
	}
}

func (self ReviewHandler) LeaveFeedback(ctx *Context) *errors.AppError {
	return nil
}

func (self ReviewHandler) GetReviews(ctx *Context) *errors.AppError {
	return nil
}
