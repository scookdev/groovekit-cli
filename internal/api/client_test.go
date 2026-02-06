package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scookdev/groovekit-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogin_Success tests successful login
func TestLogin_Success(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/tokens", r.URL.Path)

		// Verify content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Decode and verify request body
		var payload map[string]string
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", payload["email"])
		assert.Equal(t, "password123", payload["password"])

		// Send success response
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"access_token": "mock-token-123",
		})
	}))
	defer server.Close()

	// Create client with test server URL
	cfg := &config.Config{
		APIBaseURL: server.URL,
	}
	client := NewClient(cfg)

	// Test login
	token, err := client.Login("test@example.com", "password123")

	// Assert results
	require.NoError(t, err)
	assert.Equal(t, "mock-token-123", token)
}

// TestLogin_InvalidCredentials tests login with wrong credentials
func TestLogin_InvalidCredentials(t *testing.T) {
	// Create a mock HTTP server that returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "Invalid credentials"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIBaseURL: server.URL,
	}
	client := NewClient(cfg)

	token, err := client.Login("test@example.com", "wrongpassword")

	// Assert error occurred and no token returned
	require.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "login failed")
}

// TestLogin_ServerError tests handling of server errors
func TestLogin_ServerError(t *testing.T) {
	// Create a mock HTTP server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	cfg := &config.Config{
		APIBaseURL: server.URL,
	}
	client := NewClient(cfg)

	token, err := client.Login("test@example.com", "password123")

	// Assert error occurred
	require.Error(t, err)
	assert.Empty(t, token)
}

// TestLogin_NetworkError tests handling of network errors
func TestLogin_NetworkError(t *testing.T) {
	// Use an invalid URL to simulate network error
	cfg := &config.Config{
		APIBaseURL: "http://invalid-url-that-does-not-exist.local",
	}
	client := NewClient(cfg)

	token, err := client.Login("test@example.com", "password123")

	// Assert error occurred
	require.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "login request failed")
}

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		APIBaseURL:  "https://api.example.com",
		AccessToken: "test-token",
	}

	client := NewClient(cfg)

	require.NotNil(t, client)
	assert.Equal(t, "https://api.example.com", client.BaseURL)
	assert.Equal(t, "test-token", client.Token)
	assert.NotNil(t, client.HTTPClient)
}
