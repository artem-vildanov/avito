package db

import (
	"avito/internal/errors"
	"avito/internal/models"
)

type EmployeeRepository struct {
	PostgresStorage
}

func NewEmployeeRepository(storage *PostgresStorage) *EmployeeRepository {
	return &EmployeeRepository{*storage}
}

func (self EmployeeRepository) GetEmployeeById(id string) (*models.EmployeeDbModel, *errors.AppError) {
	row := self.Database.QueryRow("SELECT * FROM employee WHERE id = $1 limit 1;", id)
	return models.NewEmployeeDbModel(row, errors.EmployeeNotFoundById(id))
}

func (self EmployeeRepository) GetEmployeeByUsername(username string) (*models.EmployeeDbModel, *errors.AppError) {
	row := self.Database.QueryRow("SELECT * FROM employee WHERE username = $1", username)
	return models.NewEmployeeDbModel(row, errors.EmployeeNotFoundByUsername(username))
}

func (self EmployeeRepository) CheckEmployeeExistsById(id string) *errors.AppError {
	var userExists bool
	err := self.Database.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM employee WHERE id = $1
		)`,
		id,
	).Scan(&userExists)
	if err != nil {
		return errors.DatabaseError
	}
	if !userExists {
		return errors.EmployeeNotFoundById(id)
	}
	return nil
}

func (self EmployeeRepository) CheckEmployeeExistsByUsername(username string) *errors.AppError {
	var userExists bool
	err := self.Database.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM employee WHERE username = $1
		)`,
		username,
	).Scan(&userExists)
	if err != nil {
		return errors.DatabaseError
	}
	if !userExists {
		return errors.EmployeeNotFoundByUsername(username)
	}
	return nil
}

func (self EmployeeRepository) CheckEmployeeIsResponsible(username, organizationId string) *errors.AppError {
	var isResponsible bool
	err := self.Database.QueryRow(
		`SELECT EXISTS (
			SELECT 1
			FROM employee
			INNER JOIN organization_responsible ON employee.id = organization_responsible.user_id
			WHERE employee.username = $1
			AND organization_responsible.organization_id = $2
		);`,
		username,
		organizationId,
	).Scan(&isResponsible)
	if err != nil {
		return errors.DatabaseError
	}
	if !isResponsible {
		return errors.NotEnoughPermissions(username)
	}
	return nil
}

func (self EmployeeRepository) CheckEmployeeHasBidsForTender(employeeId, tenderId string) *errors.AppError {
	var hasBids bool
	err := self.Database.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM bid WHERE author_id = $1 and tender_id = $2
		)`,
		employeeId,
		tenderId,
	).Scan(&hasBids)
	if err != nil {
		return errors.DatabaseError
	}
	if !hasBids {
		return errors.NotEnoughPermissions(employeeId)
	}
	return nil
}

//func (self EmployeeRepository) GetEmployee(id int) *errors.AppError {
//	return nil
//}
//
//func (self EmployeeRepository) CreateEmployee(
//	username,
//	firstName,
//	lastName string,
//) (string, *errors.AppError) {
//	var id int
//	err := self.Database.QueryRow(
//		`insert into
//		employee (username, first_name, last_name)
//		values ($1, $2, $3)
//		returning id;`,
//		username,
//		firstName,
//		lastName,
//	).Scan(&id)
//
//	if err != nil {
//		return -1, errors.FailedToCreate()
//	}
//
//	return id, nil
//}
