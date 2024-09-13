package db

import (
	"avito/internal/errors"
	"avito/internal/models"
)

type ReviewRepository struct {
	PostgresStorage
}

func NewReviewRepository(storage *PostgresStorage) *ReviewRepository {
	return &ReviewRepository{*storage}
}

func (self ReviewRepository) GetReviewsByBidsAuthorUsername(
	username string,
	limit,
	offset int,
) ([]*models.ReviewDbModel, *errors.AppError) {
	rows, err := self.Database.Query(
		`SELECT * 
		FROM review
		JOIN bid ON review.bid_id = bid.id
		JOIN employee ON bid.author_id = employee.id
		WHERE employee.username = $1
		limit $2 offset $3;`,
		username,
		limit,
		offset,
	)
	if err != nil {
		return nil, errors.DatabaseError
	}
	return models.NewReviewDbModelsList(rows)
}

func (self ReviewRepository) NewReview(bid_id, description string) (*models.ReviewDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`insert into review 
		(bid_id, description) 
		values ($1, $2) 
		returning *;`,
		bid_id,
		description,
	)
	return models.NewReviewDbModel(row, errors.FailedToCreateReview)
}
