//nolint:testpackage // Tests unexported methods
package presenter

import (
	"strings"
	"testing"
)

func TestMarkdownPresenter_makeHeader(t *testing.T) {
	t.Parallel()

	p := NewMarkdownPresenter(nil)
	header := p.makeHeader()

	if len(header) != 2 {
		t.Fatalf("expected 2 header lines, got %d", len(header))
	}

	if !strings.HasPrefix(header[0], "| ") || !strings.HasSuffix(header[0], " |") {
		t.Fatalf("unexpected header row: %q", header[0])
	}

	if !strings.HasPrefix(header[1], "|") || !strings.HasSuffix(header[1], "|") {
		t.Fatalf("unexpected separator row: %q", header[1])
	}
}

func TestCsvPresenter_makeHeader(t *testing.T) {
	t.Parallel()

	p := NewCsvPresenter(nil)
	header := p.makeHeader()

	if len(header) != 1 {
		t.Fatalf("expected single header row")
	}

	if !strings.Contains(header[0], ", ") {
		t.Fatalf("expected commas in header: %q", header[0])
	}
}

func TestTsvPresenter_makeHeader(t *testing.T) {
	t.Parallel()

	p := NewTsvPresenter(nil)
	header := p.makeHeader()

	if len(header) != 1 {
		t.Fatalf("expected single header row")
	}

	if !strings.Contains(header[0], "\t") {
		t.Fatalf("expected tabs in header: %q", header[0])
	}
}
