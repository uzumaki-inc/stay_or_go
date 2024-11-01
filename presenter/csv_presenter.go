package presenter

type CsvPresenter struct {
	analyzedLibInfos []AnalyzedLibInfo
}

func (p CsvPresenter) Display() {
	Display(p)
}

func (p CsvPresenter) makeHeader() []string {
	return make([]string, 0)
}

func (p CsvPresenter) makeBody() []string {
	return make([]string, 0)
}
