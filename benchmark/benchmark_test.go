package benchmark

import (
	api_client2 "accountapi-lib-form3/pkg/api_client"
	configuration2 "accountapi-lib-form3/pkg/configuration"
	models2 "accountapi-lib-form3/pkg/models"
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type TransportFake struct {
}

func (t *TransportFake) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	json := `{"data":{"attributes":{"account_classification":"Personal","account_matching_opt_out":false,"alternative_names":["Sam Holder"],"bank_id":"400302","bank_id_code":"GBDSC","base_currency":"GBP","bic":"NWBKGB42","country":"GB","joint_account":false,"name":["Samantha Holder"],"secondary_identification":"A1B2C3D4"},"created_on":"2021-10-15T03:19:57.796Z","id":"ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6","modified_on":"2021-10-15T03:19:57.796Z","organisation_id":"ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6","type":"accounts","version":0},"links":{"self":"/v1/organisation/accounts/ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6"}}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{
		Body: r,
	}, nil
}

func BenchmarkWithoutVerbose(b *testing.B) {
	req := models2.FetchRequest{
		AccountId: "ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6",
	}

	client := &http.Client{
		Transport: &TransportFake{},
	}

	subject := configuration2.NewConfigBuilder().
		WithHttpClient(client).
		Build("fake")

	account := api_client2.NewAccountService(&subject)

	for i := 0; i < b.N; i++ {
		_, _ = account.FetchAccount(&req)
	}
}

// BenchmarkWithVerbose allows to determine impact that enabling verbose log has in comparison to
// executing API invocation without verbose log enabled.
// According to comparison between this method and BenchmarkWithoutVerbose, it is possible to say that
// enabling verbose log may reduce performance up to 90%, so verbose log should be used only to debug
// purposes.
func BenchmarkWithVerbose(b *testing.B) {
	req := models2.FetchRequest{
		AccountId: "ebb084cb-5cb7-49b5-b61c-ea0f7036e4b6",
	}

	client := &http.Client{
		Transport: &TransportFake{},
	}

	subject := configuration2.NewConfigBuilder().
		WithHttpClient(client).
		Verbose().
		Build("fake")

	account := api_client2.NewAccountService(&subject)

	for i := 0; i < b.N; i++ {
		_, _ = account.FetchAccount(&req)
	}
}
