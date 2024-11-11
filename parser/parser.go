package parser

import (
	"os"

	"github.com/konyu/StayOrGo/utils"
)

type LibInfo struct {
	Skip          bool     // スキップするかどうかのフラグ
	SkipReason    string   // スキップ理由
	Name          string   // ライブラリの名前
	Others        []string // その他のライブラリの設定値
	RepositoryUrl string   // githubのりポトリのURL
}

type Parser interface {
	Parse(file string) []LibInfo
	GetRepositoryURL(AnalyzedLibInfoList []LibInfo) []LibInfo
}

var selectedParser Parser

func SelectParser(language string) Parser {
	var parser Parser
	switch language {
	case "ruby":
		parser = RubyParser{}
	case "go":
		parser = GoParser{}
	default:
		utils.StdErrorPrintln("Error: Unsupported language: %s", language)
		os.Exit(1)
	}
	selectedParser = parser
	return parser
}

func Parse(file string) []LibInfo {
	if selectedParser == nil {
		utils.StdErrorPrintln("Error: Parser not selected")
		os.Exit(1)
	}
	return selectedParser.Parse(file)
}
