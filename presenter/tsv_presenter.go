package presenter

import (
	"strings"
)

type TsvPresenter struct {
	analyzedLibInfos []AnalyzedLibInfo
}

func NewTsvPresenter(infos []AnalyzedLibInfo) TsvPresenter {
	return TsvPresenter{analyzedLibInfos: infos}
}

func (p TsvPresenter) Display() {
	Display(p)
}

func (p TsvPresenter) makeHeader() []string {
	headerRow := strings.Join(headerString, "\t")

	return []string{headerRow}
}

func (p TsvPresenter) makeBody() []string {
	return makeBody(p.analyzedLibInfos, "\t")
}
