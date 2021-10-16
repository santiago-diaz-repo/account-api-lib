package configuration

import (
	"net/http"
	"time"
)

func NewDefaultHttpClient(timeout time.Duration, verboseLog bool) *http.Client {
	customTransport := http.DefaultTransport

	return &http.Client{
		Timeout: timeout,
		Transport: &loggingRoundTripper{
			defaultRoundTripper: customTransport,
			verboseLog:          verboseLog,
		},
	}
}
