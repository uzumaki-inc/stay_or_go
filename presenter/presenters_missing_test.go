package presenter_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uzumaki-inc/stay_or_go/parser"
	"github.com/uzumaki-inc/stay_or_go/presenter"
)

// Verify presenters when GitHubRepoInfo is missing and LibInfo is skipped.
func getTestCases() []struct {
	name           string
	presenterFunc  func([]presenter.AnalyzedLibInfo) presenter.Presenter
	expectedOutput string
} {
	return []struct {
		name           string
		presenterFunc  func([]presenter.AnalyzedLibInfo) presenter.Presenter
		expectedOutput string
	}{
		{
			name: "Markdown",
			presenterFunc: func(infos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewMarkdownPresenter(infos)
			},
			expectedOutput: `| Name | RepositoryURL | Watchers | Stars | Forks | OpenIssues | ` +
				`LastCommitDate | Archived | Score | Skip | SkipReason |
| ---- | ------------- | -------- | ----- | ----- | ---------- | ` +
				`-------------- | -------- | ----- | ---- | ---------- |
|libX|N/A|N/A|N/A|N/A|N/A|N/A|N/A|N/A|true|Not hosted on Github|
`,
		},
		{
			name: "CSV",
			presenterFunc: func(infos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewCsvPresenter(infos)
			},
			//nolint:dupword // N/A repetition is expected output format
			expectedOutput: "Name, RepositoryURL, Watchers, Stars, Forks, OpenIssues, " +
				"LastCommitDate, Archived, Score, Skip, SkipReason\n" +
				"libX, N/A, N/A, N/A, N/A, N/A, N/A, N/A, N/A, true, Not hosted on Github\n",
		},
		{
			name: "TSV",
			presenterFunc: func(infos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewTsvPresenter(infos)
			},
			//nolint:dupword // N/A repetition is expected output format
			expectedOutput: "Name\tRepositoryURL\tWatchers\tStars\tForks\tOpenIssues\t" +
				"LastCommitDate\tArchived\tScore\tSkip\tSkipReason\n" +
				"libX\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\ttrue\tNot hosted on Github\n",
		},
	}
}

func capturePresenterOutput(pres presenter.Presenter) string {
	rPipe, wPipe, _ := os.Pipe()
	old := os.Stdout

	os.Stdout = wPipe

	defer func() { os.Stdout = old }()

	pres.Display()
	wPipe.Close()

	var buf bytes.Buffer

	_, _ = buf.ReadFrom(rPipe)

	return buf.String()
}

//nolint:paralleltest // Test manipulates os.Stdout
func TestPresenters_WithSkippedLibInfoAndMissingRepoInfo(t *testing.T) {
	lib := parser.LibInfo{
		Name:          "libX",
		RepositoryURL: "",
		Skip:          true,
		SkipReason:    "Not hosted on Github",
	}

	analyzed := []presenter.AnalyzedLibInfo{{LibInfo: &lib, GitHubRepoInfo: nil}}
	cases := getTestCases()

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			p := testCase.presenterFunc(analyzed)
			output := capturePresenterOutput(p)
			assert.Equal(t, testCase.expectedOutput, output)
		})
	}
}
