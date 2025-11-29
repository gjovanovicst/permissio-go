// Package api provides API client implementations for the Permissio.io SDK.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/permissio/permissio-go/pkg/config"
	"go.uber.org/zap"
)

// BaseClient provides common HTTP functionality for API clients.
type BaseClient struct {
	config *config.Config
}

// NewBaseClient creates a new BaseClient.
func NewBaseClient(cfg *config.Config) *BaseClient {
	return &BaseClient{
		config: cfg,
	}
}

// Config returns the configuration.
func (c *BaseClient) Config() *config.Config {
	return c.config
}

// BuildURL builds a URL for the given path.
func (c *BaseClient) BuildURL(path string) string {
	return fmt.Sprintf("%s%s", c.config.ApiURL, path)
}

// BuildFactsURL builds a URL for facts endpoints.
func (c *BaseClient) BuildFactsURL(path string) string {
	if c.config.HasScope() {
		return fmt.Sprintf("%s/v1/facts/%s/%s%s",
			c.config.ApiURL,
			c.config.ProjectID,
			c.config.EnvironmentID,
			path)
	}
	return fmt.Sprintf("%s/v1%s", c.config.ApiURL, path)
}

// BuildSchemaURL builds a URL for schema endpoints.
func (c *BaseClient) BuildSchemaURL(path string) string {
	if c.config.HasScope() {
		return fmt.Sprintf("%s/v1/schema/%s/%s%s",
			c.config.ApiURL,
			c.config.ProjectID,
			c.config.EnvironmentID,
			path)
	}
	return fmt.Sprintf("%s/v1%s", c.config.ApiURL, path)
}

// Request performs an HTTP request with retry logic.
func (c *BaseClient) Request(ctx context.Context, method, url string, body interface{}, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= c.config.RetryAttempts; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := c.doRequest(ctx, method, url, body, result)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on certain errors
		if apiErr, ok := err.(*PermisError); ok {
			if apiErr.StatusCode >= 400 && apiErr.StatusCode < 500 {
				return err // Don't retry client errors
			}
		}

		if c.config.Debug && c.config.Logger != nil {
			c.config.Logger.Debug("Request failed, retrying",
				zap.Int("attempt", attempt+1),
				zap.Error(err))
		}
	}

	return lastErr
}

// doRequest performs a single HTTP request.
func (c *BaseClient) doRequest(ctx context.Context, method, url string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.Token)

	// Add custom headers
	for key, value := range c.config.CustomHeaders {
		req.Header.Set(key, value)
	}

	if c.config.Debug && c.config.Logger != nil {
		c.config.Logger.Debug("Making request",
			zap.String("method", method),
			zap.String("url", url))
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.config.Debug && c.config.Logger != nil {
		c.config.Logger.Debug("Received response",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(respBody)))
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return c.parseError(resp.StatusCode, respBody)
	}

	// Parse result
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// parseError parses an error response.
func (c *BaseClient) parseError(statusCode int, body []byte) error {
	var errResp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
		Code    string `json:"code"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return &PermisError{
			Message:    string(body),
			StatusCode: statusCode,
		}
	}

	message := errResp.Message
	if message == "" {
		message = errResp.Error
	}
	if message == "" {
		message = http.StatusText(statusCode)
	}

	return &PermisError{
		Message:    message,
		Code:       errResp.Code,
		StatusCode: statusCode,
	}
}

// Get performs a GET request.
func (c *BaseClient) Get(ctx context.Context, url string, result interface{}) error {
	return c.Request(ctx, http.MethodGet, url, nil, result)
}

// Post performs a POST request.
func (c *BaseClient) Post(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPost, url, body, result)
}

// Put performs a PUT request.
func (c *BaseClient) Put(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPut, url, body, result)
}

// Patch performs a PATCH request.
func (c *BaseClient) Patch(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodPatch, url, body, result)
}

// Delete performs a DELETE request.
func (c *BaseClient) Delete(ctx context.Context, url string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, url, nil, result)
}

// DeleteWithBody performs a DELETE request with a body.
func (c *BaseClient) DeleteWithBody(ctx context.Context, url string, body interface{}, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, url, body, result)
}

// BuildQueryParams builds query parameters from a params struct.
func BuildQueryParams(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := u.Query()
	for key, value := range params {
		if value != "" {
			q.Set(key, value)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// ListParamsToMap converts list params to a map.
func ListParamsToMap(page, perPage int, extra map[string]string) map[string]string {
	params := make(map[string]string)
	if page > 0 {
		params["page"] = strconv.Itoa(page)
	}
	if perPage > 0 {
		params["perPage"] = strconv.Itoa(perPage)
	}
	for k, v := range extra {
		params[k] = v
	}
	return params
}
