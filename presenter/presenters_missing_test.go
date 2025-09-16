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
//nolint:tparallel,paralleltest,lll
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
            //nolint:dupword
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
            //nolint:dupword
            expectedOutput: "Name, RepositoryURL, Watchers, Stars, Forks, OpenIssues, LastCommitDate, Archived, Score, Skip, SkipReason\nlibX, N/A, N/A, N/A, N/A, N/A, N/A, N/A, N/A, true, Not hosted on Github\n",
		},
		{
			name: "TSV",
			presenterFunc: func(infos []presenter.AnalyzedLibInfo) presenter.Presenter {
				return presenter.NewTsvPresenter(infos)
			},
			expectedOutput: "Name\tRepositoryURL\tWatchers\tStars\tForks\tOpenIssues\tLastCommitDate\tArchived\tScore\tSkip\tSkipReason\nlibX\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\tN/A\ttrue\tNot hosted on Github\n",
		},
	}

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
			// capture stdout
			r, w, _ := os.Pipe()
			old := os.Stdout
			os.Stdout = w
			defer func() { os.Stdout = old }()

			p := tc.presenterFunc(analyzed)
			p.Display()

			w.Close()
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)

			assert.Equal(t, tc.expectedOutput, buf.String())
		})
	}
}
