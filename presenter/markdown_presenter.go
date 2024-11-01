package presenter

import (
	"fmt"
	"reflect"
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
	rows := []string{}
	for _, info := range p.analyzedLibInfos {
		row := "| "
		// reflectを使用して、headerに対応するフィールドを取得
		val := reflect.ValueOf(info)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		for _, header := range headerString {
			method := val.MethodByName(header)
			if method.IsValid() {
				result := method.Call(nil)
				var resultStr interface{}
				if len(result) > 0 && result[0].IsValid() && !result[0].IsNil() {
					resultStr = result[0].Elem().Interface()
				} else {
					resultStr = "N/A"
				}
				row += fmt.Sprintf("%v | ", resultStr)
			} else {
				panic(fmt.Sprintf("method %s not found in %v", header, info))
			}
		}
		rows = append(rows, row)
	}
	return rows
}
