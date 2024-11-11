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

type GoParser struct{}

func (p GoParser) Parse(filePath string) []LibInfo {
	file, err := os.Open(filePath)
	if err != nil {
		utils.StdErrorPrintln("Failed to read file %v", err)
		os.Exit(1)
	}
	defer file.Close()

	replaceModules := p.collectReplaceModules(file)
	libInfoList := p.processRequireBlock(file, replaceModules)

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

	file.Seek(0, 0) // Reset file pointer for next pass
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
					newLib = LibInfo{Name: libName, Skip: true, SkipReason: "replaced module"}
				} else {
					newLib = LibInfo{Name: libName, Others: []string{parts[0], parts[1]}}
				}
				libInfoList = append(libInfoList, newLib)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		utils.StdErrorPrintln("Failed to scan file %v", err)
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

func (p GoParser) GetRepositoryURL(libInfoList []LibInfo) []LibInfo {
	for i, _ := range libInfoList {
		libInfo := &libInfoList[i]

		if libInfo.Skip {
			continue
		}
		name := libInfo.Others[0]
		version := libInfo.Others[1]

		repoURL, err := p.getGitHubRepositoryURL(name, version)
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

type GoRepository struct {
	Version string `json:"Version"`
	Time    string `json:"Time"`
	Origin  Origin `json:"Origin"`
}

type Origin struct {
	VCS  string `json:"VCS"`
	URL  string `json:"URL"`
	Ref  string `json:"Ref"`
	Hash string `json:"Hash"`
}

func (p GoParser) getGitHubRepositoryURL(name, version string) (string, error) {
	baseURL := "https://proxy.golang.org/"
	repoUrl := baseURL + name + "/@v/" + version + ".info"
	utils.DebugPrintln("Fetching: " + repoUrl)
	response, err := http.Get(repoUrl)

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

	var repo GoRepository
	err = json.Unmarshal(bodyBytes, &repo)

	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response")
	}

	repoURL := repo.Origin.URL

	if repoURL == "" || !strings.Contains(repoURL, "github.com") {
		return "", fmt.Errorf("not a GitHub repository, skipping")
	}

	return repoURL, nil
}
