package handlers

import (
	"avito/internal/db"
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/models"
	"net/http"
)

type BidHandler struct {
	bidRepos          *db.BidRepository
	bidRollbackRepos  *db.BidRollbackRepository
	tenderRepos       *db.TenderRepository
	orgRepos          *db.OrganizationRepository
	employeeRepos     *db.EmployeeRepository
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
		responsiblesRepos: db.NewResponsibleRepository(storage),
		bidApproveRepos:   db.NewBidApproveRepository(storage),
	}
}

func (b BidHandler) CreateBid(ctx *Context) *errors.AppError {
	bidCreateModel, err := GetModelFromRequest[models.BidCreateModel](ctx.Request.Body)
	if err != nil {
		return err
	}

	tenderDbModel, err := b.tenderRepos.GetTenderById(bidCreateModel.TenderId)
	if err != nil {
		return err
	}

	if bidCreateModel.AuthorType == enums.AuthorTypeOrganization {
		orgDbModel, err := b.orgRepos.GetOrganizationById(bidCreateModel.AuthorId)
		if err != nil {
			return err
		}

		if tenderDbModel.Status != enums.TenderStatusPublished &&
			tenderDbModel.OrganizationId != orgDbModel.Id {
			return errors.TenderNotPublished(tenderDbModel.Id)
		}
	} else if bidCreateModel.AuthorType == enums.AuthorTypeUser {
		employeeDbModel, err := b.employeeRepos.GetEmployeeById(bidCreateModel.AuthorId)
		if err != nil {
			return err
		}

		if tenderDbModel.Status != enums.TenderStatusPublished {
			if err := b.employeeRepos.CheckEmployeeIsResponsible(employeeDbModel.Username, tenderDbModel.OrganizationId); err != nil {
				return err
			}
		}
	}

	bidDbModel, err := b.bidRepos.CreateBid(bidCreateModel)
	if err != nil {
		return err
	}

	if err := b.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModel(bidDbModel))
}

func (b BidHandler) GetBidsListByUsername(ctx *Context) *errors.AppError {
	limit, offset, err := ctx.GetLimitAndOffsetRequestParams()
	if err != nil {
		return err
	}

	username, err := ctx.GetUsernameRequestParam()
	if err != nil {
		return err
	}

	if err := b.employeeRepos.CheckEmployeeExistsByUsername(username); err != nil {
		return err
	}

	bidDbModelsList, err := b.bidRepos.GetBidsListByUsername(username, limit, offset)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModelsList(bidDbModelsList))
}

func (b BidHandler) GetBidsListByTenderId(ctx *Context) *errors.AppError {
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

	if err := b.checkGetTenderAccess(tenderId, username); err != nil {
		return err
	}

	bidsDbModelsList, err := b.bidRepos.GetBidsListByTenderId(tenderId, limit, offset)
	if err != nil {
		return err
	}

	userDbModel, err := b.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	var bidsDtoModelsList []*models.BidDtoModel
	for _, bidModel := range bidsDbModelsList {
		if b.checkGetBidAccess(bidModel, userDbModel) {
			bidsDtoModelsList = append(bidsDtoModelsList, models.NewBidDtoModel(bidModel))
		}
	}

	return ctx.RespondWithJson(http.StatusOK, bidsDtoModelsList)
}

func (b BidHandler) GetBidStatus(ctx *Context) *errors.AppError {
	bidId, username, err := b.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	bidDbModel, err := b.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	employeeDbModel, err := b.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !b.checkGetBidAccess(bidDbModel, employeeDbModel) {
		return errors.NotEnoughPermissions(username)
	}

	return ctx.RespondWithJson(http.StatusOK, bidDbModel.Status)
}

func (b BidHandler) UpdateBidStatus(ctx *Context) *errors.AppError {
	bidId, username, err := b.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	status, err := ctx.GetBidStatusRequestParam()
	if err != nil {
		return err
	}

	bidDbModel, err := b.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	employeeDbModel, err := b.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !b.checkUpdateBidAccess(bidDbModel, employeeDbModel) {
		return errors.NotEnoughPermissions(username)
	}

	bidDbModel, err = b.bidRepos.UpdateBidStatus(bidId, status)
	if err != nil {
		return err
	}

	if err := b.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, bidDbModel.Status)
}

func (b BidHandler) UpdateBidParams(ctx *Context) *errors.AppError {
	bidId, username, err := b.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	bidUpdateModel, err := GetModelFromRequest[models.BidUpdateModel](ctx.Request.Body)
	if err != nil {
		return err
	}

	bidDbModel, err := b.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	userDbModel, err := b.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !b.checkUpdateBidAccess(bidDbModel, userDbModel) {
		return errors.NotEnoughPermissions(username)
	}

	bidDbModel, err = b.bidRepos.UpdateBidParams(bidId, bidUpdateModel)
	if err != nil {
		return err
	}

	if err := b.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModel(bidDbModel))
}

func (b BidHandler) SubmitDecision(ctx *Context) *errors.AppError {
	bidId, username, err := b.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	decision, err := ctx.GetDecisionRequestParam()
	if err != nil {
		return err
	}

	bidDbModel, err := b.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	employeeDbModel, err := b.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !b.checkGetBidAccess(bidDbModel, employeeDbModel) {
		return errors.BidNotPublished(bidId)
	}

	tenderDbModel, err := b.tenderRepos.GetTenderById(bidDbModel.TenderId)
	if err != nil {
		return err
	}

	if err := b.employeeRepos.CheckEmployeeIsResponsible(
		username,
		tenderDbModel.OrganizationId,
	); err != nil {
		return err
	}

	var bidDtoModel *models.BidDtoModel
	if decision == enums.DecisionApproved {
		if err := b.bidApproveRepos.CheckEmployeeApprovedBid(bidId, employeeDbModel.Id); err != nil {
			return err
		}

		if err := b.bidApproveRepos.AddApprove(bidId, employeeDbModel.Id); err != nil {
			return err
		}

		approvementsCount, err := b.bidApproveRepos.CountApprovementsByBidId(bidId)
		if err != nil {
			return err
		}

		responsiblesCount, err := b.responsiblesRepos.CountResponsiblesByOrgId(tenderDbModel.OrganizationId)
		if err != nil {
			return err
		}

		// если проходит условие кворума
		if responsiblesCount < 3 && approvementsCount == responsiblesCount ||
			responsiblesCount >= 3 && approvementsCount >= 3 {

			_, err := b.tenderRepos.UpdateTenderStatus(tenderDbModel.Id, enums.TenderStatusClosed)
			if err != nil {
				return err
			}

			bidDbModel, err := b.bidRepos.UpdateBidStatus(bidId, enums.BidStatusApproved)
			if err != nil {
				return err
			}

			bidsDbModelsList, err := b.bidRepos.CancelBidsByTenderId(bidDbModel.TenderId)
			if err != nil {
				return err
			}

			if err := b.bidRollbackRepos.SaveBidRollbacksList(bidsDbModelsList); err != nil {
				return err
			}

			if err := b.bidApproveRepos.RemoveApprovesByBidId(bidId); err != nil {
				return err
			}

			bidDtoModel = models.NewBidDtoModel(bidDbModel)
		} else {
			bidDtoModel = models.NewBidDtoModel(bidDbModel)
		}
	} else if decision == enums.DecisionRejected {
		bidDbModel, err := b.bidRepos.UpdateBidStatus(bidId, enums.BidStatusCanceled)
		if err != nil {
			return err
		}

		if err := b.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
			return err
		}

		if err := b.bidApproveRepos.RemoveApprovesByBidId(bidId); err != nil {
			return err
		}

		bidDtoModel = models.NewBidDtoModel(bidDbModel)
	}

	return ctx.RespondWithJson(http.StatusOK, bidDtoModel)
}

func (b BidHandler) RollbackBid(ctx *Context) *errors.AppError {
	bidId, username, err := b.getBidIdAndUsernameFromReq(ctx)
	if err != nil {
		return err
	}

	version, err := ctx.GetVersionPathParam()
	if err != nil {
		return err
	}

	bidDbModel, err := b.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	employeeDbModel, err := b.employeeRepos.GetEmployeeByUsername(username)
	if err != nil {
		return err
	}

	if !b.checkUpdateBidAccess(bidDbModel, employeeDbModel) {
		return errors.NotEnoughPermissions(username)
	}

	bidRollbackDbModel, err := b.bidRollbackRepos.GetBidRollback(bidId, version)
	if err != nil {
		return err
	}

	bidDbModel, err = b.bidRepos.RollbackBid(bidRollbackDbModel)
	if err != nil {
		return err
	}

	if err := b.bidRollbackRepos.SaveBidRollback(bidDbModel); err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModel(bidDbModel))
}

func (b BidHandler) getBidIdAndUsernameFromReq(ctx *Context) (string, string, *errors.AppError) {
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

func (b BidHandler) checkGetBidAccess(
	bidDbModel *models.BidDbModel,
	employeeDbModel *models.EmployeeDbModel,
) bool {
	if bidDbModel.Status == enums.BidStatusCreated || bidDbModel.Status == enums.BidStatusCanceled {
		return b.checkUpdateBidAccess(bidDbModel, employeeDbModel)
	}
	return true
}

func (b BidHandler) checkUpdateBidAccess(
	bidDbModel *models.BidDbModel,
	employeeDbModel *models.EmployeeDbModel,
) bool {
	if bidDbModel.AuthorType == enums.AuthorTypeOrganization {
		if err := b.employeeRepos.CheckEmployeeIsResponsible(employeeDbModel.Username, bidDbModel.AuthorId); err != nil {
			return false
		}
	} else if bidDbModel.AuthorType == enums.AuthorTypeUser && bidDbModel.AuthorId != employeeDbModel.Id {
		return false
	}

	return true
}

func (b BidHandler) checkGetTenderAccess(tenderId, username string) *errors.AppError {
	tenderDbModel, err := b.tenderRepos.GetTenderById(tenderId)
	if err != nil {
		return err
	}

	if tenderDbModel.Status != enums.TenderStatusPublished {
		if err := b.employeeRepos.CheckEmployeeIsResponsible(username, tenderDbModel.OrganizationId); err != nil {
			return err
		}
	}

	return nil
}
