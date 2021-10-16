package configuration

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type transportFake struct{}

func (t *transportFake) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	json := `{"test":"dummy response"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	return &http.Response{
		Body:   r,
		Status: "200 OK",
		Header: map[string][]string{"test": {"test header"}},
	}, nil
}

func stdOutPrinted(subject loggingRoundTripper) string {

	json := `{"test":"dummy request"}`
	jsonReq := io.NopCloser(strings.NewReader(json))
	getBody := func() (io.ReadCloser, error) {
		r := bytes.NewReader([]byte(json))
		return io.NopCloser(r), nil
	}

	d := http.Request{
		Method:  "POST",
		URL:     &url.URL{Host: "test"},
		Header:  map[string][]string{"test": {"test header"}},
		Body:    jsonReq,
		GetBody: getBody,
	}

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	subject.RoundTrip(&d)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func TestDefaultClient_ShouldReturnCustomTransport(t *testing.T) {
	want := "*configuration.loggingRoundTripper"
	subject := NewDefaultHttpClient(4*time.Second, false)
	got := reflect.TypeOf(subject.Transport).String()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}

func TestDefaultClient_ShouldReturnCustomTimeout(t *testing.T) {
	want := 2 * time.Second
	subject := NewDefaultHttpClient(want, false)
	got := subject.Timeout

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}

func TestDefaultClient_ShouldPrintSomething(t *testing.T) {

	subject := loggingRoundTripper{
		defaultRoundTripper: &transportFake{},
		verboseLog:          true,
	}

	got := len(stdOutPrinted(subject))

	if got == 0 {
		t.Errorf("wanted: value greater than 0 \n got: %d", got)
	}
}

func TestDefaultClient_ShouldPrintNothing(t *testing.T) {

	want := 0
	subject := loggingRoundTripper{
		defaultRoundTripper: &transportFake{},
		verboseLog:          false,
	}

	got := len(stdOutPrinted(subject))

	if got != want {
		t.Errorf("wanted: value greater than 0 \n got: %d", got)
	}
}
