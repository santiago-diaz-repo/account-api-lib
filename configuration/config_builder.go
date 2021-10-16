package configuration

import (
	"net/http"
	"time"
)

type ConfigBuilder interface {
	WithHttpClient(*http.Client) ConfigBuilder
	WithAPIVersion(string) ConfigBuilder
	WithPort(string) ConfigBuilder
	Verbose() ConfigBuilder
	Build(string) Config
}

type configBuilderStruct struct {
	config
}

const (
	DefaultAPIVersion = "v1"
	DefaultPort       = "80"
	DefaultTimeout    = 4 * time.Second
	DefaultVerbose    = false
)

func NewConfigBuilder() ConfigBuilder {
	configBuilder := new(configBuilderStruct)
	configBuilder.config.port = DefaultPort
	configBuilder.config.apiVersion = DefaultAPIVersion
	configBuilder.config.verboseLog = DefaultVerbose
	return configBuilder
}

func (c *configBuilderStruct) WithHttpClient(httpClient *http.Client) ConfigBuilder {
	c.config.httpClient = httpClient
	return c
}

func (c *configBuilderStruct) WithAPIVersion(apiVersion string) ConfigBuilder {
	c.config.apiVersion = apiVersion
	return c
}

func (c *configBuilderStruct) WithPort(port string) ConfigBuilder {
	c.config.port = port
	return c
}

func (c *configBuilderStruct) Verbose() ConfigBuilder {
	c.config.verboseLog = true
	return c
}

func (c *configBuilderStruct) Build(host string) Config {
	if c.config.httpClient == nil {
		c.config.httpClient = NewDefaultHttpClient(DefaultTimeout, c.verboseLog)
	}
	//c.config.httpClient.Transport

	c.config.host = host

	return &c.config
}
