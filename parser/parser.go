package parser

import (
	"github.com/konyu/StayOrGo/common"
)

type Parser interface {
	Parse(file string) []common.AnalyzedLibInfo
	GetRepositoryURL(AnalyzedLibInfoList []common.AnalyzedLibInfo) []common.AnalyzedLibInfo
}

var selectedParser Parser

func SelectParser(language string) Parser {
	var parser Parser
	switch language {
	case "ruby":
		parser = RubyParser{}
	case "go":
		// parser = GoParser{}
	default:
		panic("Error: Unsupported language: " + language)
	}
	selectedParser = parser
	return parser
}

func Parse(file string) []common.AnalyzedLibInfo {
	if selectedParser == nil {
		panic("Error: Parser not selected")
	}
	return selectedParser.Parse(file)
}
