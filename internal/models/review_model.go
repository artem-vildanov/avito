package models

import (
	"avito/internal/errors"
	"avito/internal/utils"
	"database/sql"
)

type ReviewDbModel struct {
	Id          string `json:"id"`
	BidId       string `json:"bid_id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func NewReviewDbModel(row utils.Scannable, returnErr *errors.AppError) (*ReviewDbModel, *errors.AppError) {
	model := new(ReviewDbModel)
	err := row.Scan(
		&model.Id,
		&model.BidId,
		&model.Description,
		&model.CreatedAt,
	)
	if err != nil {
		println(err.Error())
		return nil, returnErr
	}
	return model, nil
}

func NewReviewDbModelsList(rows *sql.Rows) ([]*ReviewDbModel, *errors.AppError) {
	models := make([]*ReviewDbModel, 0)
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

func NewReviewDtoModel(dbModel *ReviewDbModel) *ReviewDtoModel {
	return &ReviewDtoModel{
		Id:          dbModel.Id,
		Description: dbModel.Description,
		CreatedAt:   dbModel.CreatedAt,
	}
}

func NewReviewDtoModelList(dbModels []*ReviewDbModel) []*ReviewDtoModel {
	models := make([]*ReviewDtoModel, 0)
	for _, model := range dbModels {
		models = append(models, NewReviewDtoModel(model))
	}
	return models
}
