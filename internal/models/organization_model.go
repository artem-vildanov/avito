package models

import (
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/utils"
)

type OrganizationDbModel struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        enums.OrganizationType `json:"type"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

func NewOrganizationDbModel(row utils.Scannable, returnErr *errors.AppError) (*OrganizationDbModel, *errors.AppError) {
	model := new(OrganizationDbModel)
	err := row.Scan(
		&model.Id,
		&model.Name,
		&model.Description,
		&model.Type,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return nil, returnErr
	}
	return model, nil
}
