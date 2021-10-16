package api_client

import "accountapi-lib-form3/models"

type AccountManagement interface {
	CreateAccount(*models.CreateRequest) (*models.CreateResponse, error)
	DeleteAccount(*models.DeleteRequest) (*models.DeleteResponse, error)
	FetchAccount(*models.FetchRequest) (*models.FetchResponse, error)
}
