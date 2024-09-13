package models

import (
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/utils"
	"database/sql"
)

type BidDbModel struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      enums.BidStatus  `json:"status"`
	TenderId    string           `json:"tender_id"`
	AuthorId    string           `json:"author_id"`
	AuthorType  enums.AuthorType `json:"author_type"`
	Version     int              `json:"version"`
	CreatedAt   string           `json:"created_at"`
}

func NewBidDbModel(row utils.Scannable, returnErr *errors.AppError) (*BidDbModel, *errors.AppError) {
	model := new(BidDbModel)
	err := row.Scan(
		&model.Id,
		&model.Name,
		&model.Description,
		&model.Status,
		&model.TenderId,
		&model.AuthorId,
		&model.AuthorType,
		&model.Version,
		&model.CreatedAt,
	)
	if err != nil {
		println(err.Error())
		return nil, returnErr
	}
	return model, nil
}

func NewBidDbModelsList(rows *sql.Rows) ([]*BidDbModel, *errors.AppError) {
	var modelsList []*BidDbModel
	for rows.Next() {
		model, err := NewBidDbModel(rows, errors.DatabaseError)
		if err != nil {
			return nil, err
		}
		modelsList = append(modelsList, model)
	}
	return modelsList, nil
}

type BidCreateModel struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	TenderId    string           `json:"tenderId"`
	AuthorId    string           `json:"authorId"`
	AuthorType  enums.AuthorType `json:"authorType"`
}

type BidUpdateModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BidDtoModel struct {
	Id         string           `json:"id"`
	Name       string           `json:"name"`
	Status     enums.BidStatus  `json:"status"`
	AuthorType enums.AuthorType `json:"authorType"`
	AuthorId   string           `json:"authorId"`
	Version    int              `json:"version"`
	CreatedAt  string           `json:"createdAt"`
}

func NewBidDtoModel(dbModel *BidDbModel) *BidDtoModel {
	return &BidDtoModel{
		Id:         dbModel.Id,
		Name:       dbModel.Name,
		Status:     dbModel.Status,
		AuthorId:   dbModel.AuthorId,
		AuthorType: dbModel.AuthorType,
		Version:    dbModel.Version,
		CreatedAt:  dbModel.CreatedAt,
	}
}

func NewBidDtoModelsList(dbModels []*BidDbModel) []*BidDtoModel {
	var modelsList []*BidDtoModel
	for _, dbModel := range dbModels {
		modelsList = append(modelsList, NewBidDtoModel(dbModel))
	}
	return modelsList
}
