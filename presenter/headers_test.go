package presenter

import (
    "strings"
    "testing"
)

func TestMarkdownPresenter_makeHeader(t *testing.T) {
    t.Parallel()
    p := NewMarkdownPresenter(nil)
    h := p.makeHeader()
    if len(h) != 2 { t.Fatalf("expected 2 header lines, got %d", len(h)) }
    if !strings.HasPrefix(h[0], "| ") || !strings.HasSuffix(h[0], " |") {
        t.Fatalf("unexpected header row: %q", h[0])
    }
    if !strings.HasPrefix(h[1], "|") || !strings.HasSuffix(h[1], "|") {
        t.Fatalf("unexpected separator row: %q", h[1])
    }
}

func TestCsvPresenter_makeHeader(t *testing.T) {
    t.Parallel()
    p := NewCsvPresenter(nil)
    h := p.makeHeader()
    if len(h) != 1 { t.Fatalf("expected single header row") }
    if !strings.Contains(h[0], ", ") { t.Fatalf("expected commas in header: %q", h[0]) }
}

func TestTsvPresenter_makeHeader(t *testing.T) {
    t.Parallel()
    p := NewTsvPresenter(nil)
    h := p.makeHeader()
    if len(h) != 1 { t.Fatalf("expected single header row") }
    if !strings.Contains(h[0], "\t") { t.Fatalf("expected tabs in header: %q", h[0]) }
}

