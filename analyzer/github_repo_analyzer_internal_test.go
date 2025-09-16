package analyzer

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestParseRepoURL_Variants(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in    string
		owner string
		repo  string
	}{
		{"https://github.com/user/repo", "user", "repo"},
		{"https://github.com/user/repo/", "user", "repo"},
		{"https://github.com/user/repo.git", "user", "repo"},
		{"https://github.com/user/repo/tree/main", "user", "repo"},
		{"http://github.com/user/repo/tree/main/subdir", "user", "repo"},
	}

	for _, tc := range cases {
		o, r := parseRepoURL(tc.in)
		if o != tc.owner || r != tc.repo {
			t.Fatalf("parseRepoURL(%q) => %s/%s, want %s/%s", tc.in, o, r, tc.owner, tc.repo)
		}
	}
}

func TestCalcScore_InvalidDate_SetsSkip(t *testing.T) {
	t.Parallel()

	info := &GitHubRepoInfo{LastCommitDate: "invalid-date"}
	w := &ParameterWeights{}

	calcScore(info, w)

	if !info.Skip {
		t.Fatalf("expected Skip=true when date invalid")
	}
	if info.SkipReason == "" {
		t.Fatalf("expected SkipReason to be set")
	}
}

func TestCreateRepoInfo_MapsFields(t *testing.T) {
	t.Parallel()
	rd := &RepoData{Name: "r", SubscribersCount: 1, StargazersCount: 2, ForksCount: 3, OpenIssuesCount: 4, Archived: true}
    gi := createRepoInfo(rd, "2024-01-01T00:00:00Z")
    //nolint:lll // consolidate field checks for clarity
    if gi.RepositoryName != "r" || gi.Watchers != 1 || gi.Stars != 2 || gi.Forks != 3 || gi.OpenIssues != 4 || gi.Archived != true || gi.LastCommitDate != "2024-01-01T00:00:00Z" {
		t.Fatalf("unexpected mapping: %+v", gi)
	}
}

func TestIndexOf(t *testing.T) {
	t.Parallel()
	s := []string{"a", "b", "c"}
	if indexOf(s, "b") != 1 {
		t.Fatalf("want 1")
	}
	if indexOf(s, "x") != -1 {
		t.Fatalf("want -1")
	}
}

type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestFetchJSONData_Non200AndDecodeError(t *testing.T) {
	t.Parallel()

	var out interface{}

	// Non-200 client
	c1 := &http.Client{Transport: rtFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusTeapot, Body: io.NopCloser(strings.NewReader("teapot")), Header: make(http.Header)}, nil
	})}
	err := fetchJSONData(c1, "http://example", nil, &out)
	if !errors.Is(err, ErrUnexpectedStatusCode) {
		t.Fatalf("expected ErrUnexpectedStatusCode, got %v", err)
	}

	// 200 but invalid JSON
	c2 := &http.Client{Transport: rtFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("not-json")), Header: make(http.Header)}, nil
	})}
	err = fetchJSONData(c2, "http://example", nil, &out)
	if err == nil {
		t.Fatalf("expected decode error")
	}
}

func TestFetchGithubInfo_NoToken_SetsSkip(t *testing.T) {
	t.Parallel()
	a := NewGitHubRepoAnalyzer("", NewParameterWeights())
	infos := a.FetchGithubInfo([]string{"https://github.com/user/repo"})
	if len(infos) != 1 {
		t.Fatalf("want 1 info")
	}
	if !infos[0].Skip {
		t.Fatalf("expected Skip true when token missing")
	}
	if infos[0].SkipReason == "" {
		t.Fatalf("expected SkipReason set")
	}
}
