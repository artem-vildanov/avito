package models

import (
	"avito/internal/errors"
	"avito/internal/utils"
)

type EmployeeDbModel struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewEmployeeDbModel(row utils.Scannable, returnErr *errors.AppError) (*EmployeeDbModel, *errors.AppError) {
	model := new(EmployeeDbModel)
	err := row.Scan(
		&model.Id,
		&model.Username,
		&model.FirstName,
		&model.LastName,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		println(err.Error())
		return nil, returnErr
	}
	return model, nil
}
