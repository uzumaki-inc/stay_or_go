package cmd_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uzumaki-inc/stay_or_go/cmd"
)

// Helper-driven subprocess tests to validate error paths without affecting parent test process.
func TestRootCommand_ErrorScenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		scenario string
		expect   string
	}{
		{name: "no args", scenario: "NOARGS", expect: "Please Enter specify a language"},
		{name: "unsupported language", scenario: "UNSUPPORTED", expect: "Error: Unsupported language"},
		{name: "bad format", scenario: "BADFORMAT", expect: "Error: Unsupported output format"},
		{name: "missing token", scenario: "NOTOKEN", expect: "Please provide a GitHub token"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			capPath := filepath.Join(dir, "stderr.txt")

			//nolint:gosec // launching test subprocess intentionally
			cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestHelperProcess_CobraRoot")
			cmd.Env = append(os.Environ(),
				"GO_WANT_HELPER_PROCESS_COBRA=1",
				"COBRA_SCENARIO="+testCase.scenario,
				"COBRA_CAPTURE="+capPath,
			)

			err := cmd.Run()
			if err == nil {
				t.Fatalf("expected process to exit with error")
			}

			// read captured stderr
			data, readErr := os.ReadFile(capPath)
			if readErr != nil {
				t.Fatalf("failed reading capture: %v", readErr)
			}

			if !strings.Contains(string(data), testCase.expect) {
				t.Fatalf("expected stderr to contain %q, got: %s", testCase.expect, string(data))
			}
		})
	}
}

// Test helper that runs in a subprocess to exercise cobra command paths that call os.Exit.
//
//nolint:paralleltest,funlen // Test helper process for subprocess testing, t unused but required
func TestHelperProcess_CobraRoot(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_COBRA") != "1" {
		return
	}

	// Redirect stderr to capture file
	capPath := os.Getenv("COBRA_CAPTURE")
	f, _ := os.Create(capPath)

	defer f.Close()

	os.Stderr = f

	scenario := os.Getenv("COBRA_SCENARIO")
	goMod := `module example.com

require (
    github.com/replaced/mod v1.0.0
)

replace (
    github.com/replaced/mod v1.0.0 => ./local/mod
)`
	gemfile := `source 'https://rubygems.org'

# git specified to be skipped

gem 'nokogiri', git: 'https://example.com/sparklemotion/nokogiri.git'
`
	handlers := map[string]func(){
		"NOARGS":      func() { cmd.GetRootCmd().SetArgs([]string{}) },
		"UNSUPPORTED": func() { cmd.GetRootCmd().SetArgs([]string{"python"}) },
		"BADFORMAT":   func() { cmd.GetRootCmd().SetArgs([]string{"go", "-f", "json", "-g", "dummy"}) },
		"NOTOKEN":     func() { _ = os.Unsetenv("GITHUB_TOKEN"); cmd.GetRootCmd().SetArgs([]string{"go"}) },
		"GO_DEFAULT": func() {
			dir := t.TempDir()
			_ = os.WriteFile(dir+"/go.mod", []byte(goMod), 0o600)
			t.Chdir(dir)
			t.Setenv("GITHUB_TOKEN", "dummy")
			cmd.GetRootCmd().SetArgs([]string{"go"})
		},
		"RUBY_DEFAULT": func() {
			dir := t.TempDir()
			_ = os.WriteFile(dir+"/Gemfile", []byte(gemfile), 0o600)
			t.Chdir(dir)
			t.Setenv("GITHUB_TOKEN", "dummy")
			cmd.GetRootCmd().SetArgs([]string{"ruby"})
		},
		"GO_VERBOSE": func() {
			dir := t.TempDir()
			_ = os.WriteFile(dir+"/go.mod", []byte(goMod), 0o600)
			t.Chdir(dir)
			t.Setenv("GITHUB_TOKEN", "dummy")
			cmd.GetRootCmd().SetArgs([]string{"go", "-v"})
		},
		"GO_CSV": func() {
			dir := t.TempDir()
			_ = os.WriteFile(dir+"/go.mod", []byte(goMod), 0o600)
			t.Chdir(dir)
			t.Setenv("GITHUB_TOKEN", "dummy")
			cmd.GetRootCmd().SetArgs([]string{"go", "-f", "csv"})
		},
		"GO_CONFIG": func() {
			dir := t.TempDir()
			_ = os.WriteFile(dir+"/go.mod", []byte(goMod), 0o600)
			cfg := dir + "/weights.yml"
			content := "watestCasehers: 1\n" +
				"stars: 2\n" +
				"forks: 3\n" +
				"open_issues: 4\n" +
				"last_commit_date: -5\n" +
				"archived: -6\n"
			_ = os.WriteFile(cfg, []byte(content), 0o600)
			t.Chdir(dir)
			t.Setenv("GITHUB_TOKEN", "dummy")
			cmd.GetRootCmd().SetArgs([]string{"go", "-c", cfg})
		},
	}

	if h, ok := handlers[scenario]; ok {
		h()
	} else {
		cmd.GetRootCmd().SetArgs([]string{})
	}

	// Avoid printing to stdout in tests; ensure buffer present
	var devnull bytes.Buffer

	_ = devnull

	// This will call os.Exit in error paths, terminating subprocess with code 1.
	cmd.Execute()
}

func TestRootCommand_DefaultInputsAndVerbose(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                 string
		scenario             string
		expectErr            bool
		expectStderrContains string
	}{
		{name: "go default input", scenario: "GO_DEFAULT", expectErr: false},
		{name: "ruby default input", scenario: "RUBY_DEFAULT", expectErr: false},
		{name: "go verbose logs", scenario: "GO_VERBOSE", expectErr: false, expectStderrContains: "Selected Language: go"},
		{name: "go with csv format", scenario: "GO_CSV", expectErr: false},
		{name: "go with config file", scenario: "GO_CONFIG", expectErr: false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			capPath := filepath.Join(dir, "stderr.txt")

			//nolint:gosec // launching test subprocess intentionally
			cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestHelperProcess_CobraRoot")
			cmd.Env = append(os.Environ(),
				"GO_WANT_HELPER_PROCESS_COBRA=1",
				"COBRA_SCENARIO="+testCase.scenario,
				"COBRA_CAPTURE="+capPath,
			)

			err := cmd.Run()
			if testCase.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if testCase.expectStderrContains != "" {
				stderrData, readErr := os.ReadFile(capPath)
				if readErr != nil {
					t.Fatalf("failed to read stderr capture: %v", readErr)
				}

				if !strings.Contains(string(stderrData), testCase.expectStderrContains) {
					t.Fatalf("stderr missing expected text %q, got: %s", testCase.expectStderrContains, string(stderrData))
				}
			}
		})
	}
}
