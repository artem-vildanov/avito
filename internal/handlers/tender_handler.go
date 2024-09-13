package handlers

import (
	"avito/internal/db"
	"avito/internal/errors"
	"avito/internal/models"
	"avito/internal/services"
)

type TenderHandler struct {
	tenderRepos         *db.TenderRepository
	tenderRollbackRepos *db.TenderRollbackRepository
	employeeService     *services.EmployeeService
}

func NewTenderHandler(storage *db.PostgresStorage) *TenderHandler {
	return &TenderHandler{
		tenderRepos:         db.NewTenderRepository(storage),
		tenderRollbackRepos: db.NewTenderRollbackRepository(storage),
		employeeService:     services.NewEmployeeService(storage),
	}
}

func (self TenderHandler) GetTendersList(ctx *Context) *errors.AppError {
	limit, offset, err := ctx.GetLimitAndOffsetRequestParams()
	if err != nil {
		return err
	}

	serviceTypes, err := ctx.GetServiceTypesRequestParam()
	if err != nil {
		return err
	}

	tendersDbModels, err := self.tenderRepos.GetTendersList(limit, offset, serviceTypes)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(200, models.NewTenderDtoModelsList(tendersDbModels))
}

func (self TenderHandler) GetTendersListByUsername(ctx *Context) *errors.AppError {
	limit, err := ctx.GetLimitRequestParam()
	if err != nil {
		return err
	}
	offset, err := ctx.GetOffsetRequestParam()
	if err != nil {
		return err
	}
	username, err := ctx.GetUsernameRequestParam()
	if err != nil {
		return err
	}
	tendersDbModels, err := self.tenderRepos.GetTendersListByUsername(username, limit, offset)
	if err != nil {
		return err
	}
	return ctx.RespondWithJson(200, models.NewTenderDtoModelsList(tendersDbModels))
}

func (self TenderHandler) CreateTender(ctx *Context) *errors.AppError {
	createModel, err := GetModelFromRequest[models.TenderCreateModel](ctx.Request.Body)
	if err != nil {
		return err
	}
	if err := self.employeeService.CheckEmployeeIsResponsible(
		createModel.CreatorUsername,
		createModel.OrganizationId,
	); err != nil {
		return err
	}
	tenderDbModel, err := self.tenderRepos.CreateTender(createModel)
	if err != nil {
		return err
	}
	err = self.tenderRollbackRepos.SaveTenderRollback(tenderDbModel)
	if err != nil {
		return err
	}
	return ctx.RespondWithJson(200, models.NewTenderDtoModel(tenderDbModel))
}

func (self TenderHandler) GetTenderStatus(ctx *Context) *errors.AppError {
	username, tenderId, err := self.getUsernameAndTenderIdReqParams(ctx)
	if err != nil {
		return err
	}
	tenderDbModel, err := self.getTenderIfEmployeeResponsible(tenderId, username)
	if err != nil {
		return err
	}
	return ctx.RespondWithJson(200, tenderDbModel.Status)
}

func (self TenderHandler) UpdateTenderStatus(ctx *Context) *errors.AppError {
	username, tenderId, err := self.getUsernameAndTenderIdReqParams(ctx)
	if err != nil {
		return err
	}
	status, err := ctx.GetTenderStatusRequestParam()
	if err != nil {
		return err
	}
	if err := self.checkEmployeeTenderAccess(tenderId, username); err != nil {
		return err
	}
	tenderDbModel, err := self.tenderRepos.UpdateTenderStatus(tenderId, status)
	if err != nil {
		return err
	}
	err = self.tenderRollbackRepos.SaveTenderRollback(tenderDbModel)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(200, models.NewTenderDtoModel(tenderDbModel))
}

func (self TenderHandler) UpdateTenderParams(ctx *Context) *errors.AppError {
	username, tenderId, err := self.getUsernameAndTenderIdReqParams(ctx)
	if err != nil {
		return err
	}
	updateModel, err := GetModelFromRequest[models.TenderUpdateModel](ctx.Request.Body)
	if err != nil {
		return err
	}
	if err := self.checkEmployeeTenderAccess(tenderId, username); err != nil {
		return err
	}
	tenderDbModel, err := self.tenderRepos.UpdateTenderParams(tenderId, updateModel)
	if err != nil {
		return err
	}
	err = self.tenderRollbackRepos.SaveTenderRollback(tenderDbModel)
	if err != nil {
		return err
	}
	return ctx.RespondWithJson(200, models.NewTenderDtoModel(tenderDbModel))
}

func (self TenderHandler) RollbackTender(ctx *Context) *errors.AppError {
	username, tenderId, err := self.getUsernameAndTenderIdReqParams(ctx)
	if err != nil {
		return err
	}

	version, err := ctx.GetVersionPathParam()
	if err != nil {
		return err
	}

	if err := self.checkEmployeeTenderAccess(tenderId, username); err != nil {
		return err
	}

	tenderRollbackDbModel, err := self.tenderRollbackRepos.GetTenderRollback(tenderId, version)
	if err != nil {
		return err
	}

	tenderDbModel, err := self.tenderRepos.RollbackTender(tenderRollbackDbModel)
	if err != nil {
		return err
	}

	err = self.tenderRollbackRepos.SaveTenderRollback(tenderDbModel)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(200, models.NewTenderDtoModel(tenderDbModel))
}

func (self TenderHandler) getUsernameAndTenderIdReqParams(ctx *Context) (string, string, *errors.AppError) {
	username, err := ctx.GetUsernameRequestParam()
	if err != nil {
		return "", "", err
	}
	tenderId, err := ctx.GetTenderIdPathParam()
	if err != nil {
		return "", "", err
	}
	return username, tenderId, nil
}

func (self TenderHandler) checkEmployeeTenderAccess(tenderId, username string) *errors.AppError {
	_, err := self.getTenderIfEmployeeResponsible(tenderId, username)
	if err != nil {
		return err
	}
	return nil
}

func (self TenderHandler) getTenderIfEmployeeResponsible(tenderId, username string) (*models.TenderDbModel, *errors.AppError) {
	tenderDbModel, err := self.tenderRepos.GetTenderById(tenderId)
	if err != nil {
		return nil, err
	}
	if err := self.employeeService.CheckEmployeeIsResponsible(
		username,
		tenderDbModel.OrganizationId,
	); err != nil {
		return nil, err
	}
	return tenderDbModel, nil
}
