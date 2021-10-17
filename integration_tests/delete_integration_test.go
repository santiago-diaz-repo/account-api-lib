// +build integration

package integration_tests

import (
	"accountapi-lib-form3/pkg/api_client"
	"accountapi-lib-form3/pkg/configuration"
	"accountapi-lib-form3/pkg/error_handling"
	"accountapi-lib-form3/pkg/models"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"testing"
)

func Test_DeleteAccount(t *testing.T) {
	var createInput models.CreateRequest
	json.Unmarshal([]byte(CreationRequest), &createInput)

	id := uuid.New().String()
	createInput.Data.ID = id
	createInput.Data.OrganisationID = id

	config := configuration.NewDefaultConfigBuilder().
		WithPort("8080").
		Build()

	subject := api_client.NewAccountService(&config)
	_, err := subject.CreateAccount(&createInput)
	if err != nil {
		t.Errorf("Creating is required to test Delete, but it has generated an error: %v", err)
	} else {

		dataTable := []struct {
			TestName         string
			ID               string
			version          int
			StatusCodeWanted int
			MessageWanted    string
		}{
			{"delete-conflict", id, 1, 409, "invalid version"},
			{"delete successful", id, 0, 204, ""},
			{"wrong-ID", "123", 0, 400, "id is not a valid uuid"},
			{"ID-does-not-exist", id, 0, 404, ""},
		}

		for _, v := range dataTable {
			t.Run(v.TestName, func(t *testing.T) {

				input := models.DeleteRequest{
					AccountId: v.ID,
					Version:   v.version,
				}

				res, err := subject.DeleteAccount(&input)
				if err != nil {
					acctErr := err.(*error_handling.AccountError)

					if acctErr.GetCode() != v.StatusCodeWanted {
						t.Errorf("wanted: %d\n got: %d", v.StatusCodeWanted, res.StatusCode)
					}

					if !strings.Contains(acctErr.GetMessage(),v.MessageWanted) {
						t.Errorf("wanted: %s\n got: %s", v.MessageWanted, acctErr.GetMessage())
					}
				} else {
					if res.StatusCode != v.StatusCodeWanted {
						t.Errorf("wanted: %d\n got: %d", v.StatusCodeWanted, res.StatusCode)
					}
				}
			})
		}
	}
}
