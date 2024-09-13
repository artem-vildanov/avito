package db

import (
	"avito/internal/errors"
	"avito/internal/models"
)

type OrganizationRepository struct {
	PostgresStorage
}

func NewOrganizationRepository(storage *PostgresStorage) *OrganizationRepository {
	return &OrganizationRepository{*storage}
}

func (self OrganizationRepository) GetOrganizationById(id string) (*models.OrganizationDbModel, *errors.AppError) {
	row := self.Database.QueryRow(`select * from organization where id = $1 limit 1;`, id)
	return models.NewOrganizationDbModel(row, errors.OrganizationNotFound(id))
}

func (self OrganizationRepository) CheckOrganizationExists(id string) *errors.AppError {
	var orgExists bool
	err := self.Database.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM organization WHERE id = $1
		)`,
		id,
	).Scan(&orgExists)
	if err != nil {
		return errors.DatabaseError
	}
	if !orgExists {
		return errors.OrganizationNotFound(id)
	}
	return nil
}

//func (self OrganizationRepository) GetOrganization(id int) *errors.AppError {
//	return nil
//}
//
//func (self OrganizationRepository) CreateOrganization(
//	name,
//	description string,
//	orgType enums.OrganizationType,
//) (int, *errors.AppError) {
//	var id int
//	err := self.Database.QueryRow(
//		`insert into organization
//		(name, description, type)
//		values ($1, $2, $3)
//		returning id`,
//		name,
//		description,
//		string(orgType),
//	).Scan(&id)
//
//	if err != nil {
//		return -1, errors.FailedToCreate()
//	}
//
//	return id, nil
//}
