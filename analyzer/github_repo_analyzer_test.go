package analyzer

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGitHubRepoAnalyzer_FetchGithubInfo(t *testing.T) {
	// Mock HTTP requests for GitHub API
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Setup mock responses for repository details
	httpmock.RegisterResponder("GET", "https://api.github.com/repos/rails/rails",
		httpmock.NewStringResponder(200, `{
			"name": "rails",
			"watchers_count": 500,
			"stargazers_count": 1000,
			"forks_count": 200,
			"open_issues_count": 20,
			"default_branch": "main",
			"archived": false
		}`))

	httpmock.RegisterResponder("GET", "https://api.github.com/repos/rails/rails/pulls",
		httpmock.NewStringResponder(200, `[{}, {}, {}]`))

	httpmock.RegisterResponder("GET", "https://api.github.com/repos/rails/rails/commits/main",
		httpmock.NewStringResponder(200, `{
			"commit": {
				"committer": {
					"date": "2023-10-10T12:00:00Z"
				}
			}
		}`))

	// Set up analyzer
	weights := ParameterWeights{
		Score:            1.0,
		Forks:            0.5,
		OpenPullRequests: 0.3,
		OpenIssues:       -0.2,
		LastCommitDate:   -0.1,
		Archived:         -1.0,
	}
	analyzer := NewGitHubRepoAnalyzer("mock-token", weights)

	// Run FetchGithubInfo method
	repoUrls := []string{"https://github.com/rails/rails"}
	libraryInfoList := analyzer.FetchGithubInfo(repoUrls)

	// Assertions
	assert.Len(t, libraryInfoList, 1)
	assert.Equal(t, "rails", libraryInfoList[0].RepositoryName)
	assert.Equal(t, 500, libraryInfoList[0].Watchers)
	assert.Equal(t, 1000, libraryInfoList[0].Stars)
	assert.Equal(t, 200, libraryInfoList[0].Forks)
	assert.Equal(t, 3, libraryInfoList[0].OpenPullRequests)
	assert.Equal(t, 20, libraryInfoList[0].OpenIssues)
	assert.Equal(t, "2023-10-10T12:00:00Z", libraryInfoList[0].LastCommitDate)
	assert.False(t, libraryInfoList[0].Archived)
	assert.Equal(t, "https://github.com/rails/rails", libraryInfoList[0].GithubRepoUrl)
}

func TestGitHubRepoAnalyzer_getGitHubInfo_ErrorHandling(t *testing.T) {
	// Mock HTTP requests for GitHub API with error responses
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.github.com/repos/unknown/repo",
		httpmock.NewStringResponder(404, `{"message": "Not Found"}`))

	// Setup mock response for pulls and commits (which might be called internally)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"message": "Not Found"}`))

	// Set up analyzer
	weights := ParameterWeights{}
	analyzer := NewGitHubRepoAnalyzer("mock-token", weights)

	// Run FetchGithubInfo method
	repoUrls := []string{"https://github.com/unknown/repo"}
	libraryInfoList := analyzer.FetchGithubInfo(repoUrls)

	// Assertions
	assert.Len(t, libraryInfoList, 1)
	assert.NotNil(t, libraryInfoList[0])
	assert.True(t, libraryInfoList[0].Skip)
	assert.Equal(t, "Failed fetching https://github.com/unknown/repo from GitHub", libraryInfoList[0].SkipReason)
}
