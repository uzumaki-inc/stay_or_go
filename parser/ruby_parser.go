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
	"strings"
	"time"

	"github.com/konyu/StayOrGo/utils"
)

type RubyParser struct{}

type RubyRepository struct {
	SourceCodeURI string `json:"source_code_uri"`
	HomepageURI   string `json:"homepage_uri"`
}

func (p RubyParser) Parse(filePath string) ([]LibInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		utils.StdErrorPrintln("%v: %v", ErrFailedToReadFile, err)

		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	libInfoList, err := p.processFile(file)
	if err != nil {
		return nil, err
	}

	return libInfoList, nil
}

func (p RubyParser) processFile(file *os.File) ([]LibInfo, error) {
	var libInfoList []LibInfo

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		libInfo, err := p.processLine(line)
		if err != nil {
			utils.StdErrorPrintln("Error processing line: %v", err)

			continue
		}

		if libInfo != nil {
			libInfoList = append(libInfoList, *libInfo)
		}
	}

	if err := scanner.Err(); err != nil {
		utils.StdErrorPrintln("%v: %v", ErrFailedToScanFile, err)

		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	return libInfoList, nil
}

const minPartsLength = 2

func (p RubyParser) processLine(line string) (*LibInfo, error) {
	if !strings.HasPrefix(strings.TrimSpace(line), "gem ") {
		return nil, nil
	}

	parts := strings.Fields(line)
	if len(parts) < minPartsLength {
		return nil, ErrMissingGemName
	}

	gemName := strings.Trim(parts[1], `'" ,`)
	newLib := NewLibInfo(gemName)

	combinedParts := strings.Join(parts[2:], " ")
	splitByComma := strings.Split(combinedParts, ",")

	for _, part := range splitByComma {
		cleanedPart := strings.TrimSpace(part)
		if cleanedPart == "" {
			continue
		}

		ngKeys := []string{"source", "git", "github"}
		for _, ngKey := range ngKeys {
			if strings.HasPrefix(cleanedPart, ":"+ngKey+" ") || strings.HasPrefix(cleanedPart, ngKey+":") {
				newLib.Skip = true
				newLib.SkipReason = "does not support libraries hosted outside of Github"

				break
			}
		}

		newLib.Others = append(newLib.Others, cleanedPart)
	}

	return &newLib, nil
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
