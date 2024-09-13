package handlers

import (
	"avito/internal/db"
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/models"
	"avito/internal/services"
	"net/http"
)

type BidHandler struct {
	bidRepos          *db.BidRepository
	bidRollbackRepos  *db.BidRollbackRepository
	tenderRepos       *db.TenderRepository
	orgRepos          *db.OrganizationRepository
	employeeRepos     *db.EmployeeRepository
	employeeService   *services.EmployeeService
	responsiblesRepos *db.ResponsibleRepository
	bidApproveRepos   *db.BidApproveRepository
}

func NewBidHandler(storage *db.PostgresStorage) *BidHandler {
	return &BidHandler{
		bidRepos:          db.NewBidRepository(storage),
		bidRollbackRepos:  db.NewBidRollbackRepository(storage),
		tenderRepos:       db.NewTenderRepository(storage),
		orgRepos:          db.NewOrganizationRepository(storage),
		employeeRepos:     db.NewEmployeeRepository(storage),
		employeeService:   services.NewEmployeeService(storage),
		responsiblesRepos: db.NewResponsibleRepository(storage),
		bidApproveRepos:   db.NewBidApproveRepository(storage),
	}
}

func (self BidHandler) CreateBid(ctx *Context) *errors.AppError {
	bidCreateModel, err := GetModelFromRequest[models.BidCreateModel](ctx.Request.Body)
	if err != nil {
		return err
	}

	tenderDbModel, err := self.tenderRepos.GetTenderById(bidCreateModel.TenderId)
	if err != nil {
		return err
	}

	if bidCreateModel.AuthorType == enums.AuthorTypeOrganization {
		orgDbModel, err := self.orgRepos.GetOrganizationById(bidCreateModel.AuthorId)
		if err != nil {
			return err
		}

		if tenderDbModel.Status != enums.TenderStatusPublished &&
			tenderDbModel.OrganizationId != orgDbModel.Id {
			return errors.TenderNotPublished(tenderDbModel.Id)
		}
	} else if bidCreateModel.AuthorType == enums.AuthorTypeUser {
		employeeDbModel, err := self.employeeRepos.GetEmployeeById(bidCreateModel.AuthorId)
		if err != nil {
			return err
		}

		if tenderDbModel.Status != enums.TenderStatusPublished {
			if err := self.employeeRepos.CheckEmployeeIsResponsible(employeeDbModel.Username, tenderDbModel.OrganizationId); err != nil {
				return err
			}
		}
	}

	bidDbModel, err := self.bidRepos.CreateBid(bidCreateModel)
	if err != nil {
		return err
	}

	if err := self.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModel(bidDbModel))
}

func (self BidHandler) GetBidsListByUsername(ctx *Context) *errors.AppError {
	limit, offset, err := ctx.GetLimitAndOffsetRequestParams()
	if err != nil {
		return err
	}

	username, err := ctx.GetUsernameRequestParam()
	if err != nil {
		return err
	}

	if err := self.employeeRepos.CheckEmployeeExistsByUsername(username); err != nil {
		return err
	}

	bidDbModelsList, err := self.bidRepos.GetBidsListByUsername(username, limit, offset)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModelsList(bidDbModelsList))
}

func (self BidHandler) GetBidsListByTenderId(ctx *Context) *errors.AppError {
	limit, offset, err := ctx.GetLimitAndOffsetRequestParams()
	if err != nil {
		return err
	}

	tenderId, err := ctx.GetTenderIdPathParam()
	if err != nil {
		return err
	}

	username, err := ctx.GetUsernameRequestParam()
	if err != nil {
		return err
	}

	if err := self.checkGetTenderAccess(tenderId, username); err != nil {
		return err
	}

	bidsDbModelsList, err := self.bidRepos.GetBidsListByTenderId(tenderId, limit, offset)
	if err != nil {
		return err
	}

	userDbModel, err := self.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	var bidsDtoModelsList []*models.BidDtoModel
	for _, bidModel := range bidsDbModelsList {
		if self.checkGetBidAccess(bidModel, userDbModel) {
			bidsDtoModelsList = append(bidsDtoModelsList, models.NewBidDtoModel(bidModel))
		}
	}

	return ctx.RespondWithJson(http.StatusOK, bidsDtoModelsList)
}

func (self BidHandler) GetBidStatus(ctx *Context) *errors.AppError {
	bidId, username, err := self.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	bidDbModel, err := self.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	employeeDbModel, err := self.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !self.checkGetBidAccess(bidDbModel, employeeDbModel) {
		return errors.NotEnoughPermissions(username)
	}

	return ctx.RespondWithJson(http.StatusOK, bidDbModel.Status)
}

func (self BidHandler) UpdateBidStatus(ctx *Context) *errors.AppError {
	bidId, username, err := self.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	status, err := ctx.GetBidStatusRequestParam()
	if err != nil {
		return err
	}

	bidDbModel, err := self.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	employeeDbModel, err := self.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !self.checkUpdateBidAccess(bidDbModel, employeeDbModel) {
		return errors.NotEnoughPermissions(username)
	}

	bidDbModel, err = self.bidRepos.UpdateBidStatus(bidId, status)
	if err != nil {
		return err
	}

	if err := self.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, bidDbModel.Status)
}

func (self BidHandler) UpdateBidParams(ctx *Context) *errors.AppError {
	bidId, username, err := self.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	bidUpdateModel, err := GetModelFromRequest[models.BidUpdateModel](ctx.Request.Body)
	if err != nil {
		return err
	}

	bidDbModel, err := self.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	userDbModel, err := self.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !self.checkUpdateBidAccess(bidDbModel, userDbModel) {
		return errors.NotEnoughPermissions(username)
	}

	bidDbModel, err = self.bidRepos.UpdateBidParams(bidId, bidUpdateModel)
	if err != nil {
		return err
	}

	if err := self.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModel(bidDbModel))
}

func (self BidHandler) SubmitDecision(ctx *Context) *errors.AppError {
	bidId, username, err := self.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	decision, err := ctx.GetDecisionRequestParam()
	if err != nil {
		return err
	}

	bidDbModel, err := self.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	employeeDbModel, err := self.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	tenderDbModel, err := self.tenderRepos.GetTenderById(bidDbModel.TenderId)
	if err != nil {
		return err
	}

	if err := self.employeeRepos.CheckEmployeeIsResponsible(
		username,
		tenderDbModel.OrganizationId,
	); err != nil {
		return err
	}

	var bidDtoModel *models.BidDtoModel
	if decision == enums.DecisionApproved {
		if err := self.bidApproveRepos.CheckEmployeeApprovedBid(bidId, employeeDbModel.Id); err != nil {
			return err
		}

		if err := self.bidApproveRepos.AddApprove(bidId, employeeDbModel.Id); err != nil {
			return err
		}

		approvementsCount, err := self.bidApproveRepos.CountApprovementsByBidId(bidId)
		if err != nil {
			return err
		}

		responsiblesCount, err := self.responsiblesRepos.CountResponsiblesByOrgId(tenderDbModel.OrganizationId)
		if err != nil {
			return err
		}

		// если проходит условие кворума
		if responsiblesCount < 3 && approvementsCount == responsiblesCount ||
			responsiblesCount >= 3 && approvementsCount >= 3 {

			_, err := self.tenderRepos.UpdateTenderStatus(tenderDbModel.Id, enums.TenderStatusClosed)
			if err != nil {
				return err
			}

			bidDbModel, err := self.bidRepos.UpdateBidStatus(bidId, enums.BidStatusApproved)
			if err != nil {
				return err
			}

			bidsDbModelsList, err := self.bidRepos.CancelBidsByTenderId(bidDbModel.TenderId)
			if err != nil {
				return err
			}

			if err := self.bidRollbackRepos.SaveBidRollbacksList(bidsDbModelsList); err != nil {
				return err
			}

			if err := self.bidApproveRepos.RemoveApprovesByBidId(bidId); err != nil {
				return err
			}

			bidDtoModel = models.NewBidDtoModel(bidDbModel)
		} else {
			bidDtoModel = models.NewBidDtoModel(bidDbModel)
		}
	} else if decision == enums.DecisionRejected {
		bidDbModel, err := self.bidRepos.UpdateBidStatus(bidId, enums.BidStatusCanceled)
		if err != nil {
			return err
		}

		if err := self.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
			return err
		}

		if err := self.bidApproveRepos.RemoveApprovesByBidId(bidId); err != nil {
			return err
		}

		bidDtoModel = models.NewBidDtoModel(bidDbModel)
	}

	return ctx.RespondWithJson(http.StatusOK, bidDtoModel)
}

func (self BidHandler) RollbackBid(ctx *Context) *errors.AppError {
	bidId, username, err := self.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	version, err := ctx.GetVersionPathParam()
	if err != nil {
		return err
	}

	bidDbModel, err := self.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	employeeDbModel, err := self.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !self.checkUpdateBidAccess(bidDbModel, employeeDbModel) {
		return errors.NotEnoughPermissions(username)
	}

	bidRollbackDbModel, err := self.bidRollbackRepos.GetBidRollback(bidId, version)
	if err != nil {
		return err
	}

	bidDbModel, err = self.bidRepos.RollbackBid(bidRollbackDbModel)
	if err != nil {
		return err
	}

	if err := self.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModel(bidDbModel))
}

func (self BidHandler) getBidIdAndUsernameFromReq(ctx *Context) (string, string, *errors.AppError) {
	bidId, err := ctx.GetBidIdPathParam()
	if err != nil {
		return "", "", err
	}

	username, err := ctx.GetUsernameRequestParam()
	if err != nil {
		return "", "", err
	}

	return bidId, username, nil
}

func (self BidHandler) checkGetBidAccess(
	bidDbModel *models.BidDbModel,
	employeeDbModel *models.EmployeeDbModel,
) bool {
	if bidDbModel.Status == enums.BidStatusCreated || bidDbModel.Status == enums.BidStatusCanceled {
		return self.checkUpdateBidAccess(bidDbModel, employeeDbModel)
	}
	return true
}

func (self BidHandler) checkUpdateBidAccess(
	bidDbModel *models.BidDbModel,
	employeeDbModel *models.EmployeeDbModel,
) bool {
	if bidDbModel.AuthorType == enums.AuthorTypeOrganization {
		if err := self.employeeRepos.CheckEmployeeIsResponsible(employeeDbModel.Username, bidDbModel.AuthorId); err != nil {
			return false
		}
	} else if bidDbModel.AuthorType == enums.AuthorTypeUser && bidDbModel.AuthorId != employeeDbModel.Id {
		return false
	}

	return true
}

func (self BidHandler) checkGetTenderAccess(tenderId, username string) *errors.AppError {
	tenderDbModel, err := self.tenderRepos.GetTenderById(tenderId)
	if err != nil {
		return err
	}

	if tenderDbModel.Status != enums.TenderStatusPublished {
		if err := self.employeeRepos.CheckEmployeeIsResponsible(username, tenderDbModel.OrganizationId); err != nil {
			return err
		}
	}

	return nil
}
