package analyzer

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchGithubInfo(t *testing.T) {
	// httpmockを有効化
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// モックレスポンスを設定
	httpmock.RegisterResponder("GET", "https://api.github.com/repos/example-owner/example-repo",
		httpmock.NewStringResponder(200, `{
			"name": "example-repo",
			"subscribers_count": 10,
			"stargazers_count": 50,
			"forks_count": 5,
			"open_issues_count": 3,
			"archived": false,
			"default_branch": "main"
		}`))

	httpmock.RegisterResponder("GET", "https://api.github.com/repos/example-owner/example-repo/pulls",
		httpmock.NewStringResponder(200, `[]`)) // 空のプルリクエストリスト

	httpmock.RegisterResponder("GET", "https://api.github.com/repos/example-owner/example-repo/commits/main",
		httpmock.NewStringResponder(200, `{
			"commit": {
				"committer": {
					"date": "2023-10-01T12:00:00Z"
				}
			}
		}`))

	// テスト用のGitHubRepoAnalyzerを作成
	analyzer := NewGitHubRepoAnalyzer("dummy-token", ParameterWeights{
		Forks:            1.0,
		OpenPullRequests: 1.0,
		OpenIssues:       1.0,
		LastCommitDate:   1.0,
		Archived:         1.0,
	})

	// テスト実行
	repoUrls := []string{"https://api.github.com/repos/example-owner/example-repo"}
	repoInfos := analyzer.FetchGithubInfo(repoUrls)

	assert.Equal(t, 1, len(repoInfos), "Expected 1 repo info")

	repoInfo := repoInfos[0]

	assert.Equal(t, "example-repo", repoInfo.RepositoryName, "RepositoryName mismatch")
	assert.Equal(t, 10, repoInfo.Watchers, "Watchers mismatch")
	assert.Equal(t, 50, repoInfo.Stars, "Stars mismatch")
	assert.Equal(t, 5, repoInfo.Forks, "Forks mismatch")
	assert.Equal(t, 0, repoInfo.OpenPullRequests, "OpenPullRequests mismatch")
	assert.Equal(t, 3, repoInfo.OpenIssues, "OpenIssues mismatch")
	assert.False(t, repoInfo.Archived, "Archived should be false")
	assert.False(t, repoInfo.Skip, "Skip should be false")
}
