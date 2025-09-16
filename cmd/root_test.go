package cmd

import (
    "bytes"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
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

    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            dir := t.TempDir()
            capPath := filepath.Join(dir, "stderr.txt")

            cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess_CobraRoot")
            cmd.Env = append(os.Environ(),
                "GO_WANT_HELPER_PROCESS_COBRA=1",
                "COBRA_SCENARIO="+tc.scenario,
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
            if !strings.Contains(string(data), tc.expect) {
                t.Fatalf("expected stderr to contain %q, got: %s", tc.expect, string(data))
            }
        })
    }
}

// Test helper that runs in a subprocess to exercise cobra command paths that call os.Exit.
func TestHelperProcess_CobraRoot(t *testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS_COBRA") != "1" {
        return
    }

    // Redirect stderr to capture file
    capPath := os.Getenv("COBRA_CAPTURE")
    f, _ := os.Create(capPath)
    defer f.Close()
    os.Stderr = f

    switch os.Getenv("COBRA_SCENARIO") {
    case "NOARGS":
        rootCmd.SetArgs([]string{})
    case "UNSUPPORTED":
        rootCmd.SetArgs([]string{"python"})
    case "BADFORMAT":
        rootCmd.SetArgs([]string{"go", "-f", "json", "-g", "dummy"})
    case "NOTOKEN":
        // ensure token not provided in flag nor env
        _ = os.Unsetenv("GITHUB_TOKEN")
        rootCmd.SetArgs([]string{"go"})
    case "GO_DEFAULT":
        // create temp go.mod using default path and run successfully
        dir, _ := os.MkdirTemp("", "gomod-default-*")
        _ = os.WriteFile(dir+"/go.mod", []byte("module example.com\n\nrequire (\n    github.com/replaced/mod v1.0.0\n)\n\nreplace (\n    github.com/replaced/mod v1.0.0 => ./local/mod\n)\n"), 0o600)
        _ = os.Chdir(dir)
        _ = os.Setenv("GITHUB_TOKEN", "dummy")
        rootCmd.SetArgs([]string{"go"})
    case "RUBY_DEFAULT":
        dir, _ := os.MkdirTemp("", "gem-default-*")
        _ = os.WriteFile(dir+"/Gemfile", []byte("source 'https://rubygems.org'\n\n# git指定でSkipさせる\n\ngem 'nokogiri', git: 'https://example.com/sparklemotion/nokogiri.git'\n"), 0o600)
        _ = os.Chdir(dir)
        _ = os.Setenv("GITHUB_TOKEN", "dummy")
        rootCmd.SetArgs([]string{"ruby"})
    case "GO_VERBOSE":
        dir, _ := os.MkdirTemp("", "gomod-verbose-*")
        _ = os.WriteFile(dir+"/go.mod", []byte("module example.com\n\nrequire (\n    github.com/replaced/mod v1.0.0\n)\n\nreplace (\n    github.com/replaced/mod v1.0.0 => ./local/mod\n)\n"), 0o600)
        _ = os.Chdir(dir)
        _ = os.Setenv("GITHUB_TOKEN", "dummy")
        rootCmd.SetArgs([]string{"go", "-v"})
    case "GO_CSV":
        dir, _ := os.MkdirTemp("", "gomod-csv-*")
        _ = os.WriteFile(dir+"/go.mod", []byte("module example.com\n\nrequire (\n    github.com/replaced/mod v1.0.0\n)\n\nreplace (\n    github.com/replaced/mod v1.0.0 => ./local/mod\n)\n"), 0o600)
        _ = os.Chdir(dir)
        _ = os.Setenv("GITHUB_TOKEN", "dummy")
        rootCmd.SetArgs([]string{"go", "-f", "csv"})
    case "GO_CONFIG":
        dir, _ := os.MkdirTemp("", "gomod-config-*")
        _ = os.WriteFile(dir+"/go.mod", []byte("module example.com\n\nrequire (\n    github.com/replaced/mod v1.0.0\n)\n\nreplace (\n    github.com/replaced/mod v1.0.0 => ./local/mod\n)\n"), 0o600)
        cfg := dir+"/weights.yml"
        _ = os.WriteFile(cfg, []byte("watchers: 1\nstars: 2\nforks: 3\nopen_issues: 4\nlast_commit_date: -5\narchived: -6\n"), 0o600)
        _ = os.Chdir(dir)
        _ = os.Setenv("GITHUB_TOKEN", "dummy")
        rootCmd.SetArgs([]string{"go", "-c", cfg})
    default:
        rootCmd.SetArgs([]string{})
    }

    // Avoid printing to stdout in tests; ensure buffer present
    var devnull bytes.Buffer
    _ = devnull

    // This will call os.Exit in error paths, terminating subprocess with code 1.
    Execute()

    // Should not reach here
}

func TestRootCommand_DefaultInputsAndVerbose(t *testing.T) {
    t.Parallel()

    cases := []struct {
        name     string
        scenario string
        expectErr bool
        expectStderrContains string
    }{
        {name: "go default input", scenario: "GO_DEFAULT", expectErr: false},
        {name: "ruby default input", scenario: "RUBY_DEFAULT", expectErr: false},
        {name: "go verbose logs", scenario: "GO_VERBOSE", expectErr: false, expectStderrContains: "Selected Language: go"},
        {name: "go with csv format", scenario: "GO_CSV", expectErr: false},
        {name: "go with config file", scenario: "GO_CONFIG", expectErr: false},
    }

    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            dir := t.TempDir()
            capPath := filepath.Join(dir, "stderr.txt")

            cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess_CobraRoot")
            cmd.Env = append(os.Environ(),
                "GO_WANT_HELPER_PROCESS_COBRA=1",
                "COBRA_SCENARIO="+tc.scenario,
                "COBRA_CAPTURE="+capPath,
            )

            err := cmd.Run()
            if tc.expectErr {
                if err == nil {
                    t.Fatalf("expected error, got nil")
                }
            } else {
                if err != nil {
                    t.Fatalf("unexpected error: %v", err)
                }
            }

            if tc.expectStderrContains != "" {
                b, readErr := os.ReadFile(capPath)
                if readErr != nil {
                    t.Fatalf("failed to read stderr capture: %v", readErr)
                }
                if !strings.Contains(string(b), tc.expectStderrContains) {
                    t.Fatalf("stderr missing expected text %q, got: %s", tc.expectStderrContains, string(b))
                }
            }
        })
    }
}
