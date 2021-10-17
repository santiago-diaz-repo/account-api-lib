// +build integration

package integration_tests

import (
	"accountapi-lib-form3/pkg/api_client"
	"accountapi-lib-form3/pkg/configuration"
	"accountapi-lib-form3/pkg/error_handling"
	"accountapi-lib-form3/pkg/models"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func Test_FetchAccount(t *testing.T) {
	var createInput models.CreateRequest
	json.Unmarshal([]byte(CreationRequest), &createInput)

	id := "e6f3eed0-3f37-416a-9cfc-187c2caadb69"

	createInput.Data.ID = id
	createInput.Data.OrganisationID = id

	config := configuration.NewDefaultConfigBuilder().
		WithPort("8080").
		WithHost("accountapi").
		Build()

	subject := api_client.NewAccountService(&config)
	_, err := subject.CreateAccount(&createInput)
	if err != nil {
		t.Errorf("Creating is required to test Fetch, but it has generated an error: %v", err)
	} else {

		dataTable := []struct {
			TestName         string
			ID               string
			StatusCodeWanted int
			MessageWanted    string
		}{
			{"fetch successful", id, 200, id},
			{"wrong-ID", "123", 400, "id is not a valid uuid"},
			{"ID-does-not-exist", "a82e431e-4087-4b64-961f-df6342f78c2f", 404, fmt.Sprintf("record a82e431e-4087-4b64-961f-df6342f78c2f does not exist")},
		}

		for _, v := range dataTable {
			t.Run(v.TestName, func(t *testing.T) {

				input := models.FetchRequest{
					AccountId: v.ID,
				}

				res, err := subject.FetchAccount(&input)
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
						t.Errorf("statusCode wanted: %d\n statusCode got: %d", v.StatusCodeWanted, res.StatusCode)
					}

					if res.ResBody.Data.ID != v.MessageWanted {
						t.Errorf("id wanted: %s\n id got: %s", v.MessageWanted, res.ResBody.Data.ID)
					}
				}
			})
		}
	}

	// Cleaning environment
	_,_ = subject.DeleteAccount(&models.DeleteRequest{
		AccountId: id,
		Version: 0,
	})
}
