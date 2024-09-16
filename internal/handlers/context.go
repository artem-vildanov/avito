package handlers

import (
	"avito/internal/enums"
	"avito/internal/errors"
	"avito/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

const LimitDefaultValue uint = 5
const OffsetDefaultValue uint = 0

func GetModelFromRequest[T interface{}](requestBody io.ReadCloser) (*T, *errors.AppError) {
	var model T
	if err := json.NewDecoder(requestBody).Decode(&model); err != nil {
		return nil, errors.InvalidRequestBody()
	}
	return &model, nil
}

type Context struct {
	Response   http.ResponseWriter
	Request    *http.Request
	pathParams map[string]string
}

func NewContext(responseWriter http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Response:   responseWriter,
		Request:    request,
		pathParams: mux.Vars(request),
	}
}

func (c Context) GetBidStatusRequestParam() (enums.BidStatus, *errors.AppError) {
	statusFromReq := enums.BidStatus(c.Request.URL.Query().Get("status"))
	if statusFromReq == "" {
		return enums.BidStatusCanceled, errors.RequiredRequestParamNotProvided("status")
	}

	if !utils.Contains(enums.GetBidStatuses(), statusFromReq) {
		return enums.BidStatusCanceled, errors.InvalidRequestParam(string(statusFromReq))
	}

	return statusFromReq, nil
}

func (c Context) GetTenderStatusRequestParam() (enums.TenderStatus, *errors.AppError) {
	statusFromReq := enums.TenderStatus(c.Request.URL.Query().Get("status"))
	if statusFromReq == "" {
		return enums.TenderStatusClosed, errors.RequiredRequestParamNotProvided("status")
	}

	if !utils.Contains(enums.GetTenderStatuses(), statusFromReq) {
		return enums.TenderStatusClosed, errors.InvalidRequestParam(string(statusFromReq))
	}
	return statusFromReq, nil
}

func (c Context) GetDecisionRequestParam() (enums.Decision, *errors.AppError) {
	decisionFromReq := enums.Decision(c.Request.URL.Query().Get("decision"))
	if decisionFromReq == "" {
		return enums.DecisionRejected, errors.RequiredRequestParamNotProvided("decision")
	}
	if !utils.Contains(enums.GetDecisions(), decisionFromReq) {
		return enums.DecisionRejected, errors.InvalidRequestParam(string(decisionFromReq))
	}
	return decisionFromReq, nil
}

func (c Context) GetServiceTypesRequestParam() ([]enums.ServiceType, *errors.AppError) {
	resultServiceTypes := []enums.ServiceType{}
	serviceTypesFromReq := c.Request.URL.Query()["service_type"]

	for _, reqServiceTypeStr := range serviceTypesFromReq {
		reqServiceType := enums.ServiceType(reqServiceTypeStr)
		if utils.Contains(enums.GetServiceTypes(), reqServiceType) {
			resultServiceTypes = append(resultServiceTypes, reqServiceType)
		} else {
			return resultServiceTypes, errors.InvalidRequestParam(string(reqServiceType))
		}
	}

	return resultServiceTypes, nil
}

func (c Context) GetRequesterUsernameRequestParam() (string, *errors.AppError) {
	requesterUsername := c.Request.URL.Query().Get("requesterUsername")
	if requesterUsername == "" {
		return "", errors.RequiredRequestParamNotProvided("requesterUsername")
	}
	return requesterUsername, nil
}

func (c Context) GetAuthorUsernameRequestParam() (string, *errors.AppError) {
	authorUsername := c.Request.URL.Query().Get("authorUsername")
	if authorUsername == "" {
		return "", errors.RequiredRequestParamNotProvided("authorUsername")
	}
	return authorUsername, nil
}

func (c Context) GetBidFeedbackRequestParam() (string, *errors.AppError) {
	bidFeedback := c.Request.URL.Query().Get("bidFeedback")
	if bidFeedback == "" {
		return "", errors.RequiredRequestParamNotProvided("bidFeedback")
	}
	return bidFeedback, nil
}

// GetUsernameRequestParam can be required
func (c Context) GetUsernameRequestParam() (string, *errors.AppError) {
	username := c.Request.URL.Query().Get("username")
	if username == "" {
		return "", errors.RequiredRequestParamNotProvided("username")
	}
	return username, nil
}

func (c Context) GetLimitAndOffsetRequestParams() (uint, uint, *errors.AppError) {
	limit, err := c.GetLimitRequestParam()
	if err != nil {
		return 0, 0, err
	}
	offset, err := c.GetOffsetRequestParam()
	if err != nil {
		return 0, 0, nil
	}
	return limit, offset, nil
}

func (c Context) GetLimitRequestParam() (uint, *errors.AppError) {
	limitStr := c.Request.URL.Query().Get("limit")
	if limitStr == "" {
		return LimitDefaultValue, nil
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt < 0 {
		return 0, errors.InvalidRequestParam(limitStr)
	}

	return uint(limitInt), nil
}

func (c Context) GetOffsetRequestParam() (uint, *errors.AppError) {
	offsetStr := c.Request.URL.Query().Get("offset")
	if offsetStr == "" {
		return OffsetDefaultValue, nil
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil || offsetInt < 0 {
		return 0, errors.InvalidRequestParam(offsetStr)
	}

	return uint(offsetInt), nil
}

func (c Context) GetVersionPathParam() (uint, *errors.AppError) {
	versionStr, err := c.GetPathParam("version")
	if err != nil {
		return 0, err
	}
	versionInt, convErr := strconv.Atoi(versionStr)
	if convErr != nil || versionInt < 0 {
		return 0, errors.InvalidRequestParam(versionStr)
	}
	return uint(versionInt), nil
}

// GetBidIdPathParam always required
func (c Context) GetBidIdPathParam() (string, *errors.AppError) {
	return c.GetPathParam("bidId")
}

// GetTenderIdPathParam always required
func (c Context) GetTenderIdPathParam() (string, *errors.AppError) {
	return c.GetPathParam("tenderId")
}

func (c Context) GetPathParam(paramKey string) (string, *errors.AppError) {
	paramValue, exists := c.pathParams[paramKey]
	if !exists {
		return "", errors.RequiredRequestParamNotProvided(paramKey)
	}

	return paramValue, nil
}

func (c Context) RespondWithJson(status int, content any) *errors.AppError {
	if err := c.writeJson(status, content); err != nil {
		errMessage := fmt.Sprintf("json parse error: %s", err.Error())
		log.Println(errMessage)
		_ = c.writeJson(http.StatusInternalServerError, errMessage)
	}
	return nil
}

func (c Context) writeJson(status int, content any) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(int(status))
	return json.NewEncoder(c.Response).Encode(content)
}
