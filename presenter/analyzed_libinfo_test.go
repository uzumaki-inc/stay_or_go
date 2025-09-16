package presenter

import (
    "testing"

    "github.com/uzumaki-inc/stay_or_go/analyzer"
    "github.com/uzumaki-inc/stay_or_go/parser"
)

func TestAnalyzedLibInfo_Getters_NilAndValues(t *testing.T) {
    t.Parallel()

    // Case 1: Only LibInfo with Skip (GitHubRepoInfo is nil)
    a1 := AnalyzedLibInfo{LibInfo: &parser.LibInfo{Skip: true, SkipReason: "li"}, GitHubRepoInfo: nil}
    if a1.Name() != nil { t.Fatalf("expected nil Name") }
    if a1.RepositoryURL() != nil { t.Fatalf("expected nil RepositoryURL") }
    if a1.Watchers() != nil { t.Fatalf("expected nil Watchers") }
    if a1.Stars() != nil { t.Fatalf("expected nil Stars") }
    if a1.Forks() != nil { t.Fatalf("expected nil Forks") }
    if a1.OpenIssues() != nil { t.Fatalf("expected nil OpenIssues") }
    if a1.LastCommitDate() != nil { t.Fatalf("expected nil LastCommitDate") }
    if a1.GithubRepoURL() != nil { t.Fatalf("expected nil GithubRepoURL") }
    if a1.Archived() != nil { t.Fatalf("expected nil Archived") }
    if a1.Score() != nil { t.Fatalf("expected nil Score") }
    if a1.Skip() == nil || *a1.Skip() != true { t.Fatalf("expected Skip=true from LibInfo") }
    if v := a1.SkipReason(); v == nil || *v != "li" { t.Fatalf("expected SkipReason from LibInfo") }

    // Case 2: LibInfo with data, and RepoInfo with data
    li := parser.LibInfo{Name: "lib", RepositoryURL: "https://github.com/x/y"}
    ri := analyzer.GitHubRepoInfo{Watchers: 1, Stars: 2, Forks: 3, OpenIssues: 4, LastCommitDate: "2024-01-01T00:00:00Z", GithubRepoURL: "https://github.com/x/y", Archived: true, Score: 42}
    a2 := AnalyzedLibInfo{LibInfo: &li, GitHubRepoInfo: &ri}
    if v := a2.Name(); v == nil || *v != "lib" { t.Fatalf("unexpected Name") }
    if v := a2.RepositoryURL(); v == nil || *v != "https://github.com/x/y" { t.Fatalf("unexpected RepositoryURL") }
    if v := a2.Watchers(); v == nil || *v != 1 { t.Fatalf("unexpected Watchers") }
    if v := a2.Stars(); v == nil || *v != 2 { t.Fatalf("unexpected Stars") }
    if v := a2.Forks(); v == nil || *v != 3 { t.Fatalf("unexpected Forks") }
    if v := a2.OpenIssues(); v == nil || *v != 4 { t.Fatalf("unexpected OpenIssues") }
    if v := a2.LastCommitDate(); v == nil || *v != "2024-01-01T00:00:00Z" { t.Fatalf("unexpected LastCommitDate") }
    if v := a2.GithubRepoURL(); v == nil || *v != "https://github.com/x/y" { t.Fatalf("unexpected GithubRepoURL") }
    if v := a2.Archived(); v == nil || *v != true { t.Fatalf("unexpected Archived") }
    if v := a2.Score(); v == nil || *v != 42 { t.Fatalf("unexpected Score") }
    if v := a2.Skip(); v == nil || *v != false { t.Fatalf("unexpected Skip false") }
    if a2.SkipReason() != nil { t.Fatalf("expected nil SkipReason") }

    // Case 3: Skip due to LibInfo
    li3 := parser.LibInfo{Name: "lib3", Skip: true, SkipReason: "reason"}
    a3 := AnalyzedLibInfo{LibInfo: &li3, GitHubRepoInfo: &ri}
    if v := a3.Skip(); v == nil || *v != true { t.Fatalf("expected Skip true from LibInfo") }
    if v := a3.SkipReason(); v == nil || *v != "reason" { t.Fatalf("expected reason") }

    // Case 4: Skip due to RepoInfo
    ri4 := analyzer.GitHubRepoInfo{Skip: true, SkipReason: "repo-reason"}
    a4 := AnalyzedLibInfo{LibInfo: &li, GitHubRepoInfo: &ri4}
    if v := a4.Skip(); v == nil || *v != true { t.Fatalf("expected Skip true from RepoInfo") }
    if v := a4.SkipReason(); v == nil || *v != "repo-reason" { t.Fatalf("expected repo-reason") }
}
