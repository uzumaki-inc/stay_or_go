package presenter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uzumaki-inc/stay_or_go/analyzer"
	"github.com/uzumaki-inc/stay_or_go/parser"
	"github.com/uzumaki-inc/stay_or_go/presenter"
)

func TestAnalyzedLibInfo_EmptyLibInfo_NoRepoInfo(t *testing.T) {
	t.Parallel()

	info := presenter.AnalyzedLibInfo{LibInfo: &parser.LibInfo{Skip: true, SkipReason: "li"}, GitHubRepoInfo: nil}

	assert.Nil(t, info.Name())
	assert.Nil(t, info.RepositoryURL())
	assert.Nil(t, info.Watchers())
	assert.Nil(t, info.Stars())
	assert.Nil(t, info.Forks())
	assert.Nil(t, info.OpenIssues())
	assert.Nil(t, info.LastCommitDate())
	assert.Nil(t, info.GithubRepoURL())
	assert.Nil(t, info.Archived())
	assert.Nil(t, info.Score())

	if info.Skip() == nil {
		t.Fatalf("skip pointer is nil")
	}

	assert.True(t, *info.Skip())

	if v := info.SkipReason(); v == nil {
		t.Fatalf("skip reason nil")
	} else {
		assert.Equal(t, "li", *v)
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
	info := presenter.AnalyzedLibInfo{LibInfo: &lib, GitHubRepoInfo: &repo}

	assert.NotNil(t, info.Name())
	assert.Equal(t, "lib", *info.Name())
	assert.NotNil(t, info.RepositoryURL())
	assert.Equal(t, "https://github.com/x/y", *info.RepositoryURL())
	assert.NotNil(t, info.Watchers())
	assert.Equal(t, 1, *info.Watchers())
	assert.NotNil(t, info.Stars())
	assert.Equal(t, 2, *info.Stars())
	assert.NotNil(t, info.Forks())
	assert.Equal(t, 3, *info.Forks())
	assert.NotNil(t, info.OpenIssues())
	assert.Equal(t, 4, *info.OpenIssues())
	assert.NotNil(t, info.LastCommitDate())
	assert.Equal(t, "2024-01-01T00:00:00Z", *info.LastCommitDate())
	assert.NotNil(t, info.GithubRepoURL())
	assert.Equal(t, "https://github.com/x/y", *info.GithubRepoURL())
	assert.NotNil(t, info.Archived())
	assert.True(t, *info.Archived())
	assert.NotNil(t, info.Score())
	assert.Equal(t, 42, *info.Score())
	assert.NotNil(t, info.Skip())
	assert.False(t, *info.Skip())
	assert.Nil(t, info.SkipReason())
}

func TestAnalyzedLibInfo_SkipReason_FromLibInfo(t *testing.T) {
	t.Parallel()

	lib := parser.LibInfo{Name: "lib3", Skip: true, SkipReason: "reason"}
	repo := analyzer.GitHubRepoInfo{Score: 1}
	info := presenter.AnalyzedLibInfo{LibInfo: &lib, GitHubRepoInfo: &repo}

	if v := info.Skip(); v == nil {
		t.Fatalf("skip nil")
	} else {
		assert.True(t, *v)
	}

	if v := info.SkipReason(); v == nil {
		t.Fatalf("reason nil")
	} else {
		assert.Equal(t, "reason", *v)
	}
}

func TestAnalyzedLibInfo_SkipReason_FromRepoInfo(t *testing.T) {
	t.Parallel()

	lib := parser.LibInfo{Name: "lib"}
	repo := analyzer.GitHubRepoInfo{Skip: true, SkipReason: "repo-reason"}
	info := presenter.AnalyzedLibInfo{LibInfo: &lib, GitHubRepoInfo: &repo}

	if v := info.Skip(); v == nil {
		t.Fatalf("skip nil")
	} else {
		assert.True(t, *v)
	}

	if v := info.SkipReason(); v == nil {
		t.Fatalf("reason nil")
	} else {
		assert.Equal(t, "repo-reason", *v)
	}
}
