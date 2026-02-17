// Package api provides the HTTP client for interacting with the GrooveKit API
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

// Client represents an HTTP client for the GrooveKit API
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
	defer func() { _ = resp.Body.Close() }()

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
	defer func() { _ = resp.Body.Close() }()

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

// Get performs a GET request to the API
func (c *Client) Get(path string, result interface{}) error {
	return c.doRequest("GET", path, nil, result)
}

// Post performs a POST request to the API
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	return c.doRequest("POST", path, body, result)
}

// Put performs a PUT request to the API
func (c *Client) Put(path string, body interface{}, result interface{}) error {
	return c.doRequest("PUT", path, body, result)
}

// Delete performs a DELETE request to the API
func (c *Client) Delete(path string) error {
	return c.doRequest("DELETE", path, nil, nil)
}

// Account API method

// GetAccount returns account information with subscription and usage
func (c *Client) GetAccount() (*Account, error) {
	var account Account
	if err := c.Get("/users/me", &account); err != nil {
		return nil, err
	}
	return &account, nil
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
	payload := map[string]any{
		"job": req,
	}
	var result JobResponse
	if err := c.Post("/jobs", payload, &result); err != nil {
		return nil, err
	}
	return &result.Job, nil
}

// UpdateJob updates an existing job
func (c *Client) UpdateJob(id string, req *UpdateJobRequest) (*Job, error) {
	payload := map[string]any{
		"job": req,
	}
	var result JobResponse
	if err := c.Put("/jobs/"+id, payload, &result); err != nil {
		return nil, err
	}
	return &result.Job, nil
}

// DeleteJob deletes a job by ID
func (c *Client) DeleteJob(id string) error {
	return c.Delete("/jobs/" + id)
}

// ListJobPings returns recent pings for a job
func (c *Client) ListJobPings(id string) ([]Ping, error) {
	var result struct {
		Pings []Ping `json:"pings"`
	}
	if err := c.Get("/jobs/"+id+"/pings", &result); err != nil {
		return nil, err
	}
	return result.Pings, nil
}

// ListJobIncidents returns incident history for a job
func (c *Client) ListJobIncidents(id string) ([]Incident, error) {
	var result struct {
		Incidents []Incident `json:"incidents"`
	}
	if err := c.Get("/jobs/"+id+"/incidents", &result); err != nil {
		return nil, err
	}
	return result.Incidents, nil
}

// Monitors API methods

// ListMonitors returns all api monitors for the authenticated user
func (c *Client) ListMonitors() (*MonitorsResponse, error) {
	var result MonitorsResponse
	if err := c.Get("/api_monitors", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMonitor returns a single api monitor by ID
func (c *Client) GetMonitor(id string) (*Monitor, error) {
	var result MonitorResponse
	if err := c.Get("/api_monitors/"+id, &result); err != nil {
		return nil, err
	}
	return &result.APIMonitor, nil
}

// CreateMonitor creates a new api monitor
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

// UpdateMonitor updates an existing api monitor
func (c *Client) UpdateMonitor(id string, req *UpdateMonitorRequest) (*Monitor, error) {
	payload := map[string]interface{}{
		"api_monitor": req,
	}
	var result MonitorResponse
	if err := c.Put("/api_monitors/"+id, payload, &result); err != nil {
		return nil, err
	}
	return &result.APIMonitor, nil
}

// DeleteMonitor deletes an api monitor by ID
func (c *Client) DeleteMonitor(id string) error {
	return c.Delete("/api_monitors/" + id)
}

// ListMonitorChecks returns recent checks for an api monitor
func (c *Client) ListMonitorChecks(id string) ([]Check, error) {
	var result struct {
		APIChecks []Check `json:"api_checks"`
	}
	if err := c.Get("/api_monitors/"+id+"/api_checks", &result); err != nil {
		return nil, err
	}
	return result.APIChecks, nil
}

// ListMonitorIncidents returns incident history for an api monitor
func (c *Client) ListMonitorIncidents(id string) ([]Incident, error) {
	var result struct {
		Incidents []Incident `json:"incidents"`
	}
	if err := c.Get("/api_monitors/"+id+"/incidents", &result); err != nil {
		return nil, err
	}
	return result.Incidents, nil
}

// SSL Certificate Monitor API method

// GetCert returns a single api certificate by ID
func (c *Client) GetCert(id string) (*SslMonitor, error) {
	var result SslMonitorResponse
	if err := c.Get("/ssl_monitors/"+id, &result); err != nil {
		return nil, err
	}
	return &result.SslMonitor, nil
}

// ListCerts returns all ssl monitors for the authenticated user
func (c *Client) ListCerts() (*SslMonitorsResponse, error) {
	var result SslMonitorsResponse
	if err := c.Get("/ssl_monitors", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateCert creates a new SSL monitor
func (c *Client) CreateCert(req *CreateSslMonitorRequest) (*SslMonitor, error) {
	payload := map[string]interface{}{
		"ssl_monitor": req,
	}
	var result SslMonitorResponse
	if err := c.Post("/ssl_monitors", payload, &result); err != nil {
		return nil, err
	}
	return &result.SslMonitor, nil
}

// UpdateCert updates an existing SSL monitor
func (c *Client) UpdateCert(id string, req *UpdateSslMonitorRequest) (*SslMonitor, error) {
	payload := map[string]interface{}{
		"ssl_monitor": req,
	}
	var result SslMonitorResponse
	if err := c.Put("/ssl_monitors/"+id, payload, &result); err != nil {
		return nil, err
	}
	return &result.SslMonitor, nil
}

// DeleteCert deletes an SSL monitor by ID
func (c *Client) DeleteCert(id string) error {
	return c.Delete("/ssl_monitors/" + id)
}

// ListCertIncidents returns incident history for an SSL monitor
func (c *Client) ListCertIncidents(id string) ([]Incident, error) {
	var result struct {
		Incidents []Incident `json:"incidents"`
	}
	if err := c.Get("/ssl_monitors/"+id+"/incidents", &result); err != nil {
		return nil, err
	}
	return result.Incidents, nil
}

// Domain Monitor API methods

// ListDomains returns all domain monitors for the authenticated user
func (c *Client) ListDomains() (*DomainMonitorsResponse, error) {
	var result DomainMonitorsResponse
	if err := c.Get("/domain_monitors", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDomain returns a single domain monitor by ID
func (c *Client) GetDomain(id string) (*DomainMonitor, error) {
	var result DomainMonitorResponse
	if err := c.Get("/domain_monitors/"+id, &result); err != nil {
		return nil, err
	}
	return &result.DomainMonitor, nil
}

// CreateDomain creates a new domain monitor
func (c *Client) CreateDomain(req *CreateDomainMonitorRequest) (*DomainMonitor, error) {
	payload := map[string]interface{}{
		"domain_monitor": req,
	}
	var result DomainMonitorResponse
	if err := c.Post("/domain_monitors", payload, &result); err != nil {
		return nil, err
	}
	return &result.DomainMonitor, nil
}

// UpdateDomain updates an existing domain monitor
func (c *Client) UpdateDomain(id string, req *UpdateDomainMonitorRequest) (*DomainMonitor, error) {
	payload := map[string]interface{}{
		"domain_monitor": req,
	}
	var result DomainMonitorResponse
	if err := c.Put("/domain_monitors/"+id, payload, &result); err != nil {
		return nil, err
	}
	return &result.DomainMonitor, nil
}

// DeleteDomain deletes a domain monitor by ID
func (c *Client) DeleteDomain(id string) error {
	return c.Delete("/domain_monitors/" + id)
}

// ListDomainIncidents returns incident history for a domain monitor
func (c *Client) ListDomainIncidents(id string) ([]Incident, error) {
	var result struct {
		Incidents []Incident `json:"incidents"`
	}
	if err := c.Get("/domain_monitors/"+id+"/incidents", &result); err != nil {
		return nil, err
	}
	return result.Incidents, nil
}

// DNS Monitor API methods

// ListDnsMonitors returns all DNS monitors for the authenticated user
func (c *Client) ListDnsMonitors() (*DnsMonitorsResponse, error) {
	var result DnsMonitorsResponse
	if err := c.Get("/dns_monitors", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDnsMonitor returns a single DNS monitor by ID
func (c *Client) GetDnsMonitor(id string) (*DnsMonitor, error) {
	var result DnsMonitorResponse
	if err := c.Get("/dns_monitors/"+id, &result); err != nil {
		return nil, err
	}
	return &result.DnsMonitor, nil
}

// CreateDnsMonitor creates a new DNS monitor
func (c *Client) CreateDnsMonitor(req *CreateDnsMonitorRequest) (*DnsMonitor, error) {
	payload := map[string]interface{}{
		"dns_monitor": req,
	}
	var result DnsMonitorResponse
	if err := c.Post("/dns_monitors", payload, &result); err != nil {
		return nil, err
	}
	return &result.DnsMonitor, nil
}

// UpdateDnsMonitor updates an existing DNS monitor
func (c *Client) UpdateDnsMonitor(id string, req *UpdateDnsMonitorRequest) (*DnsMonitor, error) {
	payload := map[string]interface{}{
		"dns_monitor": req,
	}
	var result DnsMonitorResponse
	if err := c.Put("/dns_monitors/"+id, payload, &result); err != nil {
		return nil, err
	}
	return &result.DnsMonitor, nil
}

// DeleteDnsMonitor deletes a DNS monitor by ID
func (c *Client) DeleteDnsMonitor(id string) error {
	return c.Delete("/dns_monitors/" + id)
}

// ListDnsMonitorIncidents returns incident history for a DNS monitor
func (c *Client) ListDnsMonitorIncidents(id string) ([]Incident, error) {
	var result struct {
		Incidents []Incident `json:"incidents"`
	}
	if err := c.Get("/dns_monitors/"+id+"/incidents", &result); err != nil {
		return nil, err
	}
	return result.Incidents, nil
}
