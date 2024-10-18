package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type RubyParser struct{}

func (p RubyParser) Parse(file string) []LibInfo {
	var libInfoList []LibInfo

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(strings.TrimSpace(line), "gem ") {
			continue
		}

		parts := strings.Fields(line)

		if len(parts) < 2 {
			continue
		}
		gemName := parts[1]
		newLib := LibInfo{Name: gemName, Skip: false}

		// ここ以降のpartsはカンマ区切りでパースする
		combinedParts := strings.Join(parts[2:], " ")
		splitByComma := strings.Split(combinedParts, ",")

		for _, part := range splitByComma {
			cleanedPart := strings.TrimSpace(part)
			if cleanedPart == "" {
				continue
			}
			fmt.Println(cleanedPart)

			// NGキーのリスト
			ngKeys := []string{"source", "git", "github"}

			// cleanedPartがハッシュ形式を表すかチェックし、NGキーが含まれているか判定
			for _, ngKey := range ngKeys {
				if strings.HasPrefix(cleanedPart, ":"+ngKey+" ") || strings.HasPrefix(cleanedPart, ngKey+":") {
					newLib.Skip = true
					break // NGキーが見つかったらこれ以上チェックする必要はない
				}
			}
			newLib.Others = append(newLib.Others, cleanedPart)
		}

		libInfoList = append(libInfoList, newLib)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return libInfoList
}
