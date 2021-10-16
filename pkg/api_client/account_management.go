package api_client

import (
	models2 "accountapi-lib-form3/pkg/models"
)

type AccountManagement interface {
	CreateAccount(*models2.CreateRequest) (*models2.CreateResponse, error)
	DeleteAccount(*models2.DeleteRequest) (*models2.DeleteResponse, error)
	FetchAccount(*models2.FetchRequest) (*models2.FetchResponse, error)
}
