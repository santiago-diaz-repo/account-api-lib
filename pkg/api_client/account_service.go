package api_client

import (
	configuration2 "accountapi-lib-form3/pkg/configuration"
	models2 "accountapi-lib-form3/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AccountService struct {
	config *configuration2.Config
}

const (
	AccountsPath      = "/organisation/accounts"
	DateHeader        = "Date"
	AcceptHeader      = "Accept"
	JsonAPIMediaType  = "application/vnd.api+json"
	ContentTypeHeader = "Content-Type"
	ApplicationJson   = "application/json"
)

func NewAccountService(config *configuration2.Config) AccountManagement {
	return &AccountService{
		config: config,
	}
}

func (a *AccountService) CreateAccount(reqModel *models2.CreateRequest) (*models2.CreateResponse, error) {

	inp, err := json.Marshal(reqModel)
	if err != nil {
		return nil, fmt.Errorf("Account create: Failed marshalling request: %v\n", err)
	}

	inpReader := strings.NewReader(string(inp))

	endpoint := (*a.config).APIBasePath() + AccountsPath

	request, err := http.NewRequest(http.MethodPost, endpoint, inpReader)
	if err != nil {
		return nil, fmt.Errorf("Account create: Failed creating request: %v\n", err)
	}
	request.Header.Set(DateHeader, time.Now().Format(time.RFC3339))
	request.Header.Set(ContentTypeHeader, ApplicationJson)

	response, err := (*a.config).HttpClient().Do(request)
	if err != nil {
		return nil, fmt.Errorf("Account create: Failed invoking: %v\n", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Account create: Failed reading response Body: %v\n", err)
	}

	if response.StatusCode != http.StatusCreated {
		var outErr models2.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, fmt.Errorf("Account create: Failed decoding error response: %v\n", err)
		}

		return &models2.CreateResponse{
			StatusCode:   response.StatusCode,
			ErrorMessage: outErr.ErrorMessage,
		}, nil
	}

	var out models2.ResponseObject
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("Account create: Failed decoding response: %v\n", err)
	}

	return &models2.CreateResponse{
		ResBody:    &out,
		StatusCode: response.StatusCode,
	}, nil
}

func (a *AccountService) DeleteAccount(reqModel *models2.DeleteRequest) (*models2.DeleteResponse, error) {
	endpoint := fmt.Sprintf("%s%s/%s?version=%d", (*a.config).APIBasePath(), AccountsPath, reqModel.AccountId, reqModel.Version)

	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("Account delete: Failed creating request: %v\n", err)
	}
	request.Header.Set(DateHeader, time.Now().Format(time.RFC3339))

	response, err := (*a.config).HttpClient().Do(request)
	if err != nil {
		return nil, fmt.Errorf("Account delete: Failed invoking: %v\n", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Account delete: Failed reading response Body: %v\n", err)
	}

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusNotFound {
		var outErr models2.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, fmt.Errorf("Account delete: Failed decoding error response: %v\n", err)
		}

		return &models2.DeleteResponse{
			StatusCode:   response.StatusCode,
			ErrorMessage: outErr.ErrorMessage,
		}, nil
	}

	return &models2.DeleteResponse{
		StatusCode: response.StatusCode,
	}, nil
}

// FetchAccount allows to get a particular account by searching for its ID
func (a *AccountService) FetchAccount(reqModel *models2.FetchRequest) (*models2.FetchResponse, error) {
	endpoint := fmt.Sprintf("%s%s/%s", (*a.config).APIBasePath(), AccountsPath, reqModel.AccountId)

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("Account fetch: Failed creating request: %v\n", err)
	}
	request.Header.Set(DateHeader, time.Now().Format(time.RFC3339))
	request.Header.Set(AcceptHeader, JsonAPIMediaType)

	response, err := (*a.config).HttpClient().Do(request)
	if err != nil {
		return nil, fmt.Errorf("Account fetch: Failed invoking backend: %v\n", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Account fetch: Failed reading response Body: %v\n", err)
	}

	if response.StatusCode != http.StatusOK {
		var outErr models2.ResponseError
		err = json.Unmarshal(body, &outErr)
		if err != nil {
			return nil, fmt.Errorf("Account fetch: Failed decoding error response: %v\n", err)
		}

		return &models2.FetchResponse{
			StatusCode:   response.StatusCode,
			ErrorMessage: outErr.ErrorMessage,
		}, nil
	}

	var out models2.ResponseObject
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("Account fetch: Failed decoding response: %v\n", err)
	}

	return &models2.FetchResponse{
		ResBody:    &out,
		StatusCode: response.StatusCode,
	}, nil
}
