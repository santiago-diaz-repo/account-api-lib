package configuration

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type loggingRoundTripper struct {
	defaultRoundTripper http.RoundTripper
}

// RoundTrip prints information related to request and response. It is important to make sure to use
// this method to debug purposes as it reduces performance significantly.
func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	res, err := l.defaultRoundTripper.RoundTrip(req)

	fmt.Printf("[%s] Method: %s\n", time.Now().Format(time.RFC3339), req.Method)
	fmt.Printf("[%s] URL: %s\n", time.Now().Format(time.RFC3339), req.URL.String())
	fmt.Printf("[%s] Request Headers: %s\n", time.Now().Format(time.RFC3339), req.Header)

	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		getBody := req.GetBody
		copyBody, _ := getBody()
		bodyBytes, _ := ioutil.ReadAll(copyBody)
		fmt.Printf("[%s] Request Body: %s\n", time.Now().Format(time.RFC3339), bodyBytes)
	}

	if err == nil {
		fmt.Printf("[%s] StatusCode: %s\n", time.Now().Format(time.RFC3339), res.Status)
		fmt.Printf("[%s] Response Headers: %s\n", time.Now().Format(time.RFC3339), res.Header)
		bodyBytes, _ := io.ReadAll(res.Body)
		fmt.Printf("[%s] Response Body: %s\n", time.Now().Format(time.RFC3339), string(bodyBytes))
		res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return res, err
}
