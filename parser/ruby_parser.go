package parser

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/konyu/StayOrGo/utils"
)

type RubyParser struct{}

type RubyRepository struct {
	SourceCodeURI string `json:"source_code_uri"`
	HomepageURI   string `json:"homepage_uri"`
}

// Parse メソッド
func (p RubyParser) Parse(filePath string) ([]LibInfo, error) {
	lines, err := p.readLines(filePath)
	if err != nil {
		return nil, err
	}

	var libs []LibInfo

	inOtherBlock := false

	for _, line := range lines {
		if p.isOtherBlockStart(line) {
			inOtherBlock = true

			continue
		}

		if p.isBlockEnd(line) {
			inOtherBlock = false

			continue
		}

		// gem を解析
		if gemName := p.extractGemName(line); gemName != "" {
			isNgGem := p.containsInvalidKeywords(line)
			lib := p.createLibInfo(gemName, isNgGem, inOtherBlock)
			libs = append(libs, lib)
		}
	}

	return libs, nil
}

// ファイルの内容を行ごとに読み取る
func (p *RubyParser) readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return lines, nil
}

// その他のブロックの開始か判定
func (p *RubyParser) isOtherBlockStart(line string) bool {
	sourceStartRegex := regexp.MustCompile(`source\s+['"].+['"]\s+do`)
	platformsStartRegex := regexp.MustCompile(`platforms\s+[:\w,]+\s+do`)
	installIfStartRegex := regexp.MustCompile(`install_if\s+->\s+\{.*\}\s+do`)

	return sourceStartRegex.MatchString(line) ||
		platformsStartRegex.MatchString(line) ||
		installIfStartRegex.MatchString(line)
}

// ブロックの終了か判定
func (p *RubyParser) isBlockEnd(line string) bool {
	endRegex := regexp.MustCompile(`^end$`)

	return endRegex.MatchString(line)
}

// gem 名を抽出
func (p *RubyParser) extractGemName(line string) string {
	gemRegex := regexp.MustCompile(`gem ['"]([^'"]+)['"]`)

	if matches := gemRegex.FindStringSubmatch(line); matches != nil {
		return matches[1]
	}

	return ""
}

func (p *RubyParser) containsInvalidKeywords(line string) bool {
	// カンマ区切りで分割
	parts := strings.Split(line, ",")

	// 判定するキーワード
	ngKeywords := []string{"source", "git", "github"}

	// 2番目以降をチェック
	for _, part := range parts[1:] {
		trimmedPart := strings.TrimSpace(part)
		for _, keyword := range ngKeywords {
			if strings.Contains(trimmedPart, keyword) {
				return true
			}
		}
	}

	return false
}

func (p *RubyParser) createLibInfo(gemName string, isNgGem bool, inOtherBlock bool) LibInfo {
	lib := LibInfo{Name: gemName}
	if isNgGem {
		lib.Skip = true
		lib.SkipReason = "Not hosted on Github"
	} else if inOtherBlock {
		lib.Skip = true
		lib.SkipReason = "Not hosted on Github"
	}

	return lib
}

func (p RubyParser) GetRepositoryURL(libInfoList []LibInfo) []LibInfo {
	client := &http.Client{}

	for i := range libInfoList {
		// ポインタを取得
		libInfo := &libInfoList[i]
		name := libInfo.Name

		if libInfo.Skip {
			continue
		}

		repoURL, err := p.getGitHubRepositoryURL(client, name)
		if err != nil {
			libInfo.Skip = true
			libInfo.SkipReason = "Does not support libraries hosted outside of Github"

			utils.StdErrorPrintln("%s does not support libraries hosted outside of Github: %s", name, err)

			continue
		}

		libInfo.RepositoryURL = repoURL
	}

	return libInfoList
}

func (p RubyParser) getGitHubRepositoryURL(client *http.Client, name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOutSec*time.Second)
	defer cancel()

	baseURL := "https://rubygems.org/api/v1/gems/"
	repoURL := baseURL + name + ".json"
	utils.DebugPrintln("Fetching: " + repoURL)

	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", ErrFailedToGetRepository
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return "", ErrFailedToGetRepository
	}

	response, err := client.Do(req)
	if err != nil {
		return "", ErrFailedToGetRepository
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", ErrNotAGitHubRepository
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", ErrFailedToReadResponseBody
	}

	var repo RubyRepository

	err = json.Unmarshal(bodyBytes, &repo)
	if err != nil {
		return "", ErrFailedToUnmarshalJSON
	}

	repoURLfromRubyGems := repo.SourceCodeURI

	if repoURLfromRubyGems == "" {
		repoURLfromRubyGems = repo.HomepageURI
	}

	if repoURLfromRubyGems == "" || !strings.Contains(repoURLfromRubyGems, "github.com") {
		return "", ErrNotAGitHubRepository
	}

	return repoURLfromRubyGems, nil
}
