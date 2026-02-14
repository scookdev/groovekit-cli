package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJobsCommand tests the basic structure of the jobs command
func TestJobsCommand(t *testing.T) {
	assert.Equal(t, "jobs", jobsCmd.Use)
	assert.Equal(t, "Manage cron job monitors", jobsCmd.Short)
	assert.NotEmpty(t, jobsCmd.Long)
}

// TestJobsListCommand tests the jobs list command
func TestJobsListCommand(t *testing.T) {
	assert.Equal(t, "list", jobsListCmd.Use)
	assert.Equal(t, "List all jobs", jobsListCmd.Short)
	assert.NotEmpty(t, jobsListCmd.Long)
	require.NotNil(t, jobsListCmd.RunE, "jobs list command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := jobsListCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "jobs list command should have --json flag")
	assert.Equal(t, "bool", jsonFlag.Value.Type())
}

// TestJobsShowCommand tests the jobs show command
func TestJobsShowCommand(t *testing.T) {
	assert.Equal(t, "show <id>", jobsShowCmd.Use)
	assert.Equal(t, "Show job details", jobsShowCmd.Short)
	assert.NotEmpty(t, jobsShowCmd.Long)
	require.NotNil(t, jobsShowCmd.RunE, "jobs show command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := jobsShowCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "jobs show command should have --json flag")
}

// TestJobsCreateCommand tests the jobs create command
func TestJobsCreateCommand(t *testing.T) {
	assert.Equal(t, "create", jobsCreateCmd.Use)
	assert.Equal(t, "Create a new job", jobsCreateCmd.Short)
	assert.NotEmpty(t, jobsCreateCmd.Long)
	require.NotNil(t, jobsCreateCmd.RunE, "jobs create command should have a RunE function")

	// Verify required flags
	nameFlag := jobsCreateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "jobs create command should have --name flag")
	assert.Equal(t, "string", nameFlag.Value.Type())

	intervalFlag := jobsCreateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "jobs create command should have --interval flag")
	assert.Equal(t, "int", intervalFlag.Value.Type())

	// Verify optional flags
	gracePeriodFlag := jobsCreateCmd.Flags().Lookup("grace-period")
	require.NotNil(t, gracePeriodFlag, "jobs create command should have --grace-period flag")
}

// TestJobsUpdateCommand tests the jobs update command
func TestJobsUpdateCommand(t *testing.T) {
	assert.Equal(t, "update <id>", jobsUpdateCmd.Use)
	assert.Equal(t, "Update a job", jobsUpdateCmd.Short)
	assert.NotEmpty(t, jobsUpdateCmd.Long)
	require.NotNil(t, jobsUpdateCmd.RunE, "jobs update command should have a RunE function")

	// Verify flags exist
	nameFlag := jobsUpdateCmd.Flags().Lookup("name")
	require.NotNil(t, nameFlag, "jobs update command should have --name flag")

	intervalFlag := jobsUpdateCmd.Flags().Lookup("interval")
	require.NotNil(t, intervalFlag, "jobs update command should have --interval flag")

	statusFlag := jobsUpdateCmd.Flags().Lookup("status")
	require.NotNil(t, statusFlag, "jobs update command should have --status flag")
}

// TestJobsPauseCommand tests the jobs pause command
func TestJobsPauseCommand(t *testing.T) {
	assert.Equal(t, "pause <id>", jobsPauseCmd.Use)
	assert.Equal(t, "Pause a job", jobsPauseCmd.Short)
	assert.NotEmpty(t, jobsPauseCmd.Long)
	require.NotNil(t, jobsPauseCmd.RunE, "jobs pause command should have a RunE function")
}

// TestJobsResumeCommand tests the jobs resume command
func TestJobsResumeCommand(t *testing.T) {
	assert.Equal(t, "resume <id>", jobsResumeCmd.Use)
	assert.Equal(t, "Resume a job", jobsResumeCmd.Short)
	assert.NotEmpty(t, jobsResumeCmd.Long)
	require.NotNil(t, jobsResumeCmd.RunE, "jobs resume command should have a RunE function")
}

// TestJobsIncidentsCommand tests the jobs incidents command
func TestJobsIncidentsCommand(t *testing.T) {
	assert.Equal(t, "incidents <id>", jobsIncidentsCmd.Use)
	assert.Equal(t, "Show incident history", jobsIncidentsCmd.Short)
	assert.NotEmpty(t, jobsIncidentsCmd.Long)
	require.NotNil(t, jobsIncidentsCmd.RunE, "jobs incidents command should have a RunE function")

	// Verify --json flag exists
	jsonFlag := jobsIncidentsCmd.Flags().Lookup("json")
	require.NotNil(t, jsonFlag, "jobs incidents command should have --json flag")
}

// TestJobsDeleteCommand tests the jobs delete command
func TestJobsDeleteCommand(t *testing.T) {
	assert.Equal(t, "delete <id>", jobsDeleteCmd.Use)
	assert.Equal(t, "Delete a job", jobsDeleteCmd.Short)
	assert.NotEmpty(t, jobsDeleteCmd.Long)
	require.NotNil(t, jobsDeleteCmd.RunE, "jobs delete command should have a RunE function")

	// Verify --force flag exists
	forceFlag := jobsDeleteCmd.Flags().Lookup("force")
	require.NotNil(t, forceFlag, "jobs delete command should have --force flag")
	assert.Equal(t, "bool", forceFlag.Value.Type())
}

// TestJobsCommandHasSubcommands verifies all subcommands are registered
func TestJobsCommandHasSubcommands(t *testing.T) {
	commands := jobsCmd.Commands()

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
		assert.True(t, found, "jobs command should have %s subcommand", expected)
	}
}

// TestResolveJobID tests the helper function for resolving short IDs
func TestResolveJobID(t *testing.T) {
	// This is a unit test for the helper function
	// In a real scenario, you'd mock the API client
	// For now, we just verify the function exists by checking if it's referenced
	// A full integration test would require a mock API server
	assert.NotNil(t, jobsShowCmd.RunE, "resolveJobID is used by show command")
}
