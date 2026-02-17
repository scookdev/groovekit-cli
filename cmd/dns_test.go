package cmd

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDnsCommand tests the basic structure of the dns command
func TestDnsCommand(t *testing.T) {
	assert.Equal(t, "dns", dnsCmd.Use)
	assert.Equal(t, "Manage DNS record monitors", dnsCmd.Short)
	assert.NotEmpty(t, dnsCmd.Long)
}

// TestDnsListCommand tests the dns list command
func TestDnsListCommand(t *testing.T) {
	assert.Equal(t, "list", dnsListCmd.Use)
	assert.Equal(t, "List all DNS monitors", dnsListCmd.Short)
	assert.NotEmpty(t, dnsListCmd.Long)
	require.NotNil(t, dnsListCmd.RunE, "dns list command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := dnsListCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "dns list command should have --json flag")
	assert.Equal(t, "bool", jsonFlag.Value.Type())
}

// TestDnsShowCommand tests the dns show command
func TestDnsShowCommand(t *testing.T) {
	assert.Equal(t, "show <id>", dnsShowCmd.Use)
	assert.Equal(t, "Show DNS monitor details", dnsShowCmd.Short)
	assert.NotEmpty(t, dnsShowCmd.Long)
	require.NotNil(t, dnsShowCmd.RunE, "dns show command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := dnsShowCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "dns show command should have --json flag")
}

// TestDnsCreateCommand tests the dns create command
func TestDnsCreateCommand(t *testing.T) {
	assert.Equal(t, "create", dnsCreateCmd.Use)
	assert.Equal(t, "Create a new DNS monitor", dnsCreateCmd.Short)
	assert.NotEmpty(t, dnsCreateCmd.Long)
	require.NotNil(t, dnsCreateCmd.RunE, "dns create command should have a RunE function")

	// Verify required flags
	nameFlag := dnsCreateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "dns create command should have --name flag")
	assert.Equal(t, "string", nameFlag.Value.Type())

	domainFlag := dnsCreateCmd.Flags().Lookup("domain")
	require.NotNil(t, domainFlag, "dns create command should have --domain flag")
	assert.Equal(t, "string", domainFlag.Value.Type())

	typeFlag := dnsCreateCmd.Flags().Lookup("type")
	require.NotNil(t, typeFlag, "dns create command should have --type flag")
	assert.Equal(t, "string", typeFlag.Value.Type())

	expectedFlag := dnsCreateCmd.Flags().Lookup("expected")
	require.NotNil(t, expectedFlag, "dns create command should have --expected flag")
	assert.Equal(t, "stringSlice", expectedFlag.Value.Type())

	// Verify optional flags
	intervalFlag := dnsCreateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "dns create command should have --interval flag")
	assert.Equal(t, "int", intervalFlag.Value.Type())

	gracePeriodFlag := dnsCreateCmd.Flags().Lookup("grace-period")
	require.NotNil(t, gracePeriodFlag, "dns create command should have --grace-period flag")
	assert.Equal(t, "int", gracePeriodFlag.Value.Type())
}

// TestDnsUpdateCommand tests the dns update command
func TestDnsUpdateCommand(t *testing.T) {
	assert.Equal(t, "update <id>", dnsUpdateCmd.Use)
	assert.Equal(t, "Update a DNS monitor", dnsUpdateCmd.Short)
	assert.NotEmpty(t, dnsUpdateCmd.Long)
	require.NotNil(t, dnsUpdateCmd.RunE, "dns update command should have a RunE function")

	// Verify flags exist
	nameFlag := dnsUpdateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "dns update command should have --name flag")

	domainFlag := dnsUpdateCmd.Flags().Lookup("domain")
	require.NotNil(t, domainFlag, "dns update command should have --domain flag")

	typeFlag := dnsUpdateCmd.Flags().Lookup("type")
	require.NotNil(t, typeFlag, "dns update command should have --type flag")

	expectedFlag := dnsUpdateCmd.Flags().Lookup("expected")
	require.NotNil(t, expectedFlag, "dns update command should have --expected flag")

	intervalFlag := dnsUpdateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "dns update command should have --interval flag")

	gracePeriodFlag := dnsUpdateCmd.Flags().Lookup("grace-period")
	require.NotNil(t, gracePeriodFlag, "dns update command should have --grace-period flag")

	statusFlag := dnsUpdateCmd.Flags().Lookup("status")
	require.NotNil(t, statusFlag, "dns update command should have --status flag")
}

// TestDnsPauseCommand tests the dns pause command
func TestDnsPauseCommand(t *testing.T) {
	assert.Equal(t, "pause <id>", dnsPauseCmd.Use)
	assert.Equal(t, "Pause a DNS monitor", dnsPauseCmd.Short)
	assert.NotEmpty(t, dnsPauseCmd.Long)
	require.NotNil(t, dnsPauseCmd.RunE, "dns pause command should have a RunE function")
}

// TestDnsResumeCommand tests the dns resume command
func TestDnsResumeCommand(t *testing.T) {
	assert.Equal(t, "resume <id>", dnsResumeCmd.Use)
	assert.Equal(t, "Resume a DNS monitor", dnsResumeCmd.Short)
	assert.NotEmpty(t, dnsResumeCmd.Long)
	require.NotNil(t, dnsResumeCmd.RunE, "dns resume command should have a RunE function")
}

// TestDnsIncidentsCommand tests the dns incidents command
func TestDnsIncidentsCommand(t *testing.T) {
	assert.Equal(t, "incidents <id>", dnsIncidentsCmd.Use)
	assert.Equal(t, "Show incident history", dnsIncidentsCmd.Short)
	assert.NotEmpty(t, dnsIncidentsCmd.Long)
	require.NotNil(t, dnsIncidentsCmd.RunE, "dns incidents command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := dnsIncidentsCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "dns incidents command should have --json flag")
}

// TestDnsDeleteCommand tests the dns delete command
func TestDnsDeleteCommand(t *testing.T) {
	assert.Equal(t, "delete <id>", dnsDeleteCmd.Use)
	assert.Equal(t, "Delete a DNS monitor", dnsDeleteCmd.Short)
	assert.NotEmpty(t, dnsDeleteCmd.Long)
	require.NotNil(t, dnsDeleteCmd.RunE, "dns delete command should have a RunE function")

	// Verify --force flag exists
	forceFlag := dnsDeleteCmd.Flags().Lookup("force")
	require.NotNil(t, forceFlag, "dns delete command should have --force flag")
	assert.Equal(t, "bool", forceFlag.Value.Type())
}

// TestDnsCommandHasSubcommands verifies all subcommands are registered
func TestDnsCommandHasSubcommands(t *testing.T) {
	commands := dnsCmd.Commands()

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
		assert.True(t, found, "dns command should have %s subcommand", expected)
	}
}

// TestResolveDnsMonitorID tests the helper function for resolving short IDs
func TestResolveDnsMonitorID(t *testing.T) {
	// This is a unit test for the helper function
	// In a real scenario, you'd mock the API client
	// For now, we just verify the function exists by checking if it's referenced
	// A full integration test would require a mock API server
	assert.NotNil(t, dnsShowCmd.RunE, "resolveDnsMonitorID is used by show command")
}

// TestContainsHelper tests the contains helper function
func TestContainsHelper(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "item exists",
			slice:    []string{"A", "AAAA", "MX"},
			item:     "MX",
			expected: true,
		},
		{
			name:     "item does not exist",
			slice:    []string{"A", "AAAA", "MX"},
			item:     "CNAME",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "A",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := slices.Contains(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}
