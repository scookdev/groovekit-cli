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
// 
// Monitor represents a cron job monitor
type Monitor struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Status         string   `json:"status"`
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
type MonitorsResponse struct {
	Monitorss       []Monitor `json:"jobs"`
	HasMore    bool  `json:"has_more"`
	TotalCount int   `json:"total_count"`
}

// JobResponse represents the response from POST/PUT /jobs
type MonitorResponse struct {
	Monitor Monitor `json:"job"`
}

// CreateJobRequest represents the request body for creating a job
type CreateMonitorRequest struct {
	Name          string   `json:"name"`
	Status        string   `json:"status,omitempty"`
}
