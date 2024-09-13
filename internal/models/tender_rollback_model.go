package models

import (
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/utils"
)

type TenderRollbackDbModel struct {
	Id              string             `json:"id"`
	TenderId        string             `json:"tender_id"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Status          enums.TenderStatus `json:"status"`
	ServiceType     enums.ServiceType  `json:"service_type"`
	CreatorUsername string             `json:"creator_username"`
	OrganizationId  string             `json:"organization_id"`
	Version         int                `json:"version"`
	CreatedAt       string             `json:"created_at"`
}

func NewTenderRollbackDbModel(row utils.Scannable, returnErr *errors.AppError) (*TenderRollbackDbModel, *errors.AppError) {
	var model = new(TenderRollbackDbModel)
	err := row.Scan(
		&model.Id,
		&model.TenderId,
		&model.Name,
		&model.Description,
		&model.Status,
		&model.ServiceType,
		&model.CreatorUsername,
		&model.OrganizationId,
		&model.Version,
		&model.CreatedAt,
	)
	if err != nil {
		println(err.Error())
		return nil, returnErr
	}
	return model, nil
}
