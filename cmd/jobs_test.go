package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoginCommand tests the basic structure of the login command
func TestJobsCommand(t *testing.T) {
	assert.Equal(t, "login", loginCmd.Use)
	assert.Equal(t, "Login to GrooveKit", loginCmd.Short)
	assert.NotEmpty(t, loginCmd.Long)

	// Verify RunE function is set
	require.NotNil(t, loginCmd.RunE, "login command should have a RunE function")
}
