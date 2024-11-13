package presenter

import (
	"strings"
)

type CsvPresenter struct {
	analyzedLibInfos []AnalyzedLibInfo
}

func (p CsvPresenter) Display() {
	Display(p)
}

func (p CsvPresenter) makeHeader() []string {
	headerRow := strings.Join(headerString, ", ")

	return []string{headerRow}
}

func (p CsvPresenter) makeBody() []string {
	return makeBody(p.analyzedLibInfos, ", ")
}

func NewCsvPresenter(infos []AnalyzedLibInfo) CsvPresenter {
	return CsvPresenter{analyzedLibInfos: infos}
}
