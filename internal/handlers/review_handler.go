package handlers

import (
	"avito/internal/db"
	"avito/internal/errors"
	"avito/internal/models"
	"net/http"
)

type ReviewHandler struct {
	reviewRepos   *db.ReviewRepository
	bidRepos      *db.BidRepository
	tenderRepos   *db.TenderRepository
	employeeRepos *db.EmployeeRepository
}

func NewReviewHandler(storage *db.PostgresStorage) *ReviewHandler {
	return &ReviewHandler{
		reviewRepos:   db.NewReviewRepository(storage),
		bidRepos:      db.NewBidRepository(storage),
		employeeRepos: db.NewEmployeeRepository(storage),
		tenderRepos:   db.NewTenderRepository(storage),
	}
}

func (r ReviewHandler) LeaveFeedback(ctx *Context) *errors.AppError {
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

	bidDbModel, err := r.bidRepos.GetBidById(bidId)
	if err != nil {
		return err
	}

	tenderDbModel, err := r.tenderRepos.GetTenderById(bidDbModel.TenderId)
	if err != nil {
		return err
	}

	if err := r.employeeRepos.CheckEmployeeExistsByUsername(username); err != nil {
		return err
	}

	if err := r.employeeRepos.CheckEmployeeIsResponsible(username, tenderDbModel.OrganizationId); err != nil {
		return err
	}

	_, err = r.reviewRepos.NewReview(bidId, bidFeedback)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewBidDtoModel(bidDbModel))
}

func (r ReviewHandler) GetReviews(ctx *Context) *errors.AppError {
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

	tenderDbModel, err := r.tenderRepos.GetTenderById(tenderId)
	if err != nil {
		return err
	}

	authorDbModel, err := r.employeeRepos.GetEmployeeByUsername(authorUsername)
	if err != nil {
		return err
	}

	_, err = r.employeeRepos.GetEmployeeByUsername(requesterUsername)
	if err != nil {
		return err
	}

	// проверка что запрашивающий ответственнен за оргу которая создала тендер
	if err := r.employeeRepos.CheckEmployeeIsResponsible(requesterUsername, tenderDbModel.OrganizationId); err != nil {
		return err
	}

	// проверить что автор имеет предложение для тендера
	if err := r.employeeRepos.CheckEmployeeHasBidsForTender(authorDbModel.Id, tenderId); err != nil {
		return err
	}

	// если имеет, то найти все другие предложения автора
	reviewsDbModels, err := r.reviewRepos.GetReviewsByBidsAuthorUsername(authorUsername, limit, offset)
	if err != nil {
		return err
	}

	return ctx.RespondWithJson(http.StatusOK, models.NewReviewDtoModelList(reviewsDbModels))
}
