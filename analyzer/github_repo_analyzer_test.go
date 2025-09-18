package analyzer_test

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/uzumaki-inc/stay_or_go/analyzer"
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
	analyzer := analyzer.NewGitHubRepoAnalyzer("dummy-token", analyzer.ParameterWeights{
		Forks:          1.0,
		OpenIssues:     1.0,
		LastCommitDate: 1.0,
		Archived:       1.0,
	})

	// テスト実行
	repoURLs := []string{"https://github.com/example-owner/example-repo"}
	repoInfos := analyzer.FetchGithubInfo(repoURLs)

	assert.Len(t, repoInfos, 1, "Expected 1 repo info")

	repoInfo := repoInfos[0]

	assert.Equal(t, "example-repo", repoInfo.RepositoryName, "RepositoryName mismatch")
	assert.Equal(t, 10, repoInfo.Watchers, "Watchers mismatch")
	assert.Equal(t, 50, repoInfo.Stars, "Stars mismatch")
	assert.Equal(t, 5, repoInfo.Forks, "Forks mismatch")
	assert.Equal(t, 3, repoInfo.OpenIssues, "OpenIssues mismatch")
	assert.False(t, repoInfo.Archived, "Archived should be false")
	assert.False(t, repoInfo.Skip, "Skip should be false")
}

func TestScoreIncludesWatchers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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

	analyzer := analyzer.NewGitHubRepoAnalyzer("dummy-token", analyzer.ParameterWeights{
		Watchers:       2.0,
		Stars:          0.0,
		Forks:          0.0,
		OpenIssues:     0.0,
		LastCommitDate: 0.0,
		Archived:       0.0,
	})

	repoInfos := analyzer.FetchGithubInfo([]string{"https://github.com/example-owner/example-repo"})

	assert.Len(t, repoInfos, 1)
	assert.Equal(t, 20, repoInfos[0].Score)
}

func TestFetchGithubInfo_ConcurrentProcessingMaintainsOrder(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// 複数のリポジトリをセットアップ（順序を確認するため）
	repoURLs := []string{
		"https://github.com/owner/repo1",
		"https://github.com/owner/repo2",
		"https://github.com/owner/repo3",
		"https://github.com/owner/repo4",
		"https://github.com/owner/repo5",
	}

	// 各リポジトリに対してモックレスポンスを設定
	for idx := range repoURLs {
		repoName := fmt.Sprintf("repo%d", idx+1)
		httpmock.RegisterResponder("GET", "https://api.github.com/repos/owner/"+repoName,
			httpmock.NewStringResponder(200, fmt.Sprintf(`{
				"name": "%s",
				"subscribers_count": %d,
				"stargazers_count": %d,
				"forks_count": %d,
				"open_issues_count": %d,
				"archived": false,
				"default_branch": "main"
			}`, repoName, idx+1, (idx+1)*10, (idx+1)*2, (idx+1)*3)))

		httpmock.RegisterResponder("GET", "https://api.github.com/repos/owner/"+repoName+"/pulls",
			httpmock.NewStringResponder(200, `[]`)) // 空のプルリクエストリスト

		httpmock.RegisterResponder("GET", "https://api.github.com/repos/owner/"+repoName+"/commits/main",
			httpmock.NewStringResponder(200, `{
				"commit": {
					"committer": {
						"date": "2023-10-01T12:00:00Z"
					}
				}
			}`))
	}

	analyzer := analyzer.NewGitHubRepoAnalyzer("dummy-token", analyzer.ParameterWeights{
		Watchers:       0.1,
		Stars:          0.1,
		Forks:          0.1,
		OpenIssues:     0.01,
		LastCommitDate: -0.05,
		Archived:       -1000000,
	})

	// 並列処理でも順序が維持されることをテスト
	repoInfos := analyzer.FetchGithubInfo(repoURLs)

	assert.Len(t, repoInfos, 5, "Expected 5 repo infos")

	// 順序が入力と同じであることを確認
	for i, repoInfo := range repoInfos {
		expectedName := fmt.Sprintf("repo%d", i+1)
		assert.Equal(t, expectedName, repoInfo.RepositoryName,
			"Repository order not maintained at index %d", i)
		assert.Equal(t, repoURLs[i], repoInfo.GithubRepoURL,
			"URL order not maintained at index %d", i)
	}
}
