package presenter

import (
	"strings"
)

type MarkdownPresenter struct {
	analyzedLibInfos []AnalyzedLibInfo
}

func (p MarkdownPresenter) Display() {
	Display(p)
}

func (p MarkdownPresenter) makeHeader() []string {
	headerRow := "| " + strings.Join(headerString, " | ") + " |"

	separatorRow := "|"
	for _, header := range headerString {
		separatorRow += " " + strings.Repeat("-", len(header)) + " |"
	}
	return []string{headerRow, separatorRow}
}

func (p MarkdownPresenter) makeBody() []string {
	return makeBody(p.analyzedLibInfos, "|")
}
