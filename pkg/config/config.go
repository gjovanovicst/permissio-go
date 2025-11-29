// Package config provides configuration types and utilities for the Permissio.io SDK.
package config

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	// DefaultAPIURL is the default Permissio.io API URL.
	DefaultAPIURL = "https://api.permissio.io"

	// DefaultTimeout is the default request timeout.
	DefaultTimeout = 30 * time.Second

	// DefaultRetryAttempts is the default number of retry attempts.
	DefaultRetryAttempts = 3

	// APIKeyPrefix is the expected prefix for API keys.
	APIKeyPrefix = "permis_key_"
)

// Config represents the SDK configuration.
type Config struct {
	// Token is the API key for authentication (required).
	Token string

	// ApiURL is the base URL for the Permissio.io API.
	ApiURL string

	// ProjectID is the project identifier.
	ProjectID string

	// EnvironmentID is the environment identifier.
	EnvironmentID string

	// Timeout is the request timeout duration.
	Timeout time.Duration

	// Debug enables debug logging.
	Debug bool

	// RetryAttempts is the number of retry attempts for failed requests.
	RetryAttempts int

	// ThrowOnError determines if errors should cause panics (default: false).
	ThrowOnError bool

	// CustomHeaders are additional headers to include in requests.
	CustomHeaders map[string]string

	// Logger is the optional zap logger for debug output.
	Logger *zap.Logger

	// HTTPClient is the optional custom HTTP client.
	HTTPClient *http.Client
}

// HasScope returns true if both ProjectID and EnvironmentID are set.
func (c *Config) HasScope() bool {
	return c.ProjectID != "" && c.EnvironmentID != ""
}

// UpdateScope updates the ProjectID and EnvironmentID.
func (c *Config) UpdateScope(projectID, environmentID string) {
	c.ProjectID = projectID
	c.EnvironmentID = environmentID
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New("API token is required")
	}

	if !strings.HasPrefix(c.Token, APIKeyPrefix) {
		return errors.New("invalid API key format: must start with '" + APIKeyPrefix + "'")
	}

	if c.ApiURL == "" {
		return errors.New("API URL is required")
	}

	if c.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}

	if c.RetryAttempts < 0 {
		return errors.New("retry attempts must be non-negative")
	}

	return nil
}

// ConfigBuilder provides a fluent interface for building Config.
type ConfigBuilder struct {
	config *Config
}

// NewConfigBuilder creates a new ConfigBuilder with the given API token.
func NewConfigBuilder(token string) *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			Token:         token,
			ApiURL:        DefaultAPIURL,
			Timeout:       DefaultTimeout,
			RetryAttempts: DefaultRetryAttempts,
			Debug:         false,
			ThrowOnError:  false,
			CustomHeaders: make(map[string]string),
		},
	}
}

// WithApiUrl sets the API URL.
func (b *ConfigBuilder) WithApiUrl(url string) *ConfigBuilder {
	b.config.ApiURL = url
	return b
}

// WithProjectID sets the project ID.
func (b *ConfigBuilder) WithProjectID(projectID string) *ConfigBuilder {
	b.config.ProjectID = projectID
	return b
}

// WithEnvironmentID sets the environment ID.
func (b *ConfigBuilder) WithEnvironmentID(environmentID string) *ConfigBuilder {
	b.config.EnvironmentID = environmentID
	return b
}

// WithTimeout sets the request timeout.
func (b *ConfigBuilder) WithTimeout(timeout time.Duration) *ConfigBuilder {
	b.config.Timeout = timeout
	return b
}

// WithDebug enables or disables debug logging.
func (b *ConfigBuilder) WithDebug(debug bool) *ConfigBuilder {
	b.config.Debug = debug
	return b
}

// WithRetryAttempts sets the number of retry attempts.
func (b *ConfigBuilder) WithRetryAttempts(attempts int) *ConfigBuilder {
	b.config.RetryAttempts = attempts
	return b
}

// WithThrowOnError sets whether errors should cause panics.
func (b *ConfigBuilder) WithThrowOnError(throwOnError bool) *ConfigBuilder {
	b.config.ThrowOnError = throwOnError
	return b
}

// WithCustomHeader adds a custom header.
func (b *ConfigBuilder) WithCustomHeader(key, value string) *ConfigBuilder {
	b.config.CustomHeaders[key] = value
	return b
}

// WithCustomHeaders sets multiple custom headers.
func (b *ConfigBuilder) WithCustomHeaders(headers map[string]string) *ConfigBuilder {
	for k, v := range headers {
		b.config.CustomHeaders[k] = v
	}
	return b
}

// WithLogger sets the zap logger.
func (b *ConfigBuilder) WithLogger(logger *zap.Logger) *ConfigBuilder {
	b.config.Logger = logger
	return b
}

// WithHTTPClient sets the custom HTTP client.
func (b *ConfigBuilder) WithHTTPClient(client *http.Client) *ConfigBuilder {
	b.config.HTTPClient = client
	return b
}

// Build returns the built configuration.
// It applies default values but does not validate.
func (b *ConfigBuilder) Build() *Config {
	// Ensure HTTP client is set
	if b.config.HTTPClient == nil {
		b.config.HTTPClient = &http.Client{
			Timeout: b.config.Timeout,
		}
	}

	return b.config
}

// BuildWithValidation returns the built configuration after validation.
func (b *ConfigBuilder) BuildWithValidation() (*Config, error) {
	config := b.Build()
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}
