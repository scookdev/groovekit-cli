package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuthCommand tests the basic structure of the account command
func TestAccountCommand(t *testing.T) {
	assert.Equal(t, "account", accountCmd.Use)
	assert.Equal(t, "View account information", accountCmd.Short)
	assert.NotEmpty(t, accountCmd.Long)
}

// TestAccountCommandHasSubcommands verifies that show subcommand is registered
func TestAccountCommandHasSubcommands(t *testing.T) {
	commands := accountCmd.Commands()

	// Should have at least 1 subcommand (show)
	assert.GreaterOrEqual(t, len(commands), 1)

	// Find show command
	var hasShow bool
	for _, cmd := range commands {
		if cmd.Use == "show" {
			hasShow = true
		}
	}

	assert.True(t, hasShow, "account command should have show subcommand")
}
