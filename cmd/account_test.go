package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthCommand tests the basic structure of the account command
func TestAccountCommand(t *testing.T) {
	assert.Equal(t, "account", accountCmd.Use)
	assert.Equal(t, "View account information", accountCmd.Short)
	assert.NotEmpty(t, accountCmd.Long)
}

// TestAccountCommandHasSubcommands verifies that show are registered
func TestAuthCommandHasSubcommands(t *testing.T) {
	commands := accountCmd.Commands()

	// Should have at least 2 subcommands (login and logout)
	assert.GreaterOrEqual(t, len(commands), 2)

	// Find login and logout commands
	var hasShow bool
	for _, cmd := range commands {
		if cmd.Use == "show" {
			hasShow = true
		}
	}

	assert.True(t, hasShow, "account command should have show subcommand")
}
