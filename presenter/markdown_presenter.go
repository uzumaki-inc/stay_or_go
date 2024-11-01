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
	header := p.makeHeader()
	for _, line := range header {
		fmt.Println(line)
	}
	body := p.makeBody()
	for _, line := range body {
		fmt.Println(line)
	}
}

func (p MarkdownPresenter) makeHeader() []string {
	fmt.Println(headerString)
	// ヘッダー行を `|` で区切って作成
	headerRow := "| " + strings.Join(headerString, " | ") + " |"

	// 区切り線を `| ---- |` 形式で作成
	separatorRow := "|"
	for _, header := range headerString {
		separatorRow += " " + strings.Repeat("-", len(header)) + " |"
	}

	// ヘッダーと区切り線を結合して返す
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
