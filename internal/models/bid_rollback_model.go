package models

import (
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/utils"
)

type BidRollbackDbModel struct {
	Id          string           `json:"id"`
	BidId       string           `json:"bid_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      enums.BidStatus  `json:"status"`
	TenderId    string           `json:"tender_id"`
	AuthorId    string           `json:"author_id"`
	AuthorType  enums.AuthorType `json:"author_type"`
	Version     int              `json:"version"`
	CreatedAt   string           `json:"created_at"`
}

func NewBidRollbackDbModel(row utils.Scannable, returnErr *errors.AppError) (*BidRollbackDbModel, *errors.AppError) {
	model := new(BidRollbackDbModel)
	err := row.Scan(
		&model.Id,
		&model.BidId,
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
