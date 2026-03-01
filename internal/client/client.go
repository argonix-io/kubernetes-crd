package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// NotFoundError is returned when the API returns 404.
type NotFoundError struct {
	Endpoint string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", e.Endpoint)
}

// IsNotFound checks if the error is a NotFoundError.
func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

// Client is an HTTP client for the Argonix API.
type Client struct {
	BaseURL        string
	APIKey         string
	OrganizationID string
	HTTPClient     *http.Client
}

// NewClient creates a new Argonix API client.
// It auto-discovers the organization ID from the API key.
func NewClient(ctx context.Context, baseURL, apiKey string) (*Client, error) {
	if baseURL == "" {
		baseURL = "https://api.argonix.io"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	c := &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Auto-discover organization ID from API key.
	type apiKeyInfo struct {
		OrganizationID string `json:"organization_id"`
	}
	var info apiKeyInfo
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/0.1/auth/api-key-info/", baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("creating api-key-info request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", apiKey))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching api-key-info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api-key-info returned %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decoding api-key-info: %w", err)
	}

	c.OrganizationID = info.OrganizationID
	return c, nil
}

func (c *Client) url(endpoint string) string {
	return fmt.Sprintf("%s/api/0.1/organizations/%s%s", c.BaseURL, c.OrganizationID, endpoint)
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, payload interface{}, result interface{}) error {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshaling payload: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.url(endpoint), body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", c.APIKey))
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{Endpoint: endpoint}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// Create sends a POST request to create a resource.
func (c *Client) Create(ctx context.Context, endpoint string, payload interface{}, result interface{}) error {
	return c.doRequest(ctx, http.MethodPost, endpoint, payload, result)
}

// Read sends a GET request to read a single resource.
func (c *Client) Read(ctx context.Context, endpoint string, result interface{}) error {
	return c.doRequest(ctx, http.MethodGet, endpoint, nil, result)
}

// Update sends a PUT request to fully update a resource.
func (c *Client) Update(ctx context.Context, endpoint string, payload interface{}, result interface{}) error {
	return c.doRequest(ctx, http.MethodPut, endpoint, payload, result)
}

// Delete sends a DELETE request to remove a resource.
func (c *Client) Delete(ctx context.Context, endpoint string) error {
	return c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil)
}

// List sends a GET request to list resources.
func (c *Client) List(ctx context.Context, endpoint string, result interface{}) error {
	return c.doRequest(ctx, http.MethodGet, endpoint, nil, result)
}
