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

func (self Context) GetBidStatusRequestParam() (enums.BidStatus, *errors.AppError) {
	statusFromReq := self.Request.URL.Query().Get("status")
	if statusFromReq == "" {
		return enums.BidStatusCanceled, errors.RequiredRequestParamNotProvided("status")
	}
	if !utils.Contains(enums.BidStatusesList, statusFromReq) {
		return enums.BidStatusCanceled, errors.InvalidRequestParam(statusFromReq)
	}
	return enums.BidStatus(statusFromReq), nil
}

func (self Context) GetTenderStatusRequestParam() (enums.TenderStatus, *errors.AppError) {
	statusFromReq := self.Request.URL.Query().Get("status")
	if statusFromReq == "" {
		return enums.TenderStatusClosed, errors.RequiredRequestParamNotProvided("status")
	}
	if !utils.Contains(enums.TenderStatusesList, statusFromReq) {
		return enums.TenderStatusClosed, errors.InvalidRequestParam(statusFromReq)
	}
	return enums.TenderStatus(statusFromReq), nil
}

func (self Context) GetDecisionRequestParam() (enums.Decision, *errors.AppError) {
	decisionFromReq := self.Request.URL.Query().Get("decision")
	if decisionFromReq == "" {
		return enums.DecisionRejected, errors.RequiredRequestParamNotProvided("decision")
	}
	if !utils.Contains(enums.DecisionsList, decisionFromReq) {
		return enums.DecisionRejected, errors.InvalidRequestParam(decisionFromReq)
	}
	return enums.Decision(decisionFromReq), nil
}

func (self Context) GetServiceTypesRequestParam() ([]enums.ServiceType, *errors.AppError) {
	resultServiceTypes := []enums.ServiceType{}
	serviceTypesFromReq := self.Request.URL.Query()["service_type"]

	for _, reqServiceType := range serviceTypesFromReq {
		if utils.Contains(enums.ServiceTypesList, reqServiceType) {
			resultServiceTypes = append(resultServiceTypes, enums.ServiceType(reqServiceType))
		} else {
			return resultServiceTypes, errors.InvalidRequestParam(reqServiceType)
		}
	}

	return resultServiceTypes, nil
}

func (self Context) GetBidFeedbackRequestParam() (string, *errors.AppError) {
	bidFeedback := self.Request.URL.Query().Get("bidFeedback")
	if bidFeedback == "" {
		return "", errors.RequiredRequestParamNotProvided("bidFeedback")
	}
	return bidFeedback, nil
}

// GetUsernameRequestParam can be required
func (self Context) GetUsernameRequestParam() (string, *errors.AppError) {
	username := self.Request.URL.Query().Get("username")
	if username == "" {
		return "", errors.RequiredRequestParamNotProvided("username")
	}
	return username, nil
}

func (self Context) GetLimitAndOffsetRequestParams() (uint, uint, *errors.AppError) {
	limit, err := self.GetLimitRequestParam()
	if err != nil {
		return 0, 0, err
	}
	offset, err := self.GetOffsetRequestParam()
	if err != nil {
		return 0, 0, nil
	}
	return limit, offset, nil
}

func (self Context) GetLimitRequestParam() (uint, *errors.AppError) {
	limitStr := self.Request.URL.Query().Get("limit")
	if limitStr == "" {
		return LimitDefaultValue, nil
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt < 0 {
		return 0, errors.InvalidRequestParam(limitStr)
	}

	return uint(limitInt), nil
}

func (self Context) GetOffsetRequestParam() (uint, *errors.AppError) {
	offsetStr := self.Request.URL.Query().Get("offset")
	if offsetStr == "" {
		return OffsetDefaultValue, nil
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil || offsetInt < 0 {
		return 0, errors.InvalidRequestParam(offsetStr)
	}

	return uint(offsetInt), nil
}

func (self Context) GetVersionPathParam() (uint, *errors.AppError) {
	versionStr, err := self.GetPathParam("version")
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
func (self Context) GetBidIdPathParam() (string, *errors.AppError) {
	return self.GetPathParam("bidId")
}

// GetTenderIdPathParam always required
func (self Context) GetTenderIdPathParam() (string, *errors.AppError) {
	return self.GetPathParam("tenderId")
}

func (self Context) GetPathParam(paramKey string) (string, *errors.AppError) {
	paramValue, exists := self.pathParams[paramKey]
	if !exists {
		return "", errors.RequiredRequestParamNotProvided(paramKey)
	}

	return paramValue, nil
}

func (self Context) RespondWithJson(status int, content any) *errors.AppError {
	if err := self.writeJson(status, content); err != nil {
		errMessage := fmt.Sprintf("json parse error: %s", err.Error())
		log.Println(errMessage)
		_ = self.writeJson(http.StatusInternalServerError, errMessage)
	}
	return nil
}

func (self Context) writeJson(status int, content any) error {
	self.Response.Header().Set("Content-Type", "application/json")
	self.Response.WriteHeader(int(status))
	return json.NewEncoder(self.Response).Encode(content)
}
