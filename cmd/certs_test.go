package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuthCommand tests the basic structure of the auth command
func TestCertsCommand(t *testing.T) {
	assert.Equal(t, "auth", authCmd.Use)
	assert.Equal(t, "Manage authentication", authCmd.Short)
	assert.NotEmpty(t, authCmd.Long)
}
