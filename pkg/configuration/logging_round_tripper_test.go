package configuration

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
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

	_, _ = subject.RoundTrip(&d)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func TestLoggingRoundTripper_ShouldPrintSomething(t *testing.T) {

	subject := loggingRoundTripper{
		defaultRoundTripper: &transportFake{},
	}

	got := len(stdOutPrinted(subject))

	if got == 0 {
		t.Errorf("wanted: value greater than 0 \n got: %d", got)
	}
}
