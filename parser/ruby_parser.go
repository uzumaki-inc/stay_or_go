package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/konyu/StayOrGo/utils"
)

type RubyParser struct{}

type RubyRepository struct {
	SourceCodeURI string `json:"source_code_uri"`
	HomepageURI   string `json:"homepage_uri"`
}

func (p RubyParser) Parse(file string) []LibInfo {
	var libInfoList []LibInfo
	const minPartsLength = 2

	f, err := os.Open(file)
	if err != nil {
		utils.StdErrorPrintln("Failed to read file %v", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(strings.TrimSpace(line), "gem ") {
			continue
		}

		parts := strings.Fields(line)

		if len(parts) < minPartsLength {
			continue
		}

		gemName := strings.Trim(parts[1], `'" ,`)
		newLib := LibInfo{Name: gemName}

		// ここ以降のpartsはカンマ区切りでパースする
		combinedParts := strings.Join(parts[2:], " ")
		splitByComma := strings.Split(combinedParts, ",")

		for _, part := range splitByComma {
			cleanedPart := strings.TrimSpace(part)
			if cleanedPart == "" {
				continue
			}
			// NGキーのリスト
			ngKeys := []string{"source", "git", "github"}

			// cleanedPartがハッシュ形式を表すかチェックし、NGキーが含まれているか判定
			for _, ngKey := range ngKeys {
				if strings.HasPrefix(cleanedPart, ":"+ngKey+" ") || strings.HasPrefix(cleanedPart, ngKey+":") {
					newLib.Skip = true
					newLib.SkipReason = "does not support libraries hosted outside of Github"

					break // NGキーが見つかったらこれ以上チェックする必要はない
				}
			}

			newLib.Others = append(newLib.Others, cleanedPart)
		}

		libInfoList = append(libInfoList, newLib)
	}

	if err := scanner.Err(); err != nil {
		utils.StdErrorPrintln("Failed to scan file %v", err)
		os.Exit(1)
	}

	return libInfoList
}

func (p RubyParser) GetRepositoryURL(libInfoList []LibInfo) []LibInfo {
	for i := range libInfoList {
		// ポインタを取得
		libInfo := &libInfoList[i]
		name := libInfo.Name

		if libInfo.Skip {
			continue
		}

		repoURL, err := p.getGitHubRepositoryURL(name)
		if err != nil {
			libInfo.Skip = true
			libInfo.SkipReason = "Does not support libraries hosted outside of Github"

			utils.StdErrorPrintln("%s does not support libraries hosted outside of Github: %s", name, err)

			continue
		}

		libInfo.RepositoryUrl = repoURL
	}

	return libInfoList
}

func (p RubyParser) getGitHubRepositoryURL(name string) (string, error) {
	baseURL := "https://rubygems.org/api/v1/gems/"
	repoURL := baseURL + name + ".json"
	utils.DebugPrintln("Fetching: " + repoURL)
	response, err := http.Get(repoURL)

	if err != nil {
		return "", fmt.Errorf("can't get the gem repository, skipping")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("not a GitHub repository, skipping")
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body")
	}

	var repo RubyRepository
	err = json.Unmarshal(bodyBytes, &repo)

	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response")
	}

	repoURLfromRubyGems := repo.SourceCodeURI
	if repoURLfromRubyGems == "" {
		repoURLfromRubyGems = repo.HomepageURI
	}

	if repoURLfromRubyGems == "" || !strings.Contains(repoURLfromRubyGems, "github.com") {
		return "", fmt.Errorf("not a GitHub repository, skipping")
	}

	return repoURLfromRubyGems, nil
}
