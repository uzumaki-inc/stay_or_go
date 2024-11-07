package analyzer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/konyu/StayOrGo/utils"
)

type GitHubRepoInfo struct {
	RepositoryName   string
	Watchers         int
	Stars            int
	Forks            int
	OpenPullRequests int
	OpenIssues       int
	LastCommitDate   string
	GithubRepoUrl    string
	Archived         bool
	Score            int
	Skip             bool   // スキップするかどうかのフラグ
	SkipReason       string // スキップ理由
}

type GitHubRepoAnalyzer struct {
	githubToken string
	weights     ParameterWeights
}

func NewGitHubRepoAnalyzer(token string, weights ParameterWeights) *GitHubRepoAnalyzer {
	return &GitHubRepoAnalyzer{
		githubToken: token,
		weights:     weights,
	}
}

// FetchInfo fetches information for each repository
func (g *GitHubRepoAnalyzer) FetchGithubInfo(repositoryUrls []string) []GitHubRepoInfo {
	var libraryInfoList []GitHubRepoInfo
	for _, repoUrl := range repositoryUrls {
		utils.DebugPrintln("Fetching: " + repoUrl)

		libraryInfo, err := g.getGitHubInfo(repoUrl)
		if err != nil {
			libraryInfo = &GitHubRepoInfo{
				Skip:       true,
				SkipReason: "Failed fetching " + repoUrl + " from GitHub",
			}
			utils.StdErrorPrintln("Failed fetching %s, error details: %v", repoUrl, err)
		}

		libraryInfo.GithubRepoUrl = repoUrl
		libraryInfoList = append(libraryInfoList, *libraryInfo)
	}
	return libraryInfoList
}

// getGitHubInfo fetches repository info from GitHub API
func (g *GitHubRepoAnalyzer) getGitHubInfo(repoUrl string) (*GitHubRepoInfo, error) {
	if g.githubToken == "" {
		return nil, fmt.Errorf("GitHub token not set")
	}

	repoUrl = strings.TrimSuffix(repoUrl, "/")
	parts := strings.Split(repoUrl, "/")

	var owner, repo string
	if strings.Contains(repoUrl, "/tree/") {
		baseIndex := indexOf(parts, "github.com") + 1
		owner, repo = parts[baseIndex], parts[baseIndex+1]
	} else {
		owner, repo = parts[len(parts)-2], parts[len(parts)-1]
	}

	repo = strings.TrimSuffix(repo, ".git")

	client := &http.Client{}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("token %s", g.githubToken),
	}

	repoData, err := fetchJSON(client, fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), headers)
	if err != nil {
		return nil, err
	}

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

	repoInfo := &GitHubRepoInfo{
		RepositoryName:   repoData["name"].(string),
		Watchers:         int(repoData["subscribers_count"].(float64)),
		Stars:            int(repoData["stargazers_count"].(float64)),
		Forks:            int(repoData["forks_count"].(float64)),
		OpenPullRequests: openPullRequests,
		OpenIssues:       int(repoData["open_issues_count"].(float64)),
		LastCommitDate:   lastCommitDate,
		Archived:         repoData["archived"].(bool),
		Skip:             false,
		SkipReason:       "",
	}
	calcScore(repoInfo, &g.weights)

	return repoInfo, nil
}

func calcScore(repoInfo *GitHubRepoInfo, weights *ParameterWeights) {
	days, err := daysSince(repoInfo.LastCommitDate)
	if err != nil {
		repoInfo.Skip = true
		repoInfo.SkipReason = "Date Format Error: " + repoInfo.LastCommitDate
		utils.StdErrorPrintln("Date Format Error: %v", err)
	}

	var score = float64(repoInfo.Stars) * weights.Stars
	score += float64(repoInfo.Forks) * weights.Forks
	score += float64(repoInfo.OpenPullRequests) * weights.OpenPullRequests
	score += float64(repoInfo.OpenIssues) * weights.OpenIssues
	score += float64(days) * weights.LastCommitDate
	intArchived := map[bool]float64{true: 1.0, false: 0.0}[repoInfo.Archived]
	score += (intArchived) * weights.Archived

	repoInfo.Score = int(score)
}

// 日付文字列から現在日までの経過日数を返す関数
func daysSince(dateStr string) (int, error) {
	// 入力された日付文字列をパース（UTCフォーマット）
	layout := "2006-01-02T15:04:05Z"
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return 0, err
	}
	currentTime := time.Now()
	duration := currentTime.Sub(parsedTime)
	days := int(duration.Hours() / 24)

	return days, nil
}

func fetchJSONData(client *http.Client, url string, headers map[string]string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}

// fetchJSON sends a GET request and returns the parsed JSON object (map)
func fetchJSON(client *http.Client, url string, headers map[string]string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := fetchJSONData(client, url, headers, &result)
	return result, err
}

// fetchJSONArray sends a GET request and returns the parsed JSON array (slice)
func fetchJSONArray(client *http.Client, url string, headers map[string]string) ([]interface{}, error) {
	var result []interface{}
	err := fetchJSONData(client, url, headers, &result)
	return result, err
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
