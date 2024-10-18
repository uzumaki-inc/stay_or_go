package parser

type LibInfo struct {
	Name   string   // ライブラリの名前
	Skip   bool     // 次の設定値をスキップするかどうかのフラグ
	Others []string // その他のライブラリの設定値
}

type Parser interface {
	Parse(file string) []LibInfo
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
		panic("Error: Unsupported language: " + language)
	}
	selectedParser = parser
	return parser
}

func Parse(file string) []LibInfo {
	if selectedParser == nil {
		panic("Error: Parser not selected")
	}
	return selectedParser.Parse(file)
}
