package analyzer_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uzumaki-inc/stay_or_go/analyzer"
)

func TestNewParameterWeightsDefaults(t *testing.T) {
	t.Parallel()

	weights := analyzer.NewParameterWeights()

	// Default weights should match implementation defaults
	assert.InDelta(t, 0.1, weights.Watchers, 0.0001)
	assert.InDelta(t, 0.1, weights.Stars, 0.0001)
	assert.InDelta(t, 0.1, weights.Forks, 0.0001)
	assert.InDelta(t, 0.01, weights.OpenIssues, 0.0001)
	assert.InDelta(t, -0.05, weights.LastCommitDate, 0.0001)
	assert.InDelta(t, -1000000.0, weights.Archived, 0.1)
}

func TestNewParameterWeightsFromConfiFile_LoadsValues(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "weights.yml")
	content := []byte(
		"watchers: 1.5\n" +
			"stars: 2.5\n" +
			"forks: 3.5\n" +
			"open_issues: 4.5\n" +
			"last_commit_date: -6.5\n" +
			"archived: -99999\n",
	)

	err := os.WriteFile(path, content, 0o600)
	if err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	weights := analyzer.NewParameterWeightsFromConfiFile(path)

	assert.InDelta(t, 1.5, weights.Watchers, 0.0001)
	assert.InDelta(t, 2.5, weights.Stars, 0.0001)
	assert.InDelta(t, 3.5, weights.Forks, 0.0001)
	assert.InDelta(t, 4.5, weights.OpenIssues, 0.0001)
	assert.InDelta(t, -6.5, weights.LastCommitDate, 0.0001)
	assert.InDelta(t, -99999.0, weights.Archived, 0.0001)
}

// Test that an invalid path leads to os.Exit(1). Use helper process pattern.
func TestNewParameterWeightsFromConfiFile_ExitOnMissing(t *testing.T) {
	t.Parallel()

	//nolint:gosec // launching test subprocess intentionally
	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestHelperProcess_WeightsExit")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS_WEIGHTS=1")

	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected non-nil error (exit), got nil")
	}

	var exitErr *exec.ExitError

	if errors.As(err, &exitErr) {
		if exitErr.ExitCode() != 1 {
			t.Fatalf("expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Fatalf("expected *exec.ExitError, got %T", err)
	}
}

// Helper for the sub-process exit test.
//
//nolint:paralleltest // Test helper process for subprocess testing
func TestHelperProcess_WeightsExit(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_WEIGHTS") != "1" {
		return
	}
	// This should trigger os.Exit(1) internally
	_ = analyzer.NewParameterWeightsFromConfiFile("/path/does/not/exist.yml")

	t.Fatalf("should have exited before reaching here")
}

func TestNewParameterWeightsFromReader_LoadsValues(t *testing.T) {
	t.Parallel()

	// Create YAML content
	yamlContent := `
watchers: 0.2
stars: 0.3
forks: 0.15
open_issues: 0.02
last_commit_date: -0.1
archived: -2000000
`

	reader := strings.NewReader(yamlContent)

	weights, err := analyzer.NewParameterWeightsFromReader(reader)

	require.NoError(t, err)
	assert.InDelta(t, 0.2, weights.Watchers, 0.0001)
	assert.InDelta(t, 0.3, weights.Stars, 0.0001)
	assert.InDelta(t, 0.15, weights.Forks, 0.0001)
	assert.InDelta(t, 0.02, weights.OpenIssues, 0.0001)
	assert.InDelta(t, -0.1, weights.LastCommitDate, 0.0001)
	assert.InDelta(t, -2000000.0, weights.Archived, 0.1)
}

func TestNewParameterWeightsFromReader_ErrorOnInvalidYAML(t *testing.T) {
	t.Parallel()

	invalidYAML := "invalid: yaml: content: ["
	reader := strings.NewReader(invalidYAML)

	_, err := analyzer.NewParameterWeightsFromReader(reader)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}

func TestNewParameterWeightsFromFile_LoadsValues(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "weights.yml")
	content := []byte(
		"watchers: 1.5\n" +
			"stars: 2.5\n" +
			"forks: 3.5\n" +
			"open_issues: 4.5\n" +
			"last_commit_date: -6.5\n" +
			"archived: -99999\n",
	)

	err := os.WriteFile(path, content, 0o600)
	require.NoError(t, err)

	weights, err := analyzer.NewParameterWeightsFromFile(path)

	require.NoError(t, err)
	assert.InDelta(t, 1.5, weights.Watchers, 0.0001)
	assert.InDelta(t, 2.5, weights.Stars, 0.0001)
	assert.InDelta(t, 3.5, weights.Forks, 0.0001)
	assert.InDelta(t, 4.5, weights.OpenIssues, 0.0001)
	assert.InDelta(t, -6.5, weights.LastCommitDate, 0.0001)
	assert.InDelta(t, -99999.0, weights.Archived, 0.0001)
}

func TestNewParameterWeightsFromFile_ErrorOnMissingFile(t *testing.T) {
	t.Parallel()

	_, err := analyzer.NewParameterWeightsFromFile("/path/does/not/exist.yml")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}
