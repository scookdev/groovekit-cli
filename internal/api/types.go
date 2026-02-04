package api

// Job types
//
// Job represents a cron job monitor
type Job struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Interval       int      `json:"interval"`
	GracePeriod    int      `json:"grace_period"`
	Status         string   `json:"status"`
	PingToken      string   `json:"ping_token"`
	WebhookURL     string   `json:"webhook_url"`
	WebhookSecret  string   `json:"webhook_secret"`
	AllowedIPs     []string `json:"allowed_ips"`
	LastPingAt     *string  `json:"last_ping_at"`
	LastRunAt      *string  `json:"last_run_at"`
	LastAlertedAt  *string  `json:"last_alerted_at"`
	Down           bool     `json:"down"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// JobsResponse represents the response from GET /jobs
type JobsResponse struct {
	Jobs       []Job `json:"jobs"`
	HasMore    bool  `json:"has_more"`
	TotalCount int   `json:"total_count"`
}

// JobResponse represents the response from POST/PUT /jobs
type JobResponse struct {
	Job Job `json:"job"`
}

// CreateJobRequest represents the request body for creating a job
type CreateJobRequest struct {
	Name          string   `json:"name"`
	Interval      int      `json:"interval"`
	GracePeriod   int      `json:"grace_period,omitempty"`
	Status        string   `json:"status,omitempty"`
	WebhookURL    string   `json:"webhook_url,omitempty"`
	WebhookSecret string   `json:"webhook_secret,omitempty"`
	AllowedIPs    []string `json:"allowed_ips,omitempty"`
}

// Monitor types

// Monitor represents an API endpoint monitor
type Monitor struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	URL                   string        `json:"url"`
	HTTPMethod            string        `json:"http_method"`
	Headers               interface{}   `json:"headers"`
	ExpectedStatusCodes   []int         `json:"expected_status_codes"`
	Timeout               int           `json:"timeout"`
	Interval              int           `json:"interval"`
	GracePeriod           int           `json:"grace_period"`
	Status                string        `json:"status"`
	APICheckToken         string        `json:"api_check_token"`
	HasAuthHeaders        bool          `json:"has_auth_headers"`
	ValidateResponsePaths []string      `json:"validate_response_paths"`
	JSONSchema            *string       `json:"json_schema"`
	RequestBody           *string       `json:"request_body"`
	LastCheckAt           *string       `json:"last_check_at"`
	ConsecutiveFailures   int           `json:"consecutive_failures"`
	Down                  bool          `json:"down"`
	UptimePercentage      *float64      `json:"uptime_percentage"`
	AverageResponseTime   *float64      `json:"average_response_time"`
	CreatedAt             string        `json:"created_at"`
	UpdatedAt             string        `json:"updated_at"`
}

// MonitorsResponse represents the response from GET /api_monitors
type MonitorsResponse struct {
	APIMonitors []Monitor `json:"api_monitors"`
}

// MonitorResponse represents the response from POST/PUT /api_monitors
type MonitorResponse struct {
	APIMonitor Monitor `json:"api_monitor"`
}

// CreateMonitorRequest represents the request body for creating a monitor
type CreateMonitorRequest struct {
	Name                  string   `json:"name"`
	URL                   string   `json:"url"`
	HTTPMethod            string   `json:"http_method,omitempty"`
	Interval              int      `json:"interval,omitempty"`
	ExpectedStatusCodes   []int    `json:"expected_status_codes,omitempty"`
	Timeout               int      `json:"timeout,omitempty"`
	GracePeriod           int      `json:"grace_period,omitempty"`
	Status                string   `json:"status,omitempty"`
}

// ApiCheck represents an API health check result
type ApiCheck struct {
	ID              string  `json:"id"`
	APIMonitorID    string  `json:"api_monitor_id"`
	StatusCode      int     `json:"status_code"`
	ResponseTime    int     `json:"response_time"`
	Success         bool    `json:"success"`
	ErrorMessage    *string `json:"error_message"`
	ValidationError *string `json:"validation_error"`
	CreatedAt       string  `json:"created_at"`
}

// Ping represents a job heartbeat ping
type Ping struct {
	ID        string  `json:"id"`
	JobID     string  `json:"job_id"`
	PingType  string  `json:"ping_type"`
	Duration  *int    `json:"duration"`
	CreatedAt string  `json:"created_at"`
}
