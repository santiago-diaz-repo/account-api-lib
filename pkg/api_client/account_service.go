package api_client

import (
	"accountapi-lib-form3/pkg/configuration"
	"accountapi-lib-form3/pkg/error_handling"
	"accountapi-lib-form3/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AccountService struct {
	config *configuration.Config
}

const (
	AccountsPath             = "/organisation/accounts"
	DateHeader               = "Date"
	AcceptHeader             = "Accept"
	JsonAPIMediaType         = "application/vnd.api+json"
	ContentTypeHeader        = "Content-Type"
	ApplicationJson          = "application/json"
	CreateOperation          = "Create"
	DeleteOperation          = "Delete"
	FetchOperation           = "Fetch"
	CodeFailedMarshallingReq = 1
	MsgFailedMarshallingReq  = "failed marshalling request: "
	CodeFailedCreatingReq    = 2
	MsgFailedCreatingReq     = "failed creating request: "
	CodeFailedInvokingBack   = 3
	MsgFailedInvokingBack    = "failed invoking backend: "
	CodeFailedReadingRes     = 4
	MsgFailedReadingRes      = "failed reading response body: "
	CodeFailedDecodingErrRes = 5
	MsgFailedDecodingErrRes  = "failed decoding error response: "
	CodeFailedDecodingRes    = 6
	MsgFailedDecodingRes     = "failed decoding response: "
)

func NewAccountService(config *configuration.Config) AccountManagement {
	return &AccountService{
		config: config,
	}
}

// CreateAccount allows to create an account by passing around some information about it
func (a *AccountService) CreateAccount(reqModel *models.CreateRequest) (*models.CreateResponse, error) {

	inp, err := json.Marshal(reqModel)
	if err != nil {
		return nil, error_handling.NewAccountError(CreateOperation, CodeFailedMarshallingReq, MsgFailedMarshallingReq+err.Error())
	}

	inpReader := strings.NewReader(string(inp))

	endpoint := (*a.config).GetAPIBasePath() + AccountsPath

	request, err := http.NewRequest(http.MethodPost, endpoint, inpReader)
	if err != nil {
		return nil, error_handling.NewAccountError(CreateOperation, CodeFailedCreatingReq, MsgFailedCreatingReq+err.Error())
	}
	request.Header.Set(DateHeader, time.Now().Format(time.RFC3339))
	request.Header.Set(ContentTypeHeader, ApplicationJson)

	response, err := (*a.config).GetHttpClient().Do(request)
	if err != nil {
		return nil, error_handling.NewAccountError(CreateOperation, CodeFailedInvokingBack, MsgFailedInvokingBack+err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, error_handling.NewAccountError(CreateOperation, CodeFailedReadingRes, MsgFailedReadingRes+err.Error())
	}

	if response.StatusCode != http.StatusCreated {
		var outErr models.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, error_handling.NewAccountError(CreateOperation, CodeFailedDecodingErrRes, MsgFailedDecodingErrRes+err.Error())
		}

		return nil, error_handling.NewAccountError(CreateOperation, response.StatusCode, outErr.ErrorMessage)
	}

	var out models.ResponseObject
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, error_handling.NewAccountError(CreateOperation, CodeFailedDecodingRes, MsgFailedDecodingRes+err.Error())
	}

	return &models.CreateResponse{
		ResBody:    &out,
		StatusCode: response.StatusCode,
	}, nil
}

// DeleteAccount allows to delete a particular account by using its ID and version.
// As 404 error returns no message, it evaluates that statusCode to return an empty message to AccountError
func (a *AccountService) DeleteAccount(reqModel *models.DeleteRequest) (*models.DeleteResponse, error) {
	endpoint := fmt.Sprintf("%s%s/%s?version=%d", (*a.config).GetAPIBasePath(), AccountsPath, reqModel.AccountId, reqModel.Version)

	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, error_handling.NewAccountError(DeleteOperation, CodeFailedCreatingReq, MsgFailedCreatingReq+err.Error())
	}
	request.Header.Set(DateHeader, time.Now().Format(time.RFC3339))

	response, err := (*a.config).GetHttpClient().Do(request)
	if err != nil {
		return nil, error_handling.NewAccountError(DeleteOperation, CodeFailedInvokingBack, MsgFailedInvokingBack+err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, error_handling.NewAccountError(DeleteOperation, CodeFailedReadingRes, MsgFailedReadingRes+err.Error())
	}

	if response.StatusCode != http.StatusNoContent {

		if response.StatusCode == http.StatusNotFound {
			return nil, error_handling.NewAccountError(DeleteOperation, response.StatusCode, "")
		}

		var outErr models.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, error_handling.NewAccountError(DeleteOperation, CodeFailedDecodingErrRes, MsgFailedDecodingErrRes+err.Error())
		}

		return nil, error_handling.NewAccountError(DeleteOperation, response.StatusCode, outErr.ErrorMessage)
	}

	return &models.DeleteResponse{
		StatusCode: response.StatusCode,
	}, nil
}

// FetchAccount allows to get a particular account by searching for its ID
func (a *AccountService) FetchAccount(reqModel *models.FetchRequest) (*models.FetchResponse, error) {
	endpoint := fmt.Sprintf("%s%s/%s", (*a.config).GetAPIBasePath(), AccountsPath, reqModel.AccountId)

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, error_handling.NewAccountError(FetchOperation, CodeFailedCreatingReq, MsgFailedCreatingReq+err.Error())
	}
	request.Header.Set(DateHeader, time.Now().Format(time.RFC3339))
	request.Header.Set(AcceptHeader, JsonAPIMediaType)

	response, err := (*a.config).GetHttpClient().Do(request)
	if err != nil {
		return nil, error_handling.NewAccountError(FetchOperation, CodeFailedInvokingBack, MsgFailedInvokingBack+err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, error_handling.NewAccountError(FetchOperation, CodeFailedReadingRes, MsgFailedReadingRes+err.Error())
	}

	if response.StatusCode != http.StatusOK {
		var outErr models.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, error_handling.NewAccountError(FetchOperation, CodeFailedDecodingErrRes, MsgFailedDecodingErrRes+err.Error())
		}

		return nil, error_handling.NewAccountError(FetchOperation, response.StatusCode, outErr.ErrorMessage)
	}

	var out models.ResponseObject
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, error_handling.NewAccountError(CreateOperation, CodeFailedDecodingRes, MsgFailedDecodingRes+err.Error())
	}

	return &models.FetchResponse{
		ResBody:    &out,
		StatusCode: response.StatusCode,
	}, nil
}
