package configuration

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

type customTransportFake struct{}

func (t *customTransportFake) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return nil, nil
}

func TestConfigBuilder_ShouldAssignNewHttpClient(t *testing.T) {
	subject := configBuilderStruct{
		config{
			httpClient: &http.Client{},
		},
	}

	want := &http.Client{}
	subject.WithHttpClient(want)
	got := subject.config.httpClient

	if got != want {
		t.Errorf("wanted: %p\n got: %p", want, got)
	}
}

func TestConfigBuilder_ShouldAssignNewAPIVersion(t *testing.T) {
	subject := configBuilderStruct{
		config{
			apiVersion: "v1",
		},
	}

	want := "v2"
	subject.WithAPIVersion(want)
	got := subject.config.apiVersion

	if got != want {
		t.Errorf("wanted: %s\n got: %s", want, got)
	}
}

func TestConfigBuilder_ShouldAssignNewPort(t *testing.T) {
	subject := configBuilderStruct{
		config{
			port: "80",
		},
	}

	want := "8080"
	subject.WithPort(want)
	got := subject.config.port

	if got != want {
		t.Errorf("wanted: %s\n got: %s", want, got)
	}
}

func TestConfigBuilder_ShouldEnableVerbose(t *testing.T) {
	subject := configBuilderStruct{
		config{
			verboseLog: false,
		},
	}

	want := true
	subject.Verbose()
	got := subject.config.verboseLog

	if got != want {
		t.Errorf("wanted: %t\n got: %t", want, got)
	}
}

func TestConfigBuilder_ShouldReturnBasicConfigImplementation(t *testing.T) {
	want := &config{
		host:       "localhost",
		apiVersion: "v1",
		port:       "80",
		verboseLog: false,
		httpClient: &http.Client{Timeout: 4 * time.Second},
	}

	subject := NewConfigBuilder().
		Build("localhost")
	valueType := reflect.ValueOf(subject)
	got := valueType.Interface().(*config)

	if got.host != want.host {
		t.Errorf("host wanted: %s\n host got: %s", want.host, got.host)
	}

	if got.apiVersion != want.apiVersion {
		t.Errorf("apiVersion wanted: %s\n apiVersion got: %s", want.apiVersion, got.apiVersion)
	}

	if got.verboseLog != want.verboseLog {
		t.Errorf("verboseLog wanted: %t\n verboseLog got: %t", want.verboseLog, got.verboseLog)
	}

	if got.port != want.port {
		t.Errorf("port wanted: %s\n port got: %s", want.port, got.port)
	}

	if got.httpClient.Timeout != want.httpClient.Timeout {
		t.Errorf("client timeout wanted: %s\n client timeout got: %s", want.httpClient.Timeout, got.httpClient.Timeout)
	}
}

func TestConfigBuilder_ShouldReturnVerboseConfigImplementation(t *testing.T) {
	want := &config{
		host:       "test",
		apiVersion: "v2",
		port:       "8080",
		verboseLog: true,
		httpClient: &http.Client{},
	}

	subject := NewConfigBuilder().
		WithHttpClient(want.httpClient).
		WithPort("8080").
		WithAPIVersion("v2").
		Verbose().
		Build("test")

	valueType := reflect.ValueOf(subject)
	got := valueType.Interface().(*config)

	if !reflect.DeepEqual(*got, *want) {
		t.Errorf("wanted: %#v\n got: %#v", want, got)
	}

	transportGot := reflect.TypeOf(got.httpClient.Transport).String()

	if transportGot != "*configuration.loggingRoundTripper" {
		t.Errorf("transport wanted: *configuration.loggingRoundTripper \n transport got: %v", transportGot)
	}
}

func TestConfigBuilder_ShouldReturnNonVerboseConfigImplementation(t *testing.T) {
	want := &config{
		host:       "test",
		apiVersion: "v2",
		port:       "8080",
		verboseLog: false,
		httpClient: &http.Client{Transport: http.DefaultTransport},
	}

	subject := NewConfigBuilder().
		WithHttpClient(want.httpClient).
		WithPort("8080").
		WithAPIVersion("v2").
		Build("test")

	valueType := reflect.ValueOf(subject)
	got := valueType.Interface().(*config)

	if !reflect.DeepEqual(*got, *want) {
		t.Errorf("wanted: %#v\n got: %#v", want, got)
	}

	transportGot := reflect.TypeOf(got.httpClient.Transport).String()

	if transportGot != "*http.Transport" {
		t.Errorf("transport wanted: *http.Transport \n transport got: %v", transportGot)
	}
}

func TestConfigBuilder_ShouldSetVerboseLoggingNonNilTransport(t *testing.T) {
	want := "*configuration.loggingRoundTripper"
	subject := &http.Client{Transport: &customTransportFake{}}
	setVerboseLogging(subject)
	got := reflect.TypeOf(subject.Transport).String()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}

	internalTransport := reflect.ValueOf(subject.Transport).Interface().(*loggingRoundTripper)
	internalTransportGot := reflect.TypeOf(internalTransport.defaultRoundTripper).String()

	if internalTransportGot != "*configuration.customTransportFake" {
		t.Errorf("wanted: *configuration.customTransportFake\n got: %v", internalTransportGot)
	}
}

func TestConfigBuilder_ShouldSetVerboseLoggingNilTransport(t *testing.T) {
	want := "*configuration.loggingRoundTripper"
	subject := &http.Client{}
	setVerboseLogging(subject)
	got := reflect.TypeOf(subject.Transport).String()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}

	internalTransport := reflect.ValueOf(subject.Transport).Interface().(*loggingRoundTripper)
	internalTransportGot := reflect.TypeOf(internalTransport.defaultRoundTripper).String()

	if internalTransportGot != "*http.Transport" {
		t.Errorf("wanted: *http.Transport\n got: %v", internalTransportGot)
	}
}
