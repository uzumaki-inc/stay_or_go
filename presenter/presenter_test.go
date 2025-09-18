package presenter_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uzumaki-inc/stay_or_go/analyzer"
	"github.com/uzumaki-inc/stay_or_go/parser"
	"github.com/uzumaki-inc/stay_or_go/presenter"
)

// Disable parallel testing to test standard output
//
//nolint:paralleltest,funlen // Test manipulates os.Stdout, complex test cases
func TestDisplay(t *testing.T) {
	// Avoid running in parallel since this test manipulates os.Stdout
	testCases := []struct {
		name           string
		presenterFunc  func([]presenter.AnalyzedLibInfo) presenter.Presenter
		expectedOutput string
	}{
		{
			name: "MarkDown Presenter",
			presenterFunc: func(analyzedLibInfos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewMarkdownPresenter(analyzedLibInfos)
			},

			//nolint:lll
			expectedOutput: "| Name | RepositoryURL | Watchers | Stars | Forks | OpenIssues | LastCommitDate | Archived | Score | Skip | SkipReason |\n" +
				"| ---- | ------------- | -------- | ----- | ----- | ---------- | -------------- | -------- | ----- | ---- | ---------- |\n" +
				"|lib1|https://github.com/lib1|100|200|50|10|2023-10-10|false|85|false|N/A|\n" +
				"|lib2|https://github.com/lib2|150|250|60|15|2023-10-11|false|90|false|N/A|\n",
		},
		{
			name: "TSV Presenter",
			presenterFunc: func(analyzedLibInfos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewTsvPresenter(analyzedLibInfos)
			},
			//nolint:lll
			expectedOutput: "Name\tRepositoryURL\tWatchers\tStars\tForks\tOpenIssues\tLastCommitDate\tArchived\tScore\tSkip\tSkipReason\n" +
				"lib1\thttps://github.com/lib1\t100\t200\t50\t10\t2023-10-10\tfalse\t85\tfalse\tN/A\n" +
				"lib2\thttps://github.com/lib2\t150\t250\t60\t15\t2023-10-11\tfalse\t90\tfalse\tN/A\n",
		},
		{
			name: "CSV Presenter",
			presenterFunc: func(analyzedLibInfos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewCsvPresenter(analyzedLibInfos)
			},
			//nolint:lll
			expectedOutput: "Name, RepositoryURL, Watchers, Stars, Forks, OpenIssues, LastCommitDate, Archived, Score, Skip, SkipReason\n" +
				"lib1, https://github.com/lib1, 100, 200, 50, 10, 2023-10-10, false, 85, false, N/A\n" +
				"lib2, https://github.com/lib2, 150, 250, 60, 15, 2023-10-11, false, 90, false, N/A\n",
		},
	}

	//nolint:paralleltest
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			libInfo1 := parser.LibInfo{Name: "lib1", RepositoryURL: "https://github.com/lib1"}
			repoInfo1 := analyzer.GitHubRepoInfo{
				RepositoryName: "lib1", Watchers: 100, Stars: 200, Forks: 50,
				OpenIssues: 10, LastCommitDate: "2023-10-10", Archived: false, Score: 85,
			}
			libInfo2 := parser.LibInfo{Name: "lib2", RepositoryURL: "https://github.com/lib2"}
			repoInfo2 := analyzer.GitHubRepoInfo{
				RepositoryName: "lib2", Watchers: 150, Stars: 250, Forks: 60,
				OpenIssues: 15, LastCommitDate: "2023-10-11", Archived: false, Score: 90,
			}

			analyzedLibInfos := []presenter.AnalyzedLibInfo{
				{LibInfo: &libInfo1, GitHubRepoInfo: &repoInfo1},
				{LibInfo: &libInfo2, GitHubRepoInfo: &repoInfo2},
			}

			presenter := testCase.presenterFunc(analyzedLibInfos)

			// 標準出力をキャプチャするためのバッファを作成
			readPipe, writePipe, _ := os.Pipe()
			originalStdout := os.Stdout

			defer func() { os.Stdout = originalStdout }() // テスト後に元に戻す

			os.Stdout = writePipe

			// Displayメソッドを呼び出す
			presenter.Display()

			// 書き込みを閉じてから、キャプチャした出力を取得
			writePipe.Close()

			var buf bytes.Buffer

			_, err := buf.ReadFrom(readPipe)
			if err != nil {
				t.Fatalf("failed to read from pipe: %v", err)
			}

			output := buf.String()

			// 期待される出力を検証
			assert.Equal(t, testCase.expectedOutput, output)
		})
	}
}

// TestMakeAnalyzedLibInfoList_PointerBug tests that each AnalyzedLibInfo
// has its own LibInfo instance, not sharing the same pointer
func TestMakeAnalyzedLibInfoList_PointerBug(t *testing.T) {
	t.Parallel()

	// Create test data
	libInfoList := []parser.LibInfo{
		{Name: "lib1", RepositoryURL: "https://github.com/owner/lib1", Others: []string{"v1.0.0"}},
		{Name: "lib2", RepositoryURL: "https://github.com/owner/lib2", Others: []string{"v2.0.0"}},
		{Name: "lib3", RepositoryURL: "https://github.com/owner/lib3", Others: []string{"v3.0.0"}},
	}

	gitHubRepoInfos := []analyzer.GitHubRepoInfo{
		{
			GithubRepoURL:  "https://github.com/owner/lib1",
			RepositoryName: "lib1",
			Stars:          100,
		},
		{
			GithubRepoURL:  "https://github.com/owner/lib2",
			RepositoryName: "lib2",
			Stars:          200,
		},
	}

	// Call the function under test
	result := presenter.MakeAnalyzedLibInfoList(libInfoList, gitHubRepoInfos)

	// Verify that we have the expected number of results
	assert.Len(t, result, 3)

	// Verify that each AnalyzedLibInfo has a unique LibInfo pointer
	// and correct values
	assert.Equal(t, "lib1", result[0].LibInfo.Name)
	assert.Equal(t, []string{"v1.0.0"}, result[0].LibInfo.Others)
	assert.NotNil(t, result[0].GitHubRepoInfo)
	assert.Equal(t, 100, result[0].GitHubRepoInfo.Stars)

	assert.Equal(t, "lib2", result[1].LibInfo.Name)
	assert.Equal(t, []string{"v2.0.0"}, result[1].LibInfo.Others)
	assert.NotNil(t, result[1].GitHubRepoInfo)
	assert.Equal(t, 200, result[1].GitHubRepoInfo.Stars)

	assert.Equal(t, "lib3", result[2].LibInfo.Name)
	assert.Equal(t, []string{"v3.0.0"}, result[2].LibInfo.Others)
	assert.Nil(t, result[2].GitHubRepoInfo)

	// Most important: verify that modifying one LibInfo doesn't affect others
	// This would fail with the pointer bug
	result[0].LibInfo.Name = "modified"
	assert.Equal(t, "modified", result[0].LibInfo.Name)
	assert.Equal(t, "lib2", result[1].LibInfo.Name)
	assert.Equal(t, "lib3", result[2].LibInfo.Name)
}
