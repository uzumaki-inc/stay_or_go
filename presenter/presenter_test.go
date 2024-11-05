package presenter

import (
	"testing"

	"github.com/konyu/StayOrGo/analyzer"
	"github.com/konyu/StayOrGo/parser"
	"github.com/stretchr/testify/assert"
)

func TestMakeBody(t *testing.T) {
	libInfo1 := parser.LibInfo{Name: "lib1", RepositoryUrl: "https://github.com/lib1"}
	repoInfo1 := analyzer.GitHubRepoInfo{RepositoryName: "lib1", Watchers: 100, Stars: 200, Forks: 50, OpenPullRequests: 5, OpenIssues: 10, LastCommitDate: "2023-10-10", Archived: false, Score: 85}
	libInfo2 := parser.LibInfo{Name: "lib2", RepositoryUrl: "https://github.com/lib2"}
	repoInfo2 := analyzer.GitHubRepoInfo{RepositoryName: "lib2", Watchers: 150, Stars: 250, Forks: 60, OpenPullRequests: 7, OpenIssues: 15, LastCommitDate: "2023-10-11", Archived: false, Score: 90}

	analyzedLibInfos := []AnalyzedLibInfo{
		{LibInfo: &libInfo1, GitHubRepoInfo: &repoInfo1},
		{LibInfo: &libInfo2, GitHubRepoInfo: &repoInfo2},
	}

	body := makeBody(analyzedLibInfos, "|")

	assert.Len(t, body, 2)
	assert.Equal(t, "|lib1|https://github.com/lib1|100|200|50|5|10|2023-10-10|false|85|false|N/A|", body[0])
	assert.Equal(t, "|lib2|https://github.com/lib2|150|250|60|7|15|2023-10-11|false|90|false|N/A|", body[1])
}
