// +build integration

package integration_tests

import (
	"accountapi-lib-form3/pkg/api_client"
	"accountapi-lib-form3/pkg/configuration"
	"accountapi-lib-form3/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func Test_FetchAccount(t *testing.T) {
	var createInput models.CreateRequest
	json.Unmarshal([]byte(CreationRequest), &createInput)

	id := uuid.New().String()

	createInput.Data.ID = id
	createInput.Data.OrganisationID = id

	config := configuration.NewDefaultConfigBuilder().
		WithPort("8090").
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
			{"ID-does-not-exist", "1aa111aa-1111-1a1a-1111-1a1a1a1aa1a1", 404, fmt.Sprintf("record 1aa111aa-1111-1a1a-1111-1a1a1a1aa1a1 does not exist")},
		}

		for _, v := range dataTable {
			t.Run(v.TestName, func(t *testing.T) {

				input := models.FetchRequest{
					AccountId: v.ID,
				}

				res, err := subject.FetchAccount(&input)
				if err != nil {
					t.Errorf("There was a problem when executing test %s, error: %v", v.TestName, err)
				} else {

					if res.StatusCode != v.StatusCodeWanted {
						t.Errorf("statusCode wanted: %d\n statusCode got: %d", v.StatusCodeWanted, res.StatusCode)
					}

					if res.StatusCode == 200 {
						if res.ResBody.Data.ID != v.MessageWanted {
							t.Errorf("id wanted: %s\n id got: %s", v.MessageWanted, res.ResBody.Data.ID)
						}
					} else {

						if res.ErrorMessage != v.MessageWanted {
							t.Errorf("messge wanted: %s\n message got: %s", v.MessageWanted, res.ErrorMessage)
						}
					}
				}
			})
		}
	}
}
