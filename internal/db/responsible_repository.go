package db

import "avito/internal/errors"

type ResponsibleRepository struct {
	PostgresStorage
}

func NewResponsibleRepository(storage *PostgresStorage) *ResponsibleRepository {
	return &ResponsibleRepository{*storage}
}

func (self ResponsibleRepository) CountResponsiblesByOrgId(orgId string) (uint, *errors.AppError) {
	var count uint
	err := self.Database.QueryRow(
		`SELECT COUNT(*) 
		FROM organization_responsible 
		WHERE organization_id = $1;`,
		orgId,
	).Scan(&count)

	if err != nil {
		return 0, errors.DatabaseError
	}

	return count, nil
}

//func (self *ResponsibleRepository) CreateResponsible(
//	orgainzationId,
//	userId int,
//) (int, *errors.AppError) {
//	var id int
//	err := self.Database.QueryRow(
//		`insert into organization_responsible
//		(organization_id, user_id)
//		values ($1, $2)`,
//		orgainzationId,
//		userId,
//	).Scan(&id)
//
//	if err != nil {
//		return -1, errors.FailedToCreate()
//	}
//
//	return id, nil
//}
