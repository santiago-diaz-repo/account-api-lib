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
	accountsPath             = "/organisation/accounts"
	dateHeader               = "Date"
	acceptHeader             = "Accept"
	jsonAPIMediaType         = "application/vnd.api+json"
	contentTypeHeader        = "Content-Type"
	applicationJson          = "application/json"
	createOperation          = "Create"
	deleteOperation          = "Delete"
	fetchOperation           = "Fetch"
	codeFailedMarshallingReq = 1
	msgFailedMarshallingReq  = "failed marshalling request: "
	codeFailedCreatingReq    = 2
	msgFailedCreatingReq     = "failed creating request: "
	codeFailedInvokingBack   = 3
	msgFailedInvokingBack    = "failed invoking backend: "
	codeFailedReadingRes     = 4
	msgFailedReadingRes      = "failed reading response body: "
	codeFailedDecodingErrRes = 5
	msgFailedDecodingErrRes  = "failed decoding error response: "
	codeFailedDecodingRes    = 6
	msgFailedDecodingRes     = "failed decoding response: "
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
		return nil, error_handling.NewAccountError(createOperation, codeFailedMarshallingReq, msgFailedMarshallingReq+err.Error())
	}

	inpReader := strings.NewReader(string(inp))

	endpoint := (*a.config).GetAPIBasePath() + accountsPath

	request, err := http.NewRequest(http.MethodPost, endpoint, inpReader)
	if err != nil {
		return nil, error_handling.NewAccountError(createOperation, codeFailedCreatingReq, msgFailedCreatingReq+err.Error())
	}
	request.Header.Set(dateHeader, time.Now().Format(time.RFC3339))
	request.Header.Set(contentTypeHeader, applicationJson)

	response, err := (*a.config).GetHttpClient().Do(request)
	if err != nil {
		return nil, error_handling.NewAccountError(createOperation, codeFailedInvokingBack, msgFailedInvokingBack+err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, error_handling.NewAccountError(createOperation, codeFailedReadingRes, msgFailedReadingRes+err.Error())
	}

	if response.StatusCode != http.StatusCreated {
		var outErr models.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, error_handling.NewAccountError(createOperation, codeFailedDecodingErrRes, msgFailedDecodingErrRes+err.Error())
		}

		return nil, error_handling.NewAccountError(createOperation, response.StatusCode, outErr.ErrorMessage)
	}

	var out models.ResponseObject
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, error_handling.NewAccountError(createOperation, codeFailedDecodingRes, msgFailedDecodingRes+err.Error())
	}

	return &models.CreateResponse{
		ResBody:    &out,
		StatusCode: response.StatusCode,
	}, nil
}

// DeleteAccount allows to delete a particular account by using its ID and version.
// As 404 error returns no message, it evaluates that statusCode to return an empty message to AccountError
func (a *AccountService) DeleteAccount(reqModel *models.DeleteRequest) (*models.DeleteResponse, error) {
	endpoint := fmt.Sprintf("%s%s/%s?version=%d", (*a.config).GetAPIBasePath(), accountsPath, reqModel.AccountId, reqModel.Version)

	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, error_handling.NewAccountError(deleteOperation, codeFailedCreatingReq, msgFailedCreatingReq+err.Error())
	}
	request.Header.Set(dateHeader, time.Now().Format(time.RFC3339))

	response, err := (*a.config).GetHttpClient().Do(request)
	if err != nil {
		return nil, error_handling.NewAccountError(deleteOperation, codeFailedInvokingBack, msgFailedInvokingBack+err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, error_handling.NewAccountError(deleteOperation, codeFailedReadingRes, msgFailedReadingRes+err.Error())
	}

	if response.StatusCode != http.StatusNoContent {

		if response.StatusCode == http.StatusNotFound {
			return nil, error_handling.NewAccountError(deleteOperation, response.StatusCode, "")
		}

		var outErr models.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, error_handling.NewAccountError(deleteOperation, codeFailedDecodingErrRes, msgFailedDecodingErrRes+err.Error())
		}

		return nil, error_handling.NewAccountError(deleteOperation, response.StatusCode, outErr.ErrorMessage)
	}

	return &models.DeleteResponse{
		StatusCode: response.StatusCode,
	}, nil
}

// FetchAccount allows to get a particular account by searching for its ID
func (a *AccountService) FetchAccount(reqModel *models.FetchRequest) (*models.FetchResponse, error) {
	endpoint := fmt.Sprintf("%s%s/%s", (*a.config).GetAPIBasePath(), accountsPath, reqModel.AccountId)

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, error_handling.NewAccountError(fetchOperation, codeFailedCreatingReq, msgFailedCreatingReq+err.Error())
	}
	request.Header.Set(dateHeader, time.Now().Format(time.RFC3339))
	request.Header.Set(acceptHeader, jsonAPIMediaType)

	response, err := (*a.config).GetHttpClient().Do(request)
	if err != nil {
		return nil, error_handling.NewAccountError(fetchOperation, codeFailedInvokingBack, msgFailedInvokingBack+err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, error_handling.NewAccountError(fetchOperation, codeFailedReadingRes, msgFailedReadingRes+err.Error())
	}

	if response.StatusCode != http.StatusOK {
		var outErr models.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, error_handling.NewAccountError(fetchOperation, codeFailedDecodingErrRes, msgFailedDecodingErrRes+err.Error())
		}

		return nil, error_handling.NewAccountError(fetchOperation, response.StatusCode, outErr.ErrorMessage)
	}

	var out models.ResponseObject
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, error_handling.NewAccountError(createOperation, codeFailedDecodingRes, msgFailedDecodingRes+err.Error())
	}

	return &models.FetchResponse{
		ResBody:    &out,
		StatusCode: response.StatusCode,
	}, nil
}
