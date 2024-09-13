package models

import (
	"avito/internal/errors"
	"avito/internal/utils"
	"database/sql"
)

type ReviewDbModel struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func NewReviewDbModel(row utils.Scannable, returnErr *errors.AppError) (*ReviewDbModel, *errors.AppError) {
	model := new(ReviewDbModel)
	err := row.Scan(
		&model.Id,
		&model.Description,
		&model.CreatedAt,
	)
	if err != nil {
		return nil, returnErr
	}
	return model, nil
}

func NewReviewDbModelsList(rows *sql.Rows) ([]*ReviewDbModel, *errors.AppError) {
	var models []*ReviewDbModel
	for rows.Next() {
		model, err := NewReviewDbModel(rows, errors.DatabaseError)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

type ReviewDtoModel struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}
