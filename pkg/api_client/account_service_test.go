package api_client

import (
	"accountapi-lib-form3/pkg/configuration"
	"accountapi-lib-form3/pkg/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

const (
	RightJsonResponse = `{"data":{"attributes":{"account_classification":"Personal","account_matching_opt_out":false,"alternative_names":["Sam Holder"],"bank_id":"400302","bank_id_code":"GBDSC","base_currency":"GBP","bic":"NWBKGB42","country":"GB","joint_account":false,"name":["Samantha Holder"],"secondary_identification":"A1B2C3D4"},"created_on":"2021-10-15T03:19:57.796Z","id":"ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6","modified_on":"2021-10-15T03:19:57.796Z","organisation_id":"ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6","type":"accounts","version":0},"links":{"self":"/v1/organisation/accounts/ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6"}}`
	WrongJsonResponse = `{"error_message":"id is not a valid uuid"}`
	CreationRequest   = `{"data":{"id":"ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6","organisation_id":"ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6","type":"accounts","attributes":{"country":"GB","base_currency":"GBP","bank_id":"400302","bank_id_code":"GBDSC","customer_id":"234","bic":"NWBKGB42","name":["Samantha Holder"],"alternative_names":["Sam Holder"],"account_classification":"Personal","joint_account":false,"account_matching_opt_out":false,"secondary_identification":"A1B2C3D4"}}}`
	AccountId         = "ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6"
	RightPort         = "80"
)

type transportFake struct {
	respJson   string
	statusCode int
	isError    bool
}

func (t *transportFake) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	if t.isError {
		return nil, fmt.Errorf("fake error")
	}

	r := ioutil.NopCloser(bytes.NewReader([]byte(t.respJson)))
	return &http.Response{
		Body:       r,
		StatusCode: t.statusCode,
	}, nil
}

func getBuilder(resp string, statusCode int, isError bool, port string) configuration.Config {
	client := &http.Client{
		Transport: &transportFake{
			respJson:   resp,
			statusCode: statusCode,
			isError:    isError,
		},
	}

	builder := configuration.NewDefaultConfigBuilder()
	return builder.
		WithHost("fake").
		WithPort(port).
		WithHttpClient(client).
		Build()
}

func getStringStruct(data interface{}) string {
	json, _ := json.Marshal(data)
	return string(json)
}

func TestAccountService_ShouldReturnSuccessfulCreation(t *testing.T) {
	var input models.CreateRequest
	_ = json.Unmarshal([]byte(CreationRequest), &input)

	var responseObject models.ResponseObject
	_ = json.Unmarshal([]byte(RightJsonResponse), &responseObject)
	want := &models.CreateResponse{
		ResBody:    &responseObject,
		StatusCode: 201,
	}

	builder := getBuilder(RightJsonResponse, 201, false, RightPort)
	subject := NewAccountService(&builder)
	got, _ := subject.CreateAccount(&input)

	if !reflect.DeepEqual(*got, *want) {
		t.Errorf("wanted: %s\n got: %s", getStringStruct(want), getStringStruct(got))
	}
}

func TestAccountService_ShouldReturnFailedCreation(t *testing.T) {
	var input models.CreateRequest
	_ = json.Unmarshal([]byte(CreationRequest), &input)

	var responseObject models.ResponseError
	_ = json.Unmarshal([]byte(WrongJsonResponse), &responseObject)
	want := "Create: 409 - id is not a valid uuid"

	builder := getBuilder(WrongJsonResponse, 409, false, RightPort)
	subject := NewAccountService(&builder)
	_, got := subject.CreateAccount(&input)

	if got.Error() != want {
		t.Errorf("wanted: %s\n got: %s", want, got.Error())
	}
}

func TestAccountService_ShouldReturnSuccessfulDeletion(t *testing.T) {
	input := models.DeleteRequest{
		AccountId: AccountId,
		Version:   0,
	}

	want := &models.DeleteResponse{
		StatusCode: 204,
	}

	builder := getBuilder("", 204, false, RightPort)
	subject := NewAccountService(&builder)
	got, _ := subject.DeleteAccount(&input)

	if !reflect.DeepEqual(*got, *want) {
		t.Errorf("wanted: %s\n got: %s", getStringStruct(want), getStringStruct(got))
	}
}

func TestAccountService_ShouldReturnNoMessageOnDeletionNotFound(t *testing.T) {
	input := models.DeleteRequest{
		AccountId: AccountId,
		Version:   0,
	}

	want := "Delete: 404 - "

	builder := getBuilder("", 404, false, RightPort)
	subject := NewAccountService(&builder)
	_, got := subject.DeleteAccount(&input)

	if got.Error() != want {
		t.Errorf("wanted: %s\n got: %s", want, got.Error())
	}
}

func TestAccountService_ShouldReturnFailedDeletion(t *testing.T) {
	input := models.DeleteRequest{
		AccountId: AccountId,
		Version:   0,
	}

	want := "Delete: 400 - id is not a valid uuid"

	builder := getBuilder(WrongJsonResponse, 400, false, RightPort)
	subject := NewAccountService(&builder)
	_, got := subject.DeleteAccount(&input)

	if got.Error() != want {
		t.Errorf("wanted: %s\n got: %s", want, got.Error())
	}
}

func TestAccountService_ShouldReturnSuccessfulFetch(t *testing.T) {
	input := models.FetchRequest{
		AccountId: AccountId,
	}

	var responseObject models.ResponseObject
	_ = json.Unmarshal([]byte(RightJsonResponse), &responseObject)
	want := &models.FetchResponse{
		ResBody:    &responseObject,
		StatusCode: 200,
	}

	builder := getBuilder(RightJsonResponse, 200, false, RightPort)
	subject := NewAccountService(&builder)
	got, _ := subject.FetchAccount(&input)

	if !reflect.DeepEqual(*got, *want) {
		t.Errorf("wanted: %s\n got: %s", getStringStruct(want), getStringStruct(got))
	}
}

func TestAccountService_ShouldReturnFailedFetch(t *testing.T) {
	input := models.FetchRequest{
		AccountId: AccountId,
	}

	want := "Fetch: 404 - id is not a valid uuid"

	builder := getBuilder(WrongJsonResponse, 404, false, RightPort)
	subject := NewAccountService(&builder)
	_, got := subject.FetchAccount(&input)

	if got.Error() != want {
		t.Errorf("wanted: %s\n got: %s", want, got.Error())
	}
}

func TestAccountService_ShouldReturnFailureWhenCreatingRequest(t *testing.T) {

	want := "2 - failed creating request"
	builder := getBuilder("", 200, false, "80 ")
	subject := NewAccountService(&builder)
	var got error

	_, got = subject.CreateAccount(&models.CreateRequest{})

	if !strings.Contains(got.Error(), want) {
		t.Errorf("wanted: %s\n got: %s", want, got)
	}

	_, got = subject.DeleteAccount(&models.DeleteRequest{})
	if !strings.Contains(got.Error(), want) {
		t.Errorf("wanted string: %s\n message got: %s", want, got)
	}

	_, got = subject.FetchAccount(&models.FetchRequest{})
	if !strings.Contains(got.Error(), want) {
		t.Errorf("wanted string: %s\n message got: %s", want, got)
	}
}

func TestAccountService_ShouldReturnFailureWhenInvokingBackend(t *testing.T) {

	want := "3 - failed invoking backend"
	builder := getBuilder("", 200, true, "80")
	subject := NewAccountService(&builder)
	var got error

	_, got = subject.CreateAccount(&models.CreateRequest{})
	if !strings.Contains(got.Error(), want) {
		t.Errorf("wanted: %s\n got: %s", want, got)
	}

	_, got = subject.DeleteAccount(&models.DeleteRequest{})
	if !strings.Contains(got.Error(), want) {
		t.Errorf("wanted: %s\n got: %s", want, got)
	}

	_, got = subject.FetchAccount(&models.FetchRequest{})
	if !strings.Contains(got.Error(), want) {
		t.Errorf("wanted: %s\n got: %s", want, got)
	}
}

func TestAccountService_ShouldReturnFailureWhenDecodingResponse(t *testing.T) {

	dataTable := []struct {
		testName   string
		statusCode int
		want       string
	}{
		{"rightCreation", 201, "6 - failed decoding response"},
		{"wrongCreation", 409, "5 - failed decoding error response"},
		{"wrongDeletion", 400, "5 - failed decoding error response"},
		{"rightFetch", 200, "6 - failed decoding response"},
		{"wrongFetch", 404, "5 - failed decoding error response"},
	}

	for _, v := range dataTable {
		t.Run(v.testName, func(t *testing.T) {
			builder := getBuilder("EOF", v.statusCode, false, "80")
			subject := NewAccountService(&builder)
			var got error

			switch {
			case strings.Contains(v.testName, "Creation"):
				_, got = subject.CreateAccount(&models.CreateRequest{})
				if !strings.Contains(got.Error(), v.want) {
					t.Errorf("wanted string: %s\n message got: %s", v.want, got)
				}
			case strings.Contains(v.testName, "Deletion"):
				_, got = subject.DeleteAccount(&models.DeleteRequest{})
				if !strings.Contains(got.Error(), v.want) {
					t.Errorf("wanted string: %s\n message got: %s", v.want, got)
				}
			case strings.Contains(v.testName, "Fetch"):
				_, got = subject.FetchAccount(&models.FetchRequest{})
				if !strings.Contains(got.Error(), v.want) {
					t.Errorf("wanted string: %s\n message got: %s", v.want, got)
				}
			}
		})
	}
}
