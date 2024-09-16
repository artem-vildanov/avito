package handlers

import (
	"avito/internal/db"
	"avito/internal/errors"
	"avito/internal/models"
)

type TenderHandler struct {
	tenderRepos         *db.TenderRepository
	tenderRollbackRepos *db.TenderRollbackRepository
	employeeRepos       *db.EmployeeRepository
}

func NewTenderHandler(storage *db.PostgresStorage) *TenderHandler {
	return &TenderHandler{
		tenderRepos:         db.NewTenderRepository(storage),
		tenderRollbackRepos: db.NewTenderRollbackRepository(storage),
		employeeRepos:       db.NewEmployeeRepository(storage),
	}
}

func (t TenderHandler) GetTendersList(ctx *Context) *errors.AppError {
	limit, offset, err := ctx.GetLimitAndOffsetRequestParams()
	if err != nil {
		return err
	}

	serviceTypes, err := ctx.GetServiceTypesRequestParam()
	if err != nil {
		return err
	}

	tendersDbModels, err := t.tenderRepos.GetTendersList(limit, offset, serviceTypes)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(200, models.NewTenderDtoModelsList(tendersDbModels))
}

func (t TenderHandler) GetTendersListByUsername(ctx *Context) *errors.AppError {
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

	tendersDbModels, err := t.tenderRepos.GetTendersListByUsername(username, limit, offset)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(200, models.NewTenderDtoModelsList(tendersDbModels))
}

func (t TenderHandler) CreateTender(ctx *Context) *errors.AppError {
	createModel, err := GetModelFromRequest[models.TenderCreateModel](ctx.Request.Body)
	if err != nil {
		return err
	}

	if err := t.checkEmployeeExistsAndResponsible(createModel.CreatorUsername, createModel.OrganizationId); err != nil {
		return err
	}

	tenderDbModel, err := t.tenderRepos.CreateTender(createModel)
	if err != nil {
		return err
	}

	err = t.tenderRollbackRepos.SaveTenderRollback(tenderDbModel)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(200, models.NewTenderDtoModel(tenderDbModel))
}

func (t TenderHandler) GetTenderStatus(ctx *Context) *errors.AppError {
	username, tenderId, err := t.getUsernameAndTenderIdReqParams(ctx)
	if err != nil {
		return err
	}
	tenderDbModel, err := t.getTenderIfEmployeeResponsible(tenderId, username)
	if err != nil {
		return err
	}
	return ctx.RespondWithJson(200, tenderDbModel.Status)
}

func (t TenderHandler) UpdateTenderStatus(ctx *Context) *errors.AppError {
	username, tenderId, err := t.getUsernameAndTenderIdReqParams(ctx)
	if err != nil {
		return err
	}
	status, err := ctx.GetTenderStatusRequestParam()
	if err != nil {
		return err
	}
	if err := t.checkEmployeeTenderAccess(tenderId, username); err != nil {
		return err
	}
	tenderDbModel, err := t.tenderRepos.UpdateTenderStatus(tenderId, status)
	if err != nil {
		return err
	}
	err = t.tenderRollbackRepos.SaveTenderRollback(tenderDbModel)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(200, models.NewTenderDtoModel(tenderDbModel))
}

func (t TenderHandler) UpdateTenderParams(ctx *Context) *errors.AppError {
	username, tenderId, err := t.getUsernameAndTenderIdReqParams(ctx)
	if err != nil {
		return err
	}
	updateModel, err := GetModelFromRequest[models.TenderUpdateModel](ctx.Request.Body)
	if err != nil {
		return err
	}
	if err := t.checkEmployeeTenderAccess(tenderId, username); err != nil {
		return err
	}
	tenderDbModel, err := t.tenderRepos.UpdateTenderParams(tenderId, updateModel)
	if err != nil {
		return err
	}
	err = t.tenderRollbackRepos.SaveTenderRollback(tenderDbModel)
	if err != nil {
		return err
	}
	return ctx.RespondWithJson(200, models.NewTenderDtoModel(tenderDbModel))
}

func (t TenderHandler) RollbackTender(ctx *Context) *errors.AppError {
	username, tenderId, err := t.getUsernameAndTenderIdReqParams(ctx)
	if err != nil {
		return err
	}

	version, err := ctx.GetVersionPathParam()
	if err != nil {
		return err
	}

	if err := t.checkEmployeeTenderAccess(tenderId, username); err != nil {
		return err
	}

	tenderRollbackDbModel, err := t.tenderRollbackRepos.GetTenderRollback(tenderId, version)
	if err != nil {
		return err
	}

	tenderDbModel, err := t.tenderRepos.RollbackTender(tenderRollbackDbModel)
	if err != nil {
		return err
	}

	err = t.tenderRollbackRepos.SaveTenderRollback(tenderDbModel)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(200, models.NewTenderDtoModel(tenderDbModel))
}

func (t TenderHandler) getUsernameAndTenderIdReqParams(ctx *Context) (string, string, *errors.AppError) {
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

func (t TenderHandler) checkEmployeeTenderAccess(tenderId, username string) *errors.AppError {
	_, err := t.getTenderIfEmployeeResponsible(tenderId, username)
	if err != nil {
		return err
	}
	return nil
}

func (t TenderHandler) getTenderIfEmployeeResponsible(tenderId, username string) (*models.TenderDbModel, *errors.AppError) {
	tenderDbModel, err := t.tenderRepos.GetTenderById(tenderId)
	if err != nil {
		return nil, err
	}
	if err := t.checkEmployeeExistsAndResponsible(username, tenderDbModel.OrganizationId); err != nil {
		return nil, err
	}
	return tenderDbModel, nil
}

func (t TenderHandler) checkEmployeeExistsAndResponsible(username, organizationId string) *errors.AppError {
	if err := t.employeeRepos.CheckEmployeeExistsByUsername(username); err != nil {
		return err
	}
	if err := t.employeeRepos.CheckEmployeeIsResponsible(
		username,
		organizationId,
	); err != nil {
		return err
	}
	return nil
}
