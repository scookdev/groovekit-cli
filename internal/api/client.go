package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/scookdev/groovekit-cli/internal/config"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		BaseURL:    cfg.APIBaseURL,
		HTTPClient: &http.Client{},
		Token:      cfg.AccessToken,
	}
}

// Login authenticates and returns an access token
func (c *Client) Login(email, password string) (string, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/tokens", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", string(bodyBytes))
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

// doRequest is a helper method for authenticated requests
func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)

		// Check if response looks like HTML (common for Rails error pages)
		if len(bodyStr) > 0 && (bodyStr[0] == '<' || strings.Contains(bodyStr, "<!DOCTYPE")) {
			// Don't dump HTML, provide a clean error message
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		// Try to parse as JSON error
		var errResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
			if errResp.Error != "" {
				return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Error)
			}
			if errResp.Message != "" {
				return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Message)
			}
		}

		// Fallback to raw body if it's short
		if len(bodyStr) < 200 {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, bodyStr)
		}

		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}

// GET request
func (c *Client) Get(path string, result interface{}) error {
	return c.doRequest("GET", path, nil, result)
}

// POST request
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	return c.doRequest("POST", path, body, result)
}

// PUT request
func (c *Client) Put(path string, body interface{}, result interface{}) error {
	return c.doRequest("PUT", path, body, result)
}

// DELETE request
func (c *Client) Delete(path string) error {
	return c.doRequest("DELETE", path, nil, nil)
}

// Jobs API methods

// ListJobs returns all jobs for the authenticated user
func (c *Client) ListJobs() (*JobsResponse, error) {
	var result JobsResponse
	if err := c.Get("/jobs", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetJob returns a single job by ID
func (c *Client) GetJob(id string) (*Job, error) {
	var result Job
	if err := c.Get("/jobs/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateJob creates a new job
func (c *Client) CreateJob(req *CreateJobRequest) (*Job, error) {
	payload := map[string]interface{}{
		"job": req,
	}
	var result JobResponse
	if err := c.Post("/jobs", payload, &result); err != nil {
		return nil, err
	}
	return &result.Job, nil
}

// DeleteJob deletes a job by ID
func (c *Client) DeleteJob(id string) error {
	return c.Delete("/jobs/" + id)
}

// Monitors API methods

// ListMonitors returns all monitors for the authenticated user
func (c *Client) ListMonitors() (*MonitorsResponse, error) {
	var result MonitorsResponse
	if err := c.Get("/api_monitors", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMonitor returns a single monitor by ID
func (c *Client) GetMonitor(id string) (*Monitor, error) {
	var result MonitorResponse
	if err := c.Get("/api_monitors/"+id, &result); err != nil {
		return nil, err
	}
	return &result.APIMonitor, nil
}

// CreateMonitor creates a new monitor
func (c *Client) CreateMonitor(req *CreateMonitorRequest) (*Monitor, error) {
	payload := map[string]interface{}{
		"api_monitor": req,
	}
	var result MonitorResponse
	if err := c.Post("/api_monitors", payload, &result); err != nil {
		return nil, err
	}
	return &result.APIMonitor, nil
}

// DeleteMonitor deletes a monitor by ID
func (c *Client) DeleteMonitor(id string) error {
	return c.Delete("/api_monitors/" + id)
}
