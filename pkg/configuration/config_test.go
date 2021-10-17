package configuration

import (
	"net/http"
	"testing"
)

func getConfigStub(client *http.Client) *config {
	return &config{
		port:       "80",
		verboseLog: false,
		apiVersion: "v1",
		host:       "test",
		httpClient: client,
	}
}

func TestConfig_ShouldReturnAPIBasePath(t *testing.T) {
	want := "http://test:80/v1"
	subject := getConfigStub(nil)
	got := subject.GetAPIBasePath()

	if got != want {
		t.Errorf("wanted: %s\n got: %s", want, got)
	}
}

func TestConfig_ShouldReturnHttpClient(t *testing.T) {
	want := &http.Client{}
	subject := getConfigStub(want)
	got := subject.GetHttpClient()
	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}
