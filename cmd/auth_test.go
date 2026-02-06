package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthCommand tests the basic structure of the auth command
func TestAuthCommand(t *testing.T) {
	assert.Equal(t, "auth", authCmd.Use)
	assert.Equal(t, "Manage authentication", authCmd.Short)
	assert.NotEmpty(t, authCmd.Long)
}

// TestLoginCommand tests the basic structure of the login command
func TestLoginCommand(t *testing.T) {
	assert.Equal(t, "login", loginCmd.Use)
	assert.Equal(t, "Login to GrooveKit", loginCmd.Short)
	assert.NotEmpty(t, loginCmd.Long)

	// Verify RunE function is set
	require.NotNil(t, loginCmd.RunE, "login command should have a RunE function")
}

// TestLogoutCommand tests the basic structure of the logout command
func TestLogoutCommand(t *testing.T) {
	assert.Equal(t, "logout", logoutCmd.Use)
	assert.Equal(t, "Logout from GrooveKit", logoutCmd.Short)
	assert.NotEmpty(t, logoutCmd.Long)

	// Verify RunE function is set
	require.NotNil(t, logoutCmd.RunE, "logout command should have a RunE function")
}

// TestAuthCommandHasSubcommands verifies that login and logout are registered
func TestAuthCommandHasSubcommands(t *testing.T) {
	commands := authCmd.Commands()

	// Should have at least 2 subcommands (login and logout)
	assert.GreaterOrEqual(t, len(commands), 2)

	// Find login and logout commands
	var hasLogin, hasLogout bool
	for _, cmd := range commands {
		if cmd.Use == "login" {
			hasLogin = true
		}
		if cmd.Use == "logout" {
			hasLogout = true
		}
	}

	assert.True(t, hasLogin, "auth command should have login subcommand")
	assert.True(t, hasLogout, "auth command should have logout subcommand")
}
