package services

import (
	"avito/internal/db"
	"avito/internal/errors"
)

type EmployeeService struct {
	employeeRepos *db.EmployeeRepository
}

func NewEmployeeService(storage *db.PostgresStorage) *EmployeeService {
	return &EmployeeService{
		employeeRepos: db.NewEmployeeRepository(storage),
	}
}

func (self EmployeeService) CheckEmployeeIsResponsible(username, organizationId string) *errors.AppError {
	if err := self.employeeRepos.CheckEmployeeExistsByUsername(username); err != nil {
		return err
	}
	err := self.employeeRepos.CheckEmployeeIsResponsible(username, organizationId)
	if err != nil {
		return err
	}
	return nil
}
