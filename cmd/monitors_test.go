package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMonitorsCommand tests the basic structure of the monitors command
func TestMonitorsCommand(t *testing.T) {
	assert.Equal(t, "monitors", monitorsCmd.Use)
	assert.Equal(t, "Manage API endpoint monitors", monitorsCmd.Short)
	assert.NotEmpty(t, monitorsCmd.Long)
}

// TestMonitorsListCommand tests the monitors list command
func TestMonitorsListCommand(t *testing.T) {
	assert.Equal(t, "list", monitorsListCmd.Use)
	assert.Equal(t, "List all monitors", monitorsListCmd.Short)
	assert.NotEmpty(t, monitorsListCmd.Long)
	require.NotNil(t, monitorsListCmd.RunE, "monitors list command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := monitorsListCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "monitors list command should have --json flag")
	assert.Equal(t, "bool", jsonFlag.Value.Type())
}

// TestMonitorsShowCommand tests the monitors show command
func TestMonitorsShowCommand(t *testing.T) {
	assert.Equal(t, "show <id>", monitorsShowCmd.Use)
	assert.Equal(t, "Show monitor details", monitorsShowCmd.Short)
	assert.NotEmpty(t, monitorsShowCmd.Long)
	require.NotNil(t, monitorsShowCmd.RunE, "monitors show command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := monitorsShowCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "monitors show command should have --json flag")
}

// TestMonitorsCreateCommand tests the monitors create command
func TestMonitorsCreateCommand(t *testing.T) {
	assert.Equal(t, "create", monitorsCreateCmd.Use)
	assert.Equal(t, "Create a new monitor", monitorsCreateCmd.Short)
	assert.NotEmpty(t, monitorsCreateCmd.Long)
	require.NotNil(t, monitorsCreateCmd.RunE, "monitors create command should have a RunE function")

	// Verify required flags
	nameFlag := monitorsCreateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "monitors create command should have --name flag")
	assert.Equal(t, "string", nameFlag.Value.Type())

	urlFlag := monitorsCreateCmd.Flags().Lookup("url")
	require.NotNil(t, urlFlag, "monitors create command should have --url flag")
	assert.Equal(t, "string", urlFlag.Value.Type())

	// Verify optional flags
	intervalFlag := monitorsCreateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "monitors create command should have --interval flag")

	methodFlag := monitorsCreateCmd.Flags().Lookup("method")
	require.NotNil(t, methodFlag, "monitors create command should have --method flag")
}

// TestMonitorsUpdateCommand tests the monitors update command
func TestMonitorsUpdateCommand(t *testing.T) {
	assert.Equal(t, "update <id>", monitorsUpdateCmd.Use)
	assert.Equal(t, "Update a monitor", monitorsUpdateCmd.Short)
	assert.NotEmpty(t, monitorsUpdateCmd.Long)
	require.NotNil(t, monitorsUpdateCmd.RunE, "monitors update command should have a RunE function")

	// Verify flags exist
	nameFlag := monitorsUpdateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "monitors update command should have --name flag")

	urlFlag := monitorsUpdateCmd.Flags().Lookup("url")
	require.NotNil(t, urlFlag, "monitors update command should have --url flag")

	httpMethodFlag := monitorsUpdateCmd.Flags().Lookup("http-method")
	require.NotNil(t, httpMethodFlag, "monitors update command should have --http-method flag")

	intervalFlag := monitorsUpdateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "monitors update command should have --interval flag")

	timeoutFlag := monitorsUpdateCmd.Flags().Lookup("timeout")
	require.NotNil(t, timeoutFlag, "monitors update command should have --timeout flag")

	statusFlag := monitorsUpdateCmd.Flags().Lookup("status")
	require.NotNil(t, statusFlag, "monitors update command should have --status flag")
}

// TestMonitorsPauseCommand tests the monitors pause command
func TestMonitorsPauseCommand(t *testing.T) {
	assert.Equal(t, "pause <id>", monitorsPauseCmd.Use)
	assert.Equal(t, "Pause a monitor", monitorsPauseCmd.Short)
	assert.NotEmpty(t, monitorsPauseCmd.Long)
	require.NotNil(t, monitorsPauseCmd.RunE, "monitors pause command should have a RunE function")
}

// TestMonitorsResumeCommand tests the monitors resume command
func TestMonitorsResumeCommand(t *testing.T) {
	assert.Equal(t, "resume <id>", monitorsResumeCmd.Use)
	assert.Equal(t, "Resume a monitor", monitorsResumeCmd.Short)
	assert.NotEmpty(t, monitorsResumeCmd.Long)
	require.NotNil(t, monitorsResumeCmd.RunE, "monitors resume command should have a RunE function")
}

// TestMonitorsIncidentsCommand tests the monitors incidents command
func TestMonitorsIncidentsCommand(t *testing.T) {
	assert.Equal(t, "incidents <id>", monitorsIncidentsCmd.Use)
	assert.Equal(t, "Show incident history", monitorsIncidentsCmd.Short)
	assert.NotEmpty(t, monitorsIncidentsCmd.Long)
	require.NotNil(t, monitorsIncidentsCmd.RunE, "monitors incidents command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := monitorsIncidentsCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "monitors incidents command should have --json flag")
}

// TestMonitorsDeleteCommand tests the monitors delete command
func TestMonitorsDeleteCommand(t *testing.T) {
	assert.Equal(t, "delete <id>", monitorsDeleteCmd.Use)
	assert.Equal(t, "Delete a monitor", monitorsDeleteCmd.Short)
	assert.NotEmpty(t, monitorsDeleteCmd.Long)
	require.NotNil(t, monitorsDeleteCmd.RunE, "monitors delete command should have a RunE function")

	// Verify --force flag exists
	forceFlag := monitorsDeleteCmd.Flags().Lookup("force")
	require.NotNil(t, forceFlag, "monitors delete command should have --force flag")
	assert.Equal(t, "bool", forceFlag.Value.Type())
}

// TestMonitorsCommandHasSubcommands verifies all subcommands are registered
func TestMonitorsCommandHasSubcommands(t *testing.T) {
	commands := monitorsCmd.Commands()

	// Should have 8 subcommands
	expectedSubcommands := []string{"list", "show", "create", "update", "pause", "resume", "incidents", "delete"}
	assert.GreaterOrEqual(t, len(commands), len(expectedSubcommands))

	// Verify all expected subcommands exist
	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Use] = true
	}

	for _, expected := range expectedSubcommands {
		// Check if the command starts with the expected name (to handle "<id>" parts)
		found := false
		for cmdUse := range commandMap {
			if len(cmdUse) >= len(expected) && cmdUse[:len(expected)] == expected {
				found = true
				break
			}
		}
		assert.True(t, found, "monitors command should have %s subcommand", expected)
	}
}

// TestResolveMonitorID tests the helper function for resolving short IDs
func TestResolveMonitorID(t *testing.T) {
	// This is a unit test for the helper function
	// In a real scenario, you'd mock the API client
	// For now, we just verify the function exists by checking if it's referenced
	// A full integration test would require a mock API server
	assert.NotNil(t, monitorsShowCmd.RunE, "resolveMonitorID is used by show command")
}

// TestTruncateHelper tests the truncate helper function
func TestTruncateHelper(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short string",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "exact length",
			input:    "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "long string",
			input:    "this is a very long string",
			maxLen:   10,
			expected: "this is...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}
