package db

import (
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

type TenderRepository struct {
	PostgresStorage
}

func NewTenderRepository(db *PostgresStorage) *TenderRepository {
	return &TenderRepository{*db}
}

func (self TenderRepository) CheckTenderExists(id string) *errors.AppError {
	var tenderExists bool
	err := self.Database.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM tender WHERE id = $1
		)`,
		id,
	).Scan(&tenderExists)
	if err != nil {
		return errors.DatabaseError
	}
	if !tenderExists {
		return errors.TenderNotFound(id)
	}
	return nil
}

func (self TenderRepository) GetTenderById(id string) (*models.TenderDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`select * from tender 
		where id = $1 
		limit 1;`,
		id,
	)
	return models.NewTenderDbModel(row, errors.TenderNotFound(id))
}

func (self TenderRepository) GetTendersList(limit, offset uint, serviceTypes []enums.ServiceType) ([]*models.TenderDbModel, *errors.AppError) {
	var rows *sql.Rows
	var err error
	if len(serviceTypes) == 0 {
		rows, err = self.runQueryWithoutServiceTypes(limit, offset)
	} else {
		rows, err = self.runQueryWithServiceTypes(limit, offset, serviceTypes)
	}
	if err != nil {
		return nil, errors.DatabaseError
	}
	defer rows.Close()
	return models.NewTenderDbModelsList(rows)
}

func (self TenderRepository) runQueryWithServiceTypes(limit, offset uint, serviceTypes []enums.ServiceType) (*sql.Rows, error) {
	query := "SELECT * FROM tender WHERE service_type IN ("
	placeholders := make([]string, len(serviceTypes))
	args := make([]any, len(serviceTypes)+2)

	for i, serviceType := range serviceTypes {
		placeholders[i] = fmt.Sprintf("$%d", i+1) // $1, $2, $3...
		args[i] = string(serviceType)
	}

	query += strings.Join(placeholders, ", ") + ")"

	query += fmt.Sprintf("ORDER BY name ASC LIMIT $%d OFFSET $%d", len(serviceTypes)+1, len(serviceTypes)+2)
	args[len(serviceTypes)] = limit
	args[len(serviceTypes)+1] = offset

	return self.Database.Query(query, args...)
}

func (self TenderRepository) runQueryWithoutServiceTypes(limit, offset uint) (*sql.Rows, error) {
	return self.Database.Query(
		`select * from tender 
		ORDER BY name ASC
		limit $1 offset $2;`,
		limit,
		offset,
	)
}

func (self TenderRepository) GetTendersListByUsername(username string, limit, offset uint) ([]*models.TenderDbModel, *errors.AppError) {
	rows, err := self.Database.Query(
		`SELECT * FROM tender 
		where creator_username = $1
		ORDER BY name ASC
		limit $2
		offset $3;`,
		username,
		limit,
		offset,
	)
	if err != nil {
		return nil, errors.DatabaseError
	}
	defer rows.Close()
	return models.NewTenderDbModelsList(rows)
}

func (self TenderRepository) GetTenderStatus(id string) (enums.TenderStatus, *errors.AppError) {
	var status enums.TenderStatus
	err := self.Database.QueryRow(
		`select status from tender 
		where id = $1 
		limit 1;`,
		id,
	).Scan(&status)
	if err != nil {
		return "", errors.TenderNotFound(id)
	}
	return status, nil
}

func (self TenderRepository) CreateTender(
	createModel *models.TenderCreateModel,
) (*models.TenderDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`insert into tender 
		(name, description, service_type, creator_username, organization_id) 
		values ($1, $2, $3, $4, $5) 
		returning *`,
		createModel.Name,
		createModel.Description,
		string(createModel.ServiceType),
		createModel.CreatorUsername,
		createModel.OrganizationId,
	)
	return models.NewTenderDbModel(row, errors.FailedToCreateTender)
}

func (self TenderRepository) UpdateTenderParams(
	id string,
	updateModel *models.TenderUpdateModel,
) (*models.TenderDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`update tender
		set name = $1, description = $2, service_type = $3, version = version + 1
		where id = $4
		returning *`,
		updateModel.Name,
		updateModel.Description,
		string(updateModel.ServiceType),
		id,
	)
	return models.NewTenderDbModel(row, errors.FailedToUpdateTender(id))
}

func (self TenderRepository) UpdateTenderStatus(
	id string,
	status enums.TenderStatus,
) (*models.TenderDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`update tender
		set status = $1, version = version + 1
		where id = $2
		returning *`,
		string(status),
		id,
	)
	return models.NewTenderDbModel(row, errors.FailedToUpdateTender(id))
}

func (self TenderRepository) RollbackTender(rollbackModel *models.TenderRollbackDbModel) (*models.TenderDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`update tender set
		name = $1,
		description = $2,
		status = $3,
		service_type = $4,
		creator_username = $5,
		organization_id = $6,
		version = version + 1
		where id = $7
		returning *;`,
		rollbackModel.Name,
		rollbackModel.Description,
		rollbackModel.Status,
		rollbackModel.ServiceType,
		rollbackModel.CreatorUsername,
		rollbackModel.OrganizationId,
		rollbackModel.TenderId,
	)
	return models.NewTenderDbModel(row, errors.FailedToUpdateTender(rollbackModel.TenderId))
}
