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
//
//nolint:paralleltest,lll,funlen // Test manipulates os.Stdout
func TestPresenters_WithSkippedLibInfoAndMissingRepoInfo(t *testing.T) {
	// Avoid running in parallel since this test manipulates os.Stdout
	lib := parser.LibInfo{
		Name:          "libX",
		RepositoryURL: "",
		Skip:          true,
		SkipReason:    "Not hosted on Github",
	}

	analyzed := []presenter.AnalyzedLibInfo{{LibInfo: &lib, GitHubRepoInfo: nil}}

	cases := []struct {
		name           string
		presenterFunc  func([]presenter.AnalyzedLibInfo) presenter.Presenter
		expectedOutput string
	}{
		{
			name: "Markdown",
			presenterFunc: func(infos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewMarkdownPresenter(infos)
			},
			expectedOutput: `| Name | RepositoryURL | Watchers | Stars | Forks | OpenIssues | LastCommitDate | Archived | Score | Skip | SkipReason |
| ---- | ------------- | -------- | ----- | ----- | ---------- | -------------- | -------- | ----- | ---- | ---------- |
|libX|N/A|N/A|N/A|N/A|N/A|N/A|N/A|N/A|true|Not hosted on Github|
`,
		},
		{
			name: "CSV",
			presenterFunc: func(infos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewCsvPresenter(infos)
			},
			//nolint:dupword // N/A repetition is expected output format
			expectedOutput: "Name, RepositoryURL, Watchers, Stars, Forks, OpenIssues, LastCommitDate, Archived, Score, Skip, SkipReason\nlibX, N/A, N/A, N/A, N/A, N/A, N/A, N/A, N/A, true, Not hosted on Github\n",
		},
		{
			name: "TSV",
			presenterFunc: func(infos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewTsvPresenter(infos)
			},
			//nolint:dupword // N/A repetition is expected output format
			expectedOutput: "Name\tRepositoryURL\tWatchers\tStars\tForks\tOpenIssues\tLastCommitDate\tArchived\tScore\tSkip\tSkipReason\nlibX\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\ttrue\tNot hosted on Github\n",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			// capture stdout
			rPipe, wPipe, _ := os.Pipe()
			old := os.Stdout

			os.Stdout = wPipe
			defer func() { os.Stdout = old }()

			p := testCase.presenterFunc(analyzed)
			p.Display()

			wPipe.Close()

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(rPipe)

			assert.Equal(t, testCase.expectedOutput, buf.String())
		})
	}
}
