package config

import (
	"os"
	"testing"
)

func TestLoad_WithTokenEnvVar(t *testing.T) {
	// Set up test environment variable
	testToken := "test-token-12345"
	os.Setenv("GROOVEKIT_TOKEN", testToken)
	defer os.Unsetenv("GROOVEKIT_TOKEN")

	// Load config (may not have config file, that's ok)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify token from env var is used
	if cfg.AccessToken != testToken {
		t.Errorf("Expected AccessToken to be %q, got %q", testToken, cfg.AccessToken)
	}
}

func TestLoad_WithAPIURLEnvVar(t *testing.T) {
	// Set up test environment variable
	testURL := "http://localhost:3000"
	os.Setenv("GROOVEKIT_API_URL", testURL)
	defer os.Unsetenv("GROOVEKIT_API_URL")

	// Load config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify URL from env var is used
	if cfg.APIBaseURL != testURL {
		t.Errorf("Expected APIBaseURL to be %q, got %q", testURL, cfg.APIBaseURL)
	}
}

func TestLoad_EnvVarsTakePrecedence(t *testing.T) {
	// Even if a config file exists with different values,
	// env vars should take precedence
	testToken := "env-token-wins"
	testURL := "http://env-url-wins:8080"

	os.Setenv("GROOVEKIT_TOKEN", testToken)
	os.Setenv("GROOVEKIT_API_URL", testURL)
	defer os.Unsetenv("GROOVEKIT_TOKEN")
	defer os.Unsetenv("GROOVEKIT_API_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.AccessToken != testToken {
		t.Errorf("Expected env var token to take precedence, got %q", cfg.AccessToken)
	}

	if cfg.APIBaseURL != testURL {
		t.Errorf("Expected env var URL to take precedence, got %q", cfg.APIBaseURL)
	}
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name        string
		accessToken string
		want        bool
	}{
		{"with token", "some-token", true},
		{"without token", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{AccessToken: tt.accessToken}
			if got := cfg.IsAuthenticated(); got != tt.want {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}
