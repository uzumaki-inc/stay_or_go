package analyzer

import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewParameterWeightsDefaults(t *testing.T) {
    t.Parallel()

    w := NewParameterWeights()

    assert.InDelta(t, float64(defaultWatcherWeight), w.Watchers, 0.0001)
    assert.InDelta(t, float64(defaultStarWeight), w.Stars, 0.0001)
    assert.InDelta(t, float64(defaultForkWeight), w.Forks, 0.0001)
    assert.InDelta(t, float64(defaultOpenIssueWeight), w.OpenIssues, 0.0001)
    assert.InDelta(t, float64(defaultLastCommitDateWeight), w.LastCommitDate, 0.0001)
    assert.InDelta(t, float64(defaultArchivedWeight), w.Archived, 0.0001)
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
    if err := os.WriteFile(path, content, 0o600); err != nil {
        t.Fatalf("failed to write temp config: %v", err)
    }

    w := NewParameterWeightsFromConfiFile(path)

    assert.InDelta(t, 1.5, w.Watchers, 0.0001)
    assert.InDelta(t, 2.5, w.Stars, 0.0001)
    assert.InDelta(t, 3.5, w.Forks, 0.0001)
    assert.InDelta(t, 4.5, w.OpenIssues, 0.0001)
    assert.InDelta(t, -6.5, w.LastCommitDate, 0.0001)
    assert.InDelta(t, -99999.0, w.Archived, 0.0001)
}

// Test that an invalid path leads to os.Exit(1). Use helper process pattern.
func TestNewParameterWeightsFromConfiFile_ExitOnMissing(t *testing.T) {
    t.Parallel()

    cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess_WeightsExit")
    cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS_WEIGHTS=1")

    err := cmd.Run()
    if err == nil {
        t.Fatalf("expected non-nil error (exit), got nil")
    }
    if exitErr, ok := err.(*exec.ExitError); ok {
        if exitErr.ExitCode() != 1 {
            t.Fatalf("expected exit code 1, got %d", exitErr.ExitCode())
        }
    } else {
        t.Fatalf("expected *exec.ExitError, got %T", err)
    }
}

// Helper for the sub-process exit test.
func TestHelperProcess_WeightsExit(t *testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS_WEIGHTS") != "1" {
        return
    }
    // This should trigger os.Exit(1) internally
    _ = NewParameterWeightsFromConfiFile("/path/does/not/exist.yml")
    t.Fatalf("should have exited before reaching here")
}
