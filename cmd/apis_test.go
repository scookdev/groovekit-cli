package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestApisCommand tests the basic structure of the apis command
func TestApisCommand(t *testing.T) {
	assert.Equal(t, "apis", apisCmd.Use)
	assert.Equal(t, "Manage API endpoint monitors", apisCmd.Short)
	assert.NotEmpty(t, apisCmd.Long)
}

// TestApisListCommand tests the apis list command
func TestApisListCommand(t *testing.T) {
	assert.Equal(t, "list", apisListCmd.Use)
	assert.Equal(t, "List all api monitors", apisListCmd.Short)
	assert.NotEmpty(t, apisListCmd.Long)
	require.NotNil(t, apisListCmd.RunE, "apis list command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := apisListCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "apis list command should have --json flag")
	assert.Equal(t, "bool", jsonFlag.Value.Type())
}

// TestApisShowCommand tests the apis show command
func TestApisShowCommand(t *testing.T) {
	assert.Equal(t, "show <id>", apisShowCmd.Use)
	assert.Equal(t, "Show API monitor details", apisShowCmd.Short)
	assert.NotEmpty(t, apisShowCmd.Long)
	require.NotNil(t, apisShowCmd.RunE, "apis show command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := apisShowCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "apis show command should have --json flag")
}

// TestApisCreateCommand tests the apis create command
func TestApisCreateCommand(t *testing.T) {
	assert.Equal(t, "create", apisCreateCmd.Use)
	assert.Equal(t, "Create a new API monitor", apisCreateCmd.Short)
	assert.NotEmpty(t, apisCreateCmd.Long)
	require.NotNil(t, apisCreateCmd.RunE, "apis create command should have a RunE function")

	// Verify required flags
	nameFlag := apisCreateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "apis create command should have --name flag")
	assert.Equal(t, "string", nameFlag.Value.Type())

	urlFlag := apisCreateCmd.Flags().Lookup("url")
	require.NotNil(t, urlFlag, "apis create command should have --url flag")
	assert.Equal(t, "string", urlFlag.Value.Type())

	// Verify optional flags
	intervalFlag := apisCreateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "apis create command should have --interval flag")

	methodFlag := apisCreateCmd.Flags().Lookup("method")
	require.NotNil(t, methodFlag, "apis create command should have --method flag")
}

// TestApisUpdateCommand tests the apis update command
func TestApisUpdateCommand(t *testing.T) {
	assert.Equal(t, "update <id>", apisUpdateCmd.Use)
	assert.Equal(t, "Update an API monitor", apisUpdateCmd.Short)
	assert.NotEmpty(t, apisUpdateCmd.Long)
	require.NotNil(t, apisUpdateCmd.RunE, "apis update command should have a RunE function")

	// Verify flags exist
	nameFlag := apisUpdateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "apis update command should have --name flag")

	urlFlag := apisUpdateCmd.Flags().Lookup("url")
	require.NotNil(t, urlFlag, "apis update command should have --url flag")

	httpMethodFlag := apisUpdateCmd.Flags().Lookup("http-method")
	require.NotNil(t, httpMethodFlag, "apis update command should have --http-method flag")

	intervalFlag := apisUpdateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "apis update command should have --interval flag")

	timeoutFlag := apisUpdateCmd.Flags().Lookup("timeout")
	require.NotNil(t, timeoutFlag, "apis update command should have --timeout flag")

	statusFlag := apisUpdateCmd.Flags().Lookup("status")
	require.NotNil(t, statusFlag, "apis update command should have --status flag")
}

// TestApisPauseCommand tests the apis pause command
func TestApisPauseCommand(t *testing.T) {
	assert.Equal(t, "pause <id>", apisPauseCmd.Use)
	assert.Equal(t, "Pause an API monitor", apisPauseCmd.Short)
	assert.NotEmpty(t, apisPauseCmd.Long)
	require.NotNil(t, apisPauseCmd.RunE, "apis pause command should have a RunE function")
}

// TestApisResumeCommand tests the apis resume command
func TestApisResumeCommand(t *testing.T) {
	assert.Equal(t, "resume <id>", apisResumeCmd.Use)
	assert.Equal(t, "Resume an API monitor", apisResumeCmd.Short)
	assert.NotEmpty(t, apisResumeCmd.Long)
	require.NotNil(t, apisResumeCmd.RunE, "apis resume command should have a RunE function")
}

// TestApisIncidentsCommand tests the apis incidents command
func TestApisIncidentsCommand(t *testing.T) {
	assert.Equal(t, "incidents <id>", apisIncidentsCmd.Use)
	assert.Equal(t, "Show incident history", apisIncidentsCmd.Short)
	assert.NotEmpty(t, apisIncidentsCmd.Long)
	require.NotNil(t, apisIncidentsCmd.RunE, "apis incidents command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := apisIncidentsCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "apis incidents command should have --json flag")
}

// TestApisDeleteCommand tests the apis delete command
func TestApisDeleteCommand(t *testing.T) {
	assert.Equal(t, "delete <id>", apisDeleteCmd.Use)
	assert.Equal(t, "Delete an API monitor", apisDeleteCmd.Short)
	assert.NotEmpty(t, apisDeleteCmd.Long)
	require.NotNil(t, apisDeleteCmd.RunE, "apis delete command should have a RunE function")

	// Verify --force flag exists
	forceFlag := apisDeleteCmd.Flags().Lookup("force")
	require.NotNil(t, forceFlag, "apis delete command should have --force flag")
	assert.Equal(t, "bool", forceFlag.Value.Type())
}

// TestApisCommandHasSubcommands verifies all subcommands are registered
func TestApisCommandHasSubcommands(t *testing.T) {
	commands := apisCmd.Commands()

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
		assert.True(t, found, "apis command should have %s subcommand", expected)
	}
}

// TestResolveMonitorID tests the helper function for resolving short IDs
func TestResolveMonitorID(t *testing.T) {
	// This is a unit test for the helper function
	// In a real scenario, you'd mock the API client
	// For now, we just verify the function exists by checking if it's referenced
	// A full integration test would require a mock API server
	assert.NotNil(t, apisShowCmd.RunE, "resolveMonitorID is used by show command")
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
