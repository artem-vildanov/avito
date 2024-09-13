package db

import (
	"avito/internal/errors"
	"avito/internal/models"
)

type TenderRollbackRepository struct {
	PostgresStorage
}

func NewTenderRollbackRepository(storage *PostgresStorage) *TenderRollbackRepository {
	return &TenderRollbackRepository{*storage}
}

func (self TenderRollbackRepository) SaveTenderRollback(rollbackModel *models.TenderDbModel) *errors.AppError {
	_, err := self.Database.Exec(
		`insert into tender_rollback (
			tender_id, 
			name, 
			description, 
			status, 
			service_type, 
			creator_username, 
			organization_id, 
			version, 
			created_at
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		rollbackModel.Id,
		rollbackModel.Name,
		rollbackModel.Description,
		rollbackModel.Status,
		rollbackModel.ServiceType,
		rollbackModel.CreatorUsername,
		rollbackModel.OrganizationId,
		rollbackModel.Version,
		rollbackModel.CreatedAt,
	)
	if err != nil {
		return errors.FailedToSaveTenderRollback
	}
	return nil
}

func (self TenderRollbackRepository) GetTenderRollback(id string, version uint) (*models.TenderRollbackDbModel, *errors.AppError) {
	row := self.Database.QueryRow(
		`select * from tender_rollback 
		where tender_id = $1 and version = $2
		limit 1;`,
		id, version,
	)
	return models.NewTenderRollbackDbModel(row, errors.TenderRollbackNotFound(id, version))
}
