package presenter_test

import (
	"testing"

	"github.com/uzumaki-inc/stay_or_go/presenter"
)

func TestSelectPresenter_Formats(t *testing.T) {
	t.Parallel()

	infos := []presenter.AnalyzedLibInfo{}

	p := presenter.SelectPresenter("csv", infos)
	if _, ok := p.(presenter.CsvPresenter); !ok {
		t.Fatalf("expected CsvPresenter, got %T", p)
	}

	p = presenter.SelectPresenter("tsv", infos)
	if _, ok := p.(presenter.TsvPresenter); !ok {
		t.Fatalf("expected TsvPresenter, got %T", p)
	}

	// default → markdown
	p = presenter.SelectPresenter("unknown", infos)
	if _, ok := p.(presenter.MarkdownPresenter); !ok {
		t.Fatalf("expected MarkdownPresenter, got %T", p)
	}

	p = presenter.SelectPresenter("markdown", infos)
	if _, ok := p.(presenter.MarkdownPresenter); !ok {
		t.Fatalf("expected MarkdownPresenter, got %T", p)
	}
}
