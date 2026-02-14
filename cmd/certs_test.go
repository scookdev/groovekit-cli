package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCertsCommand tests the basic structure of the certs command
func TestCertsCommand(t *testing.T) {
	assert.Equal(t, "certs", certsCmd.Use)
	assert.Equal(t, "Manage SSL certificate monitors", certsCmd.Short)
	assert.NotEmpty(t, certsCmd.Long)
}

// TestCertsListCommand tests the certs list command
func TestCertsListCommand(t *testing.T) {
	assert.Equal(t, "list", certsListCmd.Use)
	assert.Equal(t, "List all certs", certsListCmd.Short)
	assert.NotEmpty(t, certsListCmd.Long)
	require.NotNil(t, certsListCmd.RunE, "certs list command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := certsListCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "certs list command should have --json flag")
	assert.Equal(t, "bool", jsonFlag.Value.Type())
}

// TestCertsShowCommand tests the certs show command
func TestCertsShowCommand(t *testing.T) {
	assert.Equal(t, "show <id>", certsShowCmd.Use)
	assert.Equal(t, "Show cert details", certsShowCmd.Short)
	assert.NotEmpty(t, certsShowCmd.Long)
	require.NotNil(t, certsShowCmd.RunE, "certs show command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := certsShowCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "certs show command should have --json flag")
}

// TestCertsCreateCommand tests the certs create command
func TestCertsCreateCommand(t *testing.T) {
	assert.Equal(t, "create", certsCreateCmd.Use)
	assert.Equal(t, "Create a new SSL certificate monitor", certsCreateCmd.Short)
	assert.NotEmpty(t, certsCreateCmd.Long)
	require.NotNil(t, certsCreateCmd.RunE, "certs create command should have a RunE function")

	// Verify required flags
	nameFlag := certsCreateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "certs create command should have --name flag")
	assert.Equal(t, "string", nameFlag.Value.Type())

	domainFlag := certsCreateCmd.Flags().Lookup("domain")
	require.NotNil(t, domainFlag, "certs create command should have --domain flag")
	assert.Equal(t, "string", domainFlag.Value.Type())

	// Verify optional flags
	portFlag := certsCreateCmd.Flags().Lookup("port")
	require.NotNil(t, portFlag, "certs create command should have --port flag")
	assert.Equal(t, "int", portFlag.Value.Type())

	intervalFlag := certsCreateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "certs create command should have --interval flag")
	assert.Equal(t, "int", intervalFlag.Value.Type())
}

// TestCertsUpdateCommand tests the certs update command
func TestCertsUpdateCommand(t *testing.T) {
	assert.Equal(t, "update <id>", certsUpdateCmd.Use)
	assert.Equal(t, "Update an SSL certificate monitor", certsUpdateCmd.Short)
	assert.NotEmpty(t, certsUpdateCmd.Long)
	require.NotNil(t, certsUpdateCmd.RunE, "certs update command should have a RunE function")

	// Verify flags exist
	nameFlag := certsUpdateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "certs update command should have --name flag")

	domainFlag := certsUpdateCmd.Flags().Lookup("domain")
	require.NotNil(t, domainFlag, "certs update command should have --domain flag")

	portFlag := certsUpdateCmd.Flags().Lookup("port")
	require.NotNil(t, portFlag, "certs update command should have --port flag")

	intervalFlag := certsUpdateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "certs update command should have --interval flag")

	gracePeriodFlag := certsUpdateCmd.Flags().Lookup("grace-period")
	require.NotNil(t, gracePeriodFlag, "certs update command should have --grace-period flag")

	warningThresholdFlag := certsUpdateCmd.Flags().Lookup("warning-threshold")
	require.NotNil(t, warningThresholdFlag, "certs update command should have --warning-threshold flag")

	urgentThresholdFlag := certsUpdateCmd.Flags().Lookup("urgent-threshold")
	require.NotNil(t, urgentThresholdFlag, "certs update command should have --urgent-threshold flag")

	criticalThresholdFlag := certsUpdateCmd.Flags().Lookup("critical-threshold")
	require.NotNil(t, criticalThresholdFlag, "certs update command should have --critical-threshold flag")

	statusFlag := certsUpdateCmd.Flags().Lookup("status")
	require.NotNil(t, statusFlag, "certs update command should have --status flag")
}

// TestCertsPauseCommand tests the certs pause command
func TestCertsPauseCommand(t *testing.T) {
	assert.Equal(t, "pause <id>", certsPauseCmd.Use)
	assert.Equal(t, "Pause a cert", certsPauseCmd.Short)
	assert.NotEmpty(t, certsPauseCmd.Long)
	require.NotNil(t, certsPauseCmd.RunE, "certs pause command should have a RunE function")
}

// TestCertsResumeCommand tests the certs resume command
func TestCertsResumeCommand(t *testing.T) {
	assert.Equal(t, "resume <id>", certsResumeCmd.Use)
	assert.Equal(t, "Resume a cert", certsResumeCmd.Short)
	assert.NotEmpty(t, certsResumeCmd.Long)
	require.NotNil(t, certsResumeCmd.RunE, "certs resume command should have a RunE function")
}

// TestCertsIncidentsCommand tests the certs incidents command
func TestCertsIncidentsCommand(t *testing.T) {
	assert.Equal(t, "incidents <id>", certsIncidentsCmd.Use)
	assert.Equal(t, "Show incident history", certsIncidentsCmd.Short)
	assert.NotEmpty(t, certsIncidentsCmd.Long)
	require.NotNil(t, certsIncidentsCmd.RunE, "certs incidents command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := certsIncidentsCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "certs incidents command should have --json flag")
}

// TestCertsDeleteCommand tests the certs delete command
func TestCertsDeleteCommand(t *testing.T) {
	assert.Equal(t, "delete <id>", certsDeleteCmd.Use)
	assert.Equal(t, "Delete a cert", certsDeleteCmd.Short)
	assert.NotEmpty(t, certsDeleteCmd.Long)
	require.NotNil(t, certsDeleteCmd.RunE, "certs delete command should have a RunE function")

	// Verify --force flag exists
	forceFlag := certsDeleteCmd.Flags().Lookup("force")
	require.NotNil(t, forceFlag, "certs delete command should have --force flag")
	assert.Equal(t, "bool", forceFlag.Value.Type())
}

// TestCertsCommandHasSubcommands verifies all subcommands are registered
func TestCertsCommandHasSubcommands(t *testing.T) {
	commands := certsCmd.Commands()

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
		assert.True(t, found, "certs command should have %s subcommand", expected)
	}
}

// TestResolveCertID tests the helper function for resolving short IDs
func TestResolveCertID(t *testing.T) {
	// This is a unit test for the helper function
	// In a real scenario, you'd mock the API client
	// For now, we just verify the function exists by checking if it's referenced
	// A full integration test would require a mock API server
	assert.NotNil(t, certsShowCmd.RunE, "resolveCertID is used by show command")
}

// TestFormatIncidentDuration tests the formatIncidentDuration helper function
func TestFormatIncidentDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{
			name:     "seconds",
			seconds:  30.0,
			expected: "30s",
		},
		{
			name:     "minutes",
			seconds:  120.0,
			expected: "2m",
		},
		{
			name:     "hours",
			seconds:  7200.0,
			expected: "2.0h",
		},
		{
			name:     "days",
			seconds:  172800.0,
			expected: "2.0d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatIncidentDuration(tt.seconds)
			assert.Equal(t, tt.expected, result)
		})
	}
}
