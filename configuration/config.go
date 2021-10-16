package configuration

import (
	"fmt"
	"net/http"
)

type Config interface {
	APIBasePath() string
	HttpClient() *http.Client
}

type config struct {
	apiVersion string
	host       string
	port       string
	httpClient *http.Client
	verboseLog bool
}

const defaultScheme = "http"

func (c *config) APIBasePath() string {
	return fmt.Sprintf("%s://%s:%s/%s", defaultScheme, c.host, c.port, c.apiVersion)
}

func (c *config) HttpClient() *http.Client {
	return c.httpClient
}
