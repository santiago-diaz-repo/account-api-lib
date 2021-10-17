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

// NewConfigBuilder creates a new default configuration, everything can be
// modified by other methods defined in this file. In order to modify timeout
// it is possible to change the default http.Client.
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

// Build returns a new configuration to invoke backend API, it is important to clarify that
// if Build receives a particular http.Client implementation and verbose logging is enabled,
// this will modify http.Client.Transport to set verbose logging up. Additionally, if http.Client.Transport
// is nil, this will assign a http.DefaultTransport.
//
// It is important to clarify that a component which uses this library has to pass around the host
// where the backend API is located.
func (c *configBuilderStruct) Build(host string) Config {
	if c.config.httpClient == nil {
		c.config.httpClient = NewDefaultHttpClient(DefaultTimeout, c.verboseLog)
	} else {
		if c.verboseLog {
			setVerboseLogging(c.httpClient)
		}
	}

	c.config.host = host

	return &c.config
}

// setVerboseLogging modifies an http.Client by adding a loggingRoundTripper to Transport,
// it is worth mentioning that this modification is based on the decorator pattern.
func setVerboseLogging(httpClient *http.Client) {
	transport := httpClient.Transport

	if transport == nil {
		transport = http.DefaultTransport
	}

	httpClient.Transport = &loggingRoundTripper{
		defaultRoundTripper: transport,
	}
}
