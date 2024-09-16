package models

import (
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/utils"
	"database/sql"
)

type TenderDbModel struct {
	Id              string             `json:"id"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Status          enums.TenderStatus `json:"status"`
	ServiceType     enums.ServiceType  `json:"service_type"`
	CreatorUsername string             `json:"creator_username"`
	OrganizationId  string             `json:"organization_id"`
	Version         int                `json:"version"`
	CreatedAt       string             `json:"created_at"`
}

func NewTenderDbModel(row utils.Scannable, returnErr *errors.AppError) (*TenderDbModel, *errors.AppError) {
	model := new(TenderDbModel)
	err := row.Scan(
		&model.Id,
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

func NewTenderDbModelsList(rows *sql.Rows) ([]*TenderDbModel, *errors.AppError) {
	modelsList := make([]*TenderDbModel, 0)
	for rows.Next() {
		model, err := NewTenderDbModel(rows, errors.DatabaseError)
		if err != nil {
			return nil, err
		}
		modelsList = append(modelsList, model)
	}
	return modelsList, nil
}

type TenderCreateModel struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	ServiceType     enums.ServiceType `json:"serviceType"`
	OrganizationId  string            `json:"organizationId"`
	CreatorUsername string            `json:"creatorUsername"`
}

type TenderUpdateModel struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ServiceType enums.ServiceType `json:"serviceType"`
}

type TenderDtoModel struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Status      enums.TenderStatus `json:"status"`
	ServiceType enums.ServiceType  `json:"serviceType"`
	Version     int                `json:"version"`
	CreatedAt   string             `json:"createdAt"`
}

func NewTenderDtoModel(dbModel *TenderDbModel) *TenderDtoModel {
	return &TenderDtoModel{
		Id:          dbModel.Id,
		Name:        dbModel.Name,
		Description: dbModel.Description,
		Status:      dbModel.Status,
		ServiceType: dbModel.ServiceType,
		Version:     dbModel.Version,
		CreatedAt:   dbModel.CreatedAt,
	}
}

func NewTenderDtoModelsList(dbModels []*TenderDbModel) []*TenderDtoModel {
	dtoModelsList := make([]*TenderDtoModel, 0)
	for _, dbModel := range dbModels {
		dtoModelsList = append(dtoModelsList, NewTenderDtoModel(dbModel))
	}
	return dtoModelsList
}
