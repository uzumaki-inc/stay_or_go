package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/konyu/StayOrGo/common"
)

type RubyParser struct{}

type Repository struct {
	SourceCodeURI string `json:"source_code_uri"`
	HomepageURI   string `json:"homepage_uri"`
}

func (p RubyParser) Parse(file string) []common.AnalyzedLibInfo {
	var AnalyzedLibInfoList []common.AnalyzedLibInfo

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
		gemName := strings.Trim(parts[1], `'" ,`)
		newLib := common.LibInfo{Name: gemName}
		newALI := common.AnalyzedLibInfo{Skip: false, LibInfo: newLib}

		// ここ以降のpartsはカンマ区切りでパースする
		combinedParts := strings.Join(parts[2:], " ")
		splitByComma := strings.Split(combinedParts, ",")

		for _, part := range splitByComma {
			cleanedPart := strings.TrimSpace(part)
			if cleanedPart == "" {
				continue
			}
			// fmt.Println(cleanedPart)

			// NGキーのリスト
			ngKeys := []string{"source", "git", "github"}

			// cleanedPartがハッシュ形式を表すかチェックし、NGキーが含まれているか判定
			for _, ngKey := range ngKeys {
				fmt.Println(ngKey)
				if strings.HasPrefix(cleanedPart, ":"+ngKey+" ") || strings.HasPrefix(cleanedPart, ngKey+":") {
					newALI.Skip = true
					newALI.SkipReason = "does not support libraries hosted outside of Github"

					break // NGキーが見つかったらこれ以上チェックする必要はない
				}
			}
			newLib.Others = append(newLib.Others, cleanedPart)
		}

		AnalyzedLibInfoList = append(AnalyzedLibInfoList, newALI)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return AnalyzedLibInfoList
}

func (p RubyParser) GetRepositoryURL(analyzedLibInfoList []common.AnalyzedLibInfo) []common.AnalyzedLibInfo {
	for i, _ := range analyzedLibInfoList {
		analyzedLibInfo := &analyzedLibInfoList[i] // ポインタを取得
		libInfo := &analyzedLibInfoList[i].LibInfo
		name := libInfo.Name

		if analyzedLibInfo.Skip {
			continue
		}

		repoURL, err := p.getGitHubRepositoryURL(name)
		if err != nil {
			analyzedLibInfo.Skip = true
			analyzedLibInfo.SkipReason = "does not support libraries hosted outside of Github"

			fmt.Printf("%s: %s\n", name, err.Error())
			continue
		}

		libInfo.RepositoryUrl = repoURL
		fmt.Printf("GitHub repository URL for %s: %s\n", name, repoURL)
	}
	return analyzedLibInfoList
}

func (p RubyParser) getGitHubRepositoryURL(name string) (string, error) {
	baseURL := "https://rubygems.org/api/v1/gems/"
	response, err := http.Get(baseURL + name + ".json")
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

	var repo Repository
	err = json.Unmarshal(bodyBytes, &repo)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response")
	}

	repoURL := repo.SourceCodeURI
	if repoURL == "" {
		repoURL = repo.HomepageURI
	}

	if repoURL == "" || !strings.Contains(repoURL, "github.com") {
		return "", fmt.Errorf("not a GitHub repository, skipping")
	}

	return repoURL, nil
}
