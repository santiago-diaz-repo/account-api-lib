// +build integration

package integration_tests

import (
	"accountapi-lib-form3/pkg/api_client"
	"accountapi-lib-form3/pkg/configuration"
	"accountapi-lib-form3/pkg/models"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"testing"
)

const (
	CreationRequest = `{"data":{"id":"change","organisation_id":"change","type":"accounts","attributes":{"country":"GB","base_currency":"GBP","bank_id":"400302","bank_id_code":"GBDSC","customer_id":"234","bic":"NWBKGB42","name":["Samantha Holder"],"alternative_names":["Sam Holder"],"account_classification":"Personal","joint_account":false,"account_matching_opt_out":false,"secondary_identification":"A1B2C3D4"}}}`
)

func Test_CreateAccount(t *testing.T) {
	var input models.CreateRequest
	json.Unmarshal([]byte(CreationRequest), &input)

	id := uuid.New().String()

	config := configuration.NewConfigBuilder().
		WithPort("8090").
		Build("localhost")

	subject := api_client.NewAccountService(&config)

	dataTable := []struct {
		TestName         string
		ID               string
		OrganisationId   string
		StatusCodeWanted int
		MessageWanted    string
	}{
		{"creation-successful", id, id, 201, id},
		{"wrong-ID", "123", id, 400, "id in body must be of type uuid"},
		{"wrong-organisation-ID", id, "123", 400, "organisation_id in body must be of type uuid"},
		{"id-already-exists", id, id, 409, "Account cannot be created as it violates a duplicate constraint"},
	}

	for _, v := range dataTable {
		t.Run(v.TestName, func(t *testing.T) {
			input.Data.ID = v.ID
			input.Data.OrganisationID = v.OrganisationId
			res, err := subject.CreateAccount(&input)
			if err != nil {
				t.Errorf("There was a problem when executing test %s, error: %v", v.TestName, err)
			} else {

				if res.StatusCode != v.StatusCodeWanted {
					t.Errorf("statusCode wanted: %d\n statusCode got: %d", v.StatusCodeWanted, res.StatusCode)
				}

				if res.StatusCode == 201 {
					if res.ErrorMessage != "" {
						t.Errorf("there shouldn't be an error message \n got: %s", res.ErrorMessage)
					}

					if res.ResBody.Data.OrganisationID != v.MessageWanted {
						t.Errorf("organisationId wanted: %s\n organisationId got: %s", v.MessageWanted, res.ResBody.Data.OrganisationID)
					}

				} else {
					if !strings.Contains(res.ErrorMessage, v.MessageWanted) {
						t.Errorf("wanted: %s\n got: %s", v.MessageWanted, res.ErrorMessage)
					}
				}
			}
		})
	}
}
