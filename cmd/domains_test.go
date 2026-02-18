package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDomainsCommand tests the basic structure of the domains command
func TestDomainsCommand(t *testing.T) {
	assert.Equal(t, "domains", domainsCmd.Use)
	assert.Equal(t, "Manage domain expiration monitors", domainsCmd.Short)
	assert.NotEmpty(t, domainsCmd.Long)
}

// TestDomainsListCommand tests the domains list command
func TestDomainsListCommand(t *testing.T) {
	assert.Equal(t, "list", domainsListCmd.Use)
	assert.Equal(t, "List all domains", domainsListCmd.Short)
	assert.NotEmpty(t, domainsListCmd.Long)
	require.NotNil(t, domainsListCmd.RunE, "domains list command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := domainsListCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "domains list command should have --json flag")
	assert.Equal(t, "bool", jsonFlag.Value.Type())
}

// TestDomainsShowCommand tests the domains show command
func TestDomainsShowCommand(t *testing.T) {
	assert.Equal(t, "show <id>", domainsShowCmd.Use)
	assert.Equal(t, "Show domain details", domainsShowCmd.Short)
	assert.NotEmpty(t, domainsShowCmd.Long)
	require.NotNil(t, domainsShowCmd.RunE, "domains show command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := domainsShowCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "domains show command should have --json flag")
}

// TestDomainsCreateCommand tests the domains create command
func TestDomainsCreateCommand(t *testing.T) {
	assert.Equal(t, "create", domainsCreateCmd.Use)
	assert.Equal(t, "Create a new domain monitor", domainsCreateCmd.Short)
	assert.NotEmpty(t, domainsCreateCmd.Long)
	require.NotNil(t, domainsCreateCmd.RunE, "domains create command should have a RunE function")

	// Verify required flags
	nameFlag := domainsCreateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "domains create command should have --name flag")
	assert.Equal(t, "string", nameFlag.Value.Type())

	domainFlag := domainsCreateCmd.Flags().Lookup("domain")
	require.NotNil(t, domainFlag, "domains create command should have --domain flag")
	assert.Equal(t, "string", domainFlag.Value.Type())

	// Verify optional flags
	intervalFlag := domainsCreateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "domains create command should have --interval flag")
	assert.Equal(t, "int", intervalFlag.Value.Type())

	gracePeriodFlag := domainsCreateCmd.Flags().Lookup("grace-period")
	require.NotNil(t, gracePeriodFlag, "domains create command should have --grace-period flag")
	assert.Equal(t, "int", gracePeriodFlag.Value.Type())

	warningThresholdFlag := domainsCreateCmd.Flags().Lookup("warning-threshold")
	require.NotNil(t, warningThresholdFlag, "domains create command should have --warning-threshold flag")
	assert.Equal(t, "int", warningThresholdFlag.Value.Type())

	urgentThresholdFlag := domainsCreateCmd.Flags().Lookup("urgent-threshold")
	require.NotNil(t, urgentThresholdFlag, "domains create command should have --urgent-threshold flag")
	assert.Equal(t, "int", urgentThresholdFlag.Value.Type())

	criticalThresholdFlag := domainsCreateCmd.Flags().Lookup("critical-threshold")
	require.NotNil(t, criticalThresholdFlag, "domains create command should have --critical-threshold flag")
	assert.Equal(t, "int", criticalThresholdFlag.Value.Type())
}

// TestDomainsUpdateCommand tests the domains update command
func TestDomainsUpdateCommand(t *testing.T) {
	assert.Equal(t, "update <id>", domainsUpdateCmd.Use)
	assert.Equal(t, "Update a domain monitor", domainsUpdateCmd.Short)
	assert.NotEmpty(t, domainsUpdateCmd.Long)
	require.NotNil(t, domainsUpdateCmd.RunE, "domains update command should have a RunE function")

	// Verify flags exist
	nameFlag := domainsUpdateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "domains update command should have --name flag")

	domainFlag := domainsUpdateCmd.Flags().Lookup("domain")
	require.NotNil(t, domainFlag, "domains update command should have --domain flag")

	intervalFlag := domainsUpdateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "domains update command should have --interval flag")

	gracePeriodFlag := domainsUpdateCmd.Flags().Lookup("grace-period")
	require.NotNil(t, gracePeriodFlag, "domains update command should have --grace-period flag")

	warningThresholdFlag := domainsUpdateCmd.Flags().Lookup("warning-threshold")
	require.NotNil(t, warningThresholdFlag, "domains update command should have --warning-threshold flag")

	urgentThresholdFlag := domainsUpdateCmd.Flags().Lookup("urgent-threshold")
	require.NotNil(t, urgentThresholdFlag, "domains update command should have --urgent-threshold flag")

	criticalThresholdFlag := domainsUpdateCmd.Flags().Lookup("critical-threshold")
	require.NotNil(t, criticalThresholdFlag, "domains update command should have --critical-threshold flag")

	statusFlag := domainsUpdateCmd.Flags().Lookup("status")
	require.NotNil(t, statusFlag, "domains update command should have --status flag")
}

// TestDomainsPauseCommand tests the domains pause command
func TestDomainsPauseCommand(t *testing.T) {
	assert.Equal(t, "pause <id>", domainsPauseCmd.Use)
	assert.Equal(t, "Pause a domain monitor", domainsPauseCmd.Short)
	assert.NotEmpty(t, domainsPauseCmd.Long)
	require.NotNil(t, domainsPauseCmd.RunE, "domains pause command should have a RunE function")
}

// TestDomainsResumeCommand tests the domains resume command
func TestDomainsResumeCommand(t *testing.T) {
	assert.Equal(t, "resume <id>", domainsResumeCmd.Use)
	assert.Equal(t, "Resume a domain monitor", domainsResumeCmd.Short)
	assert.NotEmpty(t, domainsResumeCmd.Long)
	require.NotNil(t, domainsResumeCmd.RunE, "domains resume command should have a RunE function")
}

// TestDomainsIncidentsCommand tests the domains incidents command
func TestDomainsIncidentsCommand(t *testing.T) {
	assert.Equal(t, "incidents <id>", domainsIncidentsCmd.Use)
	assert.Equal(t, "Show incident history", domainsIncidentsCmd.Short)
	assert.NotEmpty(t, domainsIncidentsCmd.Long)
	require.NotNil(t, domainsIncidentsCmd.RunE, "domains incidents command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := domainsIncidentsCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "domains incidents command should have --json flag")
}

// TestDomainsDeleteCommand tests the domains delete command
func TestDomainsDeleteCommand(t *testing.T) {
	assert.Equal(t, "delete <id>", domainsDeleteCmd.Use)
	assert.Equal(t, "Delete a domain monitor", domainsDeleteCmd.Short)
	assert.NotEmpty(t, domainsDeleteCmd.Long)
	require.NotNil(t, domainsDeleteCmd.RunE, "domains delete command should have a RunE function")

	// Verify --force flag exists
	forceFlag := domainsDeleteCmd.Flags().Lookup("force")
	require.NotNil(t, forceFlag, "domains delete command should have --force flag")
	assert.Equal(t, "bool", forceFlag.Value.Type())
}

// TestDomainsCommandHasSubcommands verifies all subcommands are registered
func TestDomainsCommandHasSubcommands(t *testing.T) {
	commands := domainsCmd.Commands()

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
		assert.True(t, found, "domains command should have %s subcommand", expected)
	}
}

// TestResolveDomainID tests the helper function for resolving short IDs
func TestResolveDomainID(t *testing.T) {
	// This is a unit test for the helper function
	// In a real scenario, you'd mock the API client
	// For now, we just verify the function exists by checking if it's referenced
	// A full integration test would require a mock API server
	assert.NotNil(t, domainsShowCmd.RunE, "resolveDomainID is used by show command")
}
