package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
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
	if err := c.Get("/monitors", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetJob returns a single job by ID
func (c *Client) GetMonitor(id string) (*Monitor, error) {
	var result Monitor
	if err := c.Get("/monitors/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateJob creates a new job
func (c *Client) CreateMonitor(req *CreateMonitorRequest) (*Monitor, error) {
	payload := map[string]interface{}{
		"job": req,
	}
	var result MonitorResponse
	if err := c.Post("/monitors", payload, &result); err != nil {
		return nil, err
	}
	return &result.Monitor, nil
}

// DeleteJob deletes a job by ID
func (c *Client) DeleteMonitor(id string) error {
	return c.Delete("/monitors/" + id)
}
