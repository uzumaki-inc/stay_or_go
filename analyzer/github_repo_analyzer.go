package analyzer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/konyu/StayOrGo/common"
)

type GitHubRepoAnalyzer struct {
	githubToken string
}

func NewGitHubRepoAnalyzer(token string) *GitHubRepoAnalyzer {
	return &GitHubRepoAnalyzer{
		githubToken: token,
	}
}

// FetchInfo fetches information for each repository
func (g *GitHubRepoAnalyzer) FetchGithubInfo(analyzedLibInfo []common.AnalyzedLibInfo) []common.GitHubRepoInfo {
	var libraryInfoList []common.GitHubRepoInfo
	for _, info := range analyzedLibInfo {
		if info.Skip {
			continue
		}
		name := info.LibInfo.Name
		repoUrl := info.LibInfo.RepositoryUrl

		fmt.Printf("Getting GitHub info for %s: %s\n", info.LibInfo.Name, info.LibInfo.RepositoryUrl)

		libraryInfo, err := g.getGitHubInfo(repoUrl)
		if err != nil {
			info.Skip = true
			info.SkipReason = "Failed getting" + name + "GitHub info:"
			fmt.Printf("Failed getting %s GitHub info: %v\n", name, err)
			continue
		}

		if libraryInfo != nil {
			libraryInfo.LibraryName = name
			libraryInfo.GithubRepoUrl = repoUrl
			libraryInfoList = append(libraryInfoList, *libraryInfo)
		}
	}
	return libraryInfoList
}

// getGitHubInfo fetches repository info from GitHub API
func (g *GitHubRepoAnalyzer) getGitHubInfo(repoUrl string) (*common.GitHubRepoInfo, error) {
	// githubToken := os.Getenv("GITHUB_TOKEN")
	if g.githubToken == "" {
		return nil, fmt.Errorf("GitHub token not set")
	}

	// Pre-process the repo URL
	repoUrl = strings.TrimSuffix(repoUrl, "/")
	parts := strings.Split(repoUrl, "/")

	var owner, repo string
	if strings.Contains(repoUrl, "/tree/") {
		baseIndex := indexOf(parts, "github.com") + 1
		owner, repo = parts[baseIndex], parts[baseIndex+1]
	} else {
		owner, repo = parts[len(parts)-2], parts[len(parts)-1]
	}

	// Remove .git if present
	repo = strings.TrimSuffix(repo, ".git")

	client := &http.Client{}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("token %s", g.githubToken),
	}

	repoData, err := fetchJSON(client, fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), headers)
	if err != nil {
		return nil, err
	}
	// fmt.Println(g.githubToken)
	// fmt.Println("Repo Requests Data:", repoData)

	pullRequestsData, err := fetchJSONArray(client, fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo), headers)
	if err != nil {
		return nil, err
	}

	openPullRequests := len(pullRequestsData)

	defaultBranch := repoData["default_branch"].(string)

	commitData, err := fetchJSON(client, fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, defaultBranch), headers)
	if err != nil {
		return nil, err
	}

	lastCommitDate := commitData["commit"].(map[string]interface{})["committer"].(map[string]interface{})["date"].(string)

	return &common.GitHubRepoInfo{
		RepositoryName:   repoData["name"].(string),
		Watchers:         int(repoData["watchers_count"].(float64)),
		Stars:            int(repoData["stargazers_count"].(float64)),
		Forks:            int(repoData["forks_count"].(float64)),
		OpenPullRequests: openPullRequests,
		OpenIssues:       int(repoData["open_issues_count"].(float64)),
		LastCommitDate:   lastCommitDate,
		Archived:         repoData["archived"].(bool),
	}, nil
}

// fetchJSON sends a GET request and returns the parsed JSON object (map)
func fetchJSON(client *http.Client, url string, headers map[string]string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// fetchJSONArray sends a GET request and returns the parsed JSON array (slice)
func fetchJSONArray(client *http.Client, url string, headers map[string]string) ([]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// indexOf returns the index of the element in a slice, or -1 if not found
func indexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}
