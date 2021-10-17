package configuration

import (
	"fmt"
	"net/http"
)

type Config interface {
	GetAPIBasePath() string
	GetHttpClient() *http.Client
}

type config struct {
	apiVersion string
	host       string
	port       string
	httpClient *http.Client
	verboseLog bool
}

// defaultScheme can be changed when service consumption has to be through another protocol such as secure http (https)
const defaultScheme = "http"

func (c *config) GetAPIBasePath() string {
	return fmt.Sprintf("%s://%s:%s/%s", defaultScheme, c.host, c.port, c.apiVersion)
}

func (c *config) GetHttpClient() *http.Client {
	return c.httpClient
}
