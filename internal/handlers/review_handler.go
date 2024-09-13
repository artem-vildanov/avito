package handlers

import (
	"avito/internal/db"
	"avito/internal/errors"
	"avito/internal/models"
	"avito/internal/services"
	"net/http"
)

type ReviewHandler struct {
	reviewRepos     *db.ReviewRepository
	bidRepos        *db.BidRepository
	employeeService *services.EmployeeService
	tenderRepos     *db.TenderRepository
	employeeRepos   *db.EmployeeRepository
}

func NewReviewHandler(storage *db.PostgresStorage) *ReviewHandler {
	return &ReviewHandler{
		reviewRepos:     db.NewReviewRepository(storage),
		bidRepos:        db.NewBidRepository(storage),
		employeeService: services.NewEmployeeService(storage),
		employeeRepos:   db.NewEmployeeRepository(storage),
		tenderRepos:     db.NewTenderRepository(storage),
	}
}

func (self ReviewHandler) LeaveFeedback(ctx *Context) *errors.AppError {
	bidId, err := ctx.GetBidIdPathParam()
	if err != nil {
		return err
	}

	username, err := ctx.GetUsernameRequestParam()
	if err != nil {
		return err
	}

	bidFeedback, err := ctx.GetBidFeedbackRequestParam()
	if err != nil {
		return err
	}

	bidDbModel, err := self.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	tenderDbModel, err := self.tenderRepos.GetTenderById(bidDbModel.TenderId)
	if err != nil {
		return err
	}

	if err := self.employeeService.CheckEmployeeIsResponsible(username, tenderDbModel.OrganizationId); err != nil {
		return err
	}

	_, err = self.reviewRepos.NewReview(bidId, bidFeedback)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModel(bidDbModel))
}

func (self ReviewHandler) GetReviews(ctx *Context) *errors.AppError {
	limit, offset, err := ctx.GetLimitAndOffsetRequestParams()
	if err != nil {
		return err
	}

	tenderId, err := ctx.GetTenderIdPathParam()
	if err != nil {
		return err
	}

	authorUsername, err := ctx.GetAuthorUsernameRequestParam()
	if err != nil {
		return err
	}

	requesterUsername, err := ctx.GetRequesterUsernameRequestParam()
	if err != nil {
		return err
	}

	tenderDbModel, err := self.tenderRepos.GetTenderById(tenderId)
	if err != nil {
		return err
	}

	authorDbModel, err := self.employeeRepos.GetEmployeeByUsername(authorUsername)
	if err != nil {
		return err
	}

	_, err = self.employeeRepos.GetEmployeeByUsername(requesterUsername)
	if err != nil {
		return err
	}

	// проверка что запрашивающий ответственнен за оргу которая создала тендер
	if err := self.employeeRepos.CheckEmployeeIsResponsible(requesterUsername, tenderDbModel.OrganizationId); err != nil {
		return err
	}

	// проверить что автор имеет предложение для тендера
	if err := self.employeeRepos.CheckEmployeeHasBidsForTender(authorDbModel.Id, tenderId); err != nil {
		return err
	}

	// если имеет, то найти все другие предложения автора
	reviewsDbModels, err := self.reviewRepos.GetReviewsByBidsAuthorUsername(authorUsername, limit, offset)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewReviewDtoModelList(reviewsDbModels))
}
