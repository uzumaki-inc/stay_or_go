package presenter

import (
	"testing"

	"github.com/uzumaki-inc/stay_or_go/analyzer"
	"github.com/uzumaki-inc/stay_or_go/parser"
)

func TestAnalyzedLibInfo_EmptyLibInfo_NoRepoInfo(t *testing.T) {
	t.Parallel()

	info := AnalyzedLibInfo{LibInfo: &parser.LibInfo{Skip: true, SkipReason: "li"}, GitHubRepoInfo: nil}

	if info.Name() != nil {
		t.Fatalf("expected nil Name")
	}
	if info.RepositoryURL() != nil {
		t.Fatalf("expected nil RepositoryURL")
	}
	if info.Watchers() != nil {
		t.Fatalf("expected nil Watchers")
	}
	if info.Stars() != nil {
		t.Fatalf("expected nil Stars")
	}
	if info.Forks() != nil {
		t.Fatalf("expected nil Forks")
	}
	if info.OpenIssues() != nil {
		t.Fatalf("expected nil OpenIssues")
	}
	if info.LastCommitDate() != nil {
		t.Fatalf("expected nil LastCommitDate")
	}
	if info.GithubRepoURL() != nil {
		t.Fatalf("expected nil GithubRepoURL")
	}
	if info.Archived() != nil {
		t.Fatalf("expected nil Archived")
	}
	if info.Score() != nil {
		t.Fatalf("expected nil Score")
	}
	if info.Skip() == nil || *info.Skip() != true {
		t.Fatalf("expected Skip=true from LibInfo")
	}
	if v := info.SkipReason(); v == nil || *v != "li" {
		t.Fatalf("expected SkipReason from LibInfo")
	}
}

func TestAnalyzedLibInfo_WithValues_AllGetters(t *testing.T) {
	t.Parallel()

	lib := parser.LibInfo{Name: "lib", RepositoryURL: "https://github.com/x/y"}
	repo := analyzer.GitHubRepoInfo{
		Watchers:       1,
		Stars:          2,
		Forks:          3,
		OpenIssues:     4,
		LastCommitDate: "2024-01-01T00:00:00Z",
		GithubRepoURL:  "https://github.com/x/y",
		Archived:       true,
		Score:          42,
	}
	info := AnalyzedLibInfo{LibInfo: &lib, GitHubRepoInfo: &repo}

	if v := info.Name(); v == nil || *v != "lib" {
		t.Fatalf("unexpected Name")
	}
	if v := info.RepositoryURL(); v == nil || *v != "https://github.com/x/y" {
		t.Fatalf("unexpected RepositoryURL")
	}
	if v := info.Watchers(); v == nil || *v != 1 {
		t.Fatalf("unexpected Watchers")
	}
	if v := info.Stars(); v == nil || *v != 2 {
		t.Fatalf("unexpected Stars")
	}
	if v := info.Forks(); v == nil || *v != 3 {
		t.Fatalf("unexpected Forks")
	}
	if v := info.OpenIssues(); v == nil || *v != 4 {
		t.Fatalf("unexpected OpenIssues")
	}
	if v := info.LastCommitDate(); v == nil || *v != "2024-01-01T00:00:00Z" {
		t.Fatalf("unexpected LastCommitDate")
	}
	if v := info.GithubRepoURL(); v == nil || *v != "https://github.com/x/y" {
		t.Fatalf("unexpected GithubRepoURL")
	}
	if v := info.Archived(); v == nil || *v != true {
		t.Fatalf("unexpected Archived")
	}
	if v := info.Score(); v == nil || *v != 42 {
		t.Fatalf("unexpected Score")
	}
	if v := info.Skip(); v == nil || *v != false {
		t.Fatalf("unexpected Skip false")
	}
	if info.SkipReason() != nil {
		t.Fatalf("expected nil SkipReason")
	}
}

func TestAnalyzedLibInfo_SkipReason_FromLibInfo(t *testing.T) {
	t.Parallel()

	lib := parser.LibInfo{Name: "lib3", Skip: true, SkipReason: "reason"}
	repo := analyzer.GitHubRepoInfo{Score: 1}
	info := AnalyzedLibInfo{LibInfo: &lib, GitHubRepoInfo: &repo}

	if v := info.Skip(); v == nil || *v != true {
		t.Fatalf("expected Skip true from LibInfo")
	}
	if v := info.SkipReason(); v == nil || *v != "reason" {
		t.Fatalf("expected reason")
	}
}

func TestAnalyzedLibInfo_SkipReason_FromRepoInfo(t *testing.T) {
	t.Parallel()

	lib := parser.LibInfo{Name: "lib"}
	repo := analyzer.GitHubRepoInfo{Skip: true, SkipReason: "repo-reason"}
	info := AnalyzedLibInfo{LibInfo: &lib, GitHubRepoInfo: &repo}

	if v := info.Skip(); v == nil || *v != true {
		t.Fatalf("expected Skip true from RepoInfo")
	}
	if v := info.SkipReason(); v == nil || *v != "repo-reason" {
		t.Fatalf("expected repo-reason")
	}
}
