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

func (p GoParser) Parse(file string) []LibInfo {
	var libInfoList []LibInfo

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "\t") || strings.HasSuffix(line, "// indirect") {
			continue
		}

		parts := strings.Fields(line)

		if len(parts) < 2 {
			continue
		}

		libParts := strings.Split(parts[0], "/")
		libName := libParts[len(libParts)-1]
		newLib := LibInfo{Name: libName, Others: []string{parts[0], parts[1]}}

		libInfoList = append(libInfoList, newLib)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return libInfoList
}

func (p GoParser) GetRepositoryURL(libInfoList []LibInfo) []LibInfo {
	for i, _ := range libInfoList {
		libInfo := &libInfoList[i]
		name := libInfo.Others[0]
		version := libInfo.Others[1]

		if libInfo.Skip {
			continue
		}

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
