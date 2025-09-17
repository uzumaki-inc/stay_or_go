package parser

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/uzumaki-inc/stay_or_go/utils"
)

type GoParser struct{}

func (p GoParser) Parse(filePath string) ([]LibInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		utils.StdErrorPrintln("%v: %v", ErrFailedToReadFile, err)
		os.Exit(1)
	}
	defer file.Close()

	replaceModules := p.collectReplaceModules(file)
	libInfoList := p.processRequireBlock(file, replaceModules)

	return libInfoList, nil
}

func (p GoParser) GetRepositoryURL(libInfoList []LibInfo) []LibInfo {
	client := &http.Client{}

	for i := range libInfoList {
		libInfo := &libInfoList[i]

		if libInfo.Skip {
			continue
		}

		name := libInfo.Others[0]
		version := libInfo.Others[1]

		repoURL, err := p.getGitHubRepositoryURL(client, name, version)
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

func (p GoParser) collectReplaceModules(file *os.File) []string {
	var replaceModules []string

	var inReplaceBlock bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "replace (" {
			inReplaceBlock = true

			continue
		}

		if line == ")" && inReplaceBlock {
			inReplaceBlock = false

			continue
		}

		if inReplaceBlock {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				replaceModules = append(replaceModules, parts[0])
			}
		}
	}

	_, err := file.Seek(0, 0) // Reset file pointer for next pass
	if err != nil {
		utils.StdErrorPrintln("%v: %v", ErrFailedToResetFilePointer, err)
		os.Exit(1)
	}

	return replaceModules
}

func (p GoParser) processRequireBlock(file *os.File, replaceModules []string) []LibInfo {
	var libInfoList []LibInfo

	var inRequireBlock bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "require (" {
			inRequireBlock = true

			continue
		}

		if line == ")" && inRequireBlock {
			inRequireBlock = false

			continue
		}

		if inRequireBlock && !strings.Contains(line, "// indirect") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				module := parts[0]
				libParts := strings.Split(parts[0], "/")
				libName := libParts[len(libParts)-1]

				var newLib LibInfo

				if contains(replaceModules, module) {
					newLib = NewLibInfo(libName, WithSkip(true), WithSkipReason("replaced module"))
				} else {
					newLib = NewLibInfo(libName, WithOthers([]string{parts[0], parts[1]}))
				}

				libInfoList = append(libInfoList, newLib)
			}
		}
	}

	err := scanner.Err()
	if err != nil {
		utils.StdErrorPrintln("%v: %v", ErrFailedToScanFile, err)
		os.Exit(1)
	}

	return libInfoList
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

type GoRepository struct {
	Version string `json:"version"`
	Time    string `json:"time"`
	Origin  Origin `json:"origin"`
}

type Origin struct {
	VCS  string `json:"vcs"`
	URL  string `json:"url"`
	Ref  string `json:"ref"`
	Hash string `json:"hash"`
}

func (p GoParser) getGitHubRepositoryURL(
	client *http.Client,
	name,
	version string,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOutSec*time.Second)
	defer cancel()

	baseURL := "https://proxy.golang.org/"
	repoURL := baseURL + name + "/@v/" + version + ".info"
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

	repoURLfromGithub, err := extractRepoURL(bodyBytes, name)
	if err != nil {
		return "", err
	}

	return repoURLfromGithub, nil
}

func extractRepoURL(bodyBytes []byte, name string) (string, error) {
	var repo GoRepository

	err := json.Unmarshal(bodyBytes, &repo)
	if err != nil {
		return "", ErrFailedToUnmarshalJSON
	}

	repoURLfromGithub := repo.Origin.URL

	// If there is no URL, use the package name
	if repoURLfromGithub == "" && strings.Contains(name, "github.com") {
		repoURLfromGithub = "https://" + name
	}

	if repoURLfromGithub == "" || !strings.Contains(repoURLfromGithub, "github.com") {
		return "", ErrNotAGitHubRepository
	}

	return repoURLfromGithub, nil
}
