//nolint:testpackage // Tests unexported functions
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/uzumaki-inc/stay_or_go/analyzer"
	"github.com/uzumaki-inc/stay_or_go/parser"
	"github.com/uzumaki-inc/stay_or_go/presenter"
)

// Stubs
type stubAnalyzer struct{ called bool }

func (s *stubAnalyzer) FetchGithubInfo(_ []string) []analyzer.GitHubRepoInfo {
	s.called = true

	return []analyzer.GitHubRepoInfo{{GithubRepoURL: "https://github.com/u/a"}}
}

type recorderParser struct {
	lastFile string
	list     []parser.LibInfo
}

func (r *recorderParser) Parse(file string) ([]parser.LibInfo, error) {
	r.lastFile = file

	return r.list, nil
}
func (r *recorderParser) GetRepositoryURL(list []parser.LibInfo) []parser.LibInfo { return list }

type recorderPresenter struct{ displayed bool }

func (r *recorderPresenter) Display() { r.displayed = true }

func TestRun_Success_Go_DefaultFile(t *testing.T) {
	t.Parallel()
	// Prepare deps
	recParser := &recorderParser{list: []parser.LibInfo{{Name: "a", RepositoryURL: "https://github.com/u/a"}}}
	stubAnal := &stubAnalyzer{}
	recPresenter := &recorderPresenter{}

	deps := Deps{
		NewAnalyzer:     func(_ string, _ analyzer.ParameterWeights) AnalyzerPort { return stubAnal },
		SelectParser:    func(_ string) (parser.Parser, error) { return recParser, nil },
		SelectPresenter: func(_ string, _ []presenter.AnalyzedLibInfo) PresenterPort { return recPresenter },
	}

	// Unset env to ensure token from argument is used
	_ = os.Unsetenv("GITHUB_TOKEN")

	err := run("go", "", "markdown", "tok", "", false, deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if recParser.lastFile != "go.mod" {
		t.Fatalf("expected default go.mod, got %q", recParser.lastFile)
	}

	if !stubAnal.called {
		t.Fatalf("expected analyzer called")
	}

	if !recPresenter.displayed {
		t.Fatalf("expected presenter.Display called")
	}
}

func TestRun_Success_Ruby_DefaultFile(t *testing.T) {
	t.Parallel()

	recParser := &recorderParser{list: []parser.LibInfo{{Name: "a", RepositoryURL: "https://github.com/u/a"}}}
	stubAnal := &stubAnalyzer{}
	recPresenter := &recorderPresenter{}

	deps := Deps{
		NewAnalyzer:     func(_ string, _ analyzer.ParameterWeights) AnalyzerPort { return stubAnal },
		SelectParser:    func(_ string) (parser.Parser, error) { return recParser, nil },
		SelectPresenter: func(_ string, _ []presenter.AnalyzedLibInfo) PresenterPort { return recPresenter },
	}
	_ = os.Unsetenv("GITHUB_TOKEN")

	err := run("ruby", "", "markdown", "tok", "", false, deps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if recParser.lastFile != "Gemfile" {
		t.Fatalf("expected default Gemfile, got %q", recParser.lastFile)
	}
}

func TestRun_SkipAll_DoesNotCallAnalyzer(t *testing.T) {
	t.Parallel()

	recParser := &recorderParser{list: []parser.LibInfo{{Name: "a", Skip: true, SkipReason: "skip"}}}
	called := false
	recPresenter := &recorderPresenter{}

	deps := Deps{
		NewAnalyzer:     func(_ string, _ analyzer.ParameterWeights) AnalyzerPort { return &stubAnalyzer{called: true} },
		SelectParser:    func(_ string) (parser.Parser, error) { return recParser, nil },
		SelectPresenter: func(_ string, _ []presenter.AnalyzedLibInfo) PresenterPort { return recPresenter },
	}

	// Wrap NewAnalyzer to detect if it's used later via FetchGithubInfo
	deps.NewAnalyzer = func(_ string, _ analyzer.ParameterWeights) AnalyzerPort {
		return AnalyzerPort(rtFuncAnalyzer(func(_ []string) []analyzer.GitHubRepoInfo {
			called = true

			return nil
		}))
	}

	_ = os.Unsetenv("GITHUB_TOKEN")

	err := run("go", "", "markdown", "tok", "", false, deps)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	if called {
		t.Fatalf("analyzer should not be called when all skipped")
	}

	if !recPresenter.displayed {
		t.Fatalf("presenter should be called even when skipped")
	}
}

// Analyzer adapter via function for testing
type rtFuncAnalyzer func([]string) []analyzer.GitHubRepoInfo

func (f rtFuncAnalyzer) FetchGithubInfo(urls []string) []analyzer.GitHubRepoInfo { return f(urls) }

func TestRun_UnsupportedAndFormatAndTokenErrors(t *testing.T) {
	t.Parallel()

	deps := Deps{}

	err := run("python", "", "markdown", "tok", "", false, deps)
	if err == nil {
		t.Fatalf("expected unsupported language error")
	}

	err = run("go", "", "json", "tok", "", false, deps)
	if err == nil {
		t.Fatalf("expected unsupported format error")
	}

	_ = os.Unsetenv("GITHUB_TOKEN")

	err = run("go", "", "markdown", "", "", false, deps)
	if err == nil {
		t.Fatalf("expected missing token error")
	}
}

func TestRun_WithConfigFileBranch(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfg := filepath.Join(dir, "weights.yml")

	err := os.WriteFile(cfg, []byte("watchers: 1\n"), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	recParser := &recorderParser{list: []parser.LibInfo{{Name: "a", RepositoryURL: "https://github.com/u/a"}}}
	recPresenter := &recorderPresenter{}

	deps := Deps{
		NewAnalyzer:     func(_ string, _ analyzer.ParameterWeights) AnalyzerPort { return &stubAnalyzer{} },
		SelectParser:    func(_ string) (parser.Parser, error) { return recParser, nil },
		SelectPresenter: func(_ string, _ []presenter.AnalyzedLibInfo) PresenterPort { return recPresenter },
	}
	_ = os.Unsetenv("GITHUB_TOKEN")

	err = run("go", "", "markdown", "tok", cfg, false, deps)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	if !recPresenter.displayed {
		t.Fatalf("expected presenter called")
	}
}
