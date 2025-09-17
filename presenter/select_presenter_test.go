package presenter_test

import (
	"testing"

	"github.com/uzumaki-inc/stay_or_go/presenter"
)

func TestSelectPresenter_Formats(t *testing.T) {
	t.Parallel()

	infos := []presenter.AnalyzedLibInfo{}

	pres := presenter.SelectPresenter("csv", infos)
	if _, ok := pres.(presenter.CsvPresenter); !ok {
		t.Fatalf("expected CsvPresenter, got %T", pres)
	}

	pres = presenter.SelectPresenter("tsv", infos)
	if _, ok := pres.(presenter.TsvPresenter); !ok {
		t.Fatalf("expected TsvPresenter, got %T", pres)
	}

	// default → markdown
	pres = presenter.SelectPresenter("unknown", infos)
	if _, ok := pres.(presenter.MarkdownPresenter); !ok {
		t.Fatalf("expected MarkdownPresenter, got %T", pres)
	}

	pres = presenter.SelectPresenter("markdown", infos)
	if _, ok := pres.(presenter.MarkdownPresenter); !ok {
		t.Fatalf("expected MarkdownPresenter, got %T", pres)
	}
}
