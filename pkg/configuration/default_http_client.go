package configuration

import (
	"net/http"
	"time"
)

func NewDefaultHttpClient(timeout time.Duration, verboseLog bool) *http.Client {
	customTransport := http.DefaultTransport

	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: customTransport,
	}

	if verboseLog {
		httpClient.Transport = &loggingRoundTripper{
			defaultRoundTripper: customTransport,
		}
	}

	return httpClient
}
