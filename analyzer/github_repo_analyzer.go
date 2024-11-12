package analyzer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/konyu/StayOrGo/utils"
)

var (
	ErrGitHubTokenNotSet           = errors.New("GitHub token not set")
	ErrFailedToAssertDefaultBranch = errors.New("failed to assert type for default_branch")
	ErrFailedToAssertDate          = errors.New("failed to assert type for date")

	ErrFailedToAssertName             = errors.New("failed to assert type for name")
	ErrFailedToAssertSubscribersCount = errors.New("failed to assert type for subscribers_count")
	ErrFailedToAssertStargazersCount  = errors.New("failed to assert type for stargazers_count")
	ErrFailedToAssertForksCount       = errors.New("failed to assert type for forks_count")
	ErrFailedToAssertOpenIssuesCount  = errors.New("failed to assert type for open_issues_count")
	ErrFailedToAssertArchived         = errors.New("failed to assert type for archived")
)

type GitHubRepoInfo struct {
	RepositoryName string
	Watchers       int
	Stars          int
	Forks          int
	OpenIssues     int
	LastCommitDate string
	GithubRepoURL  string
	Archived       bool
	Score          int
	Skip           bool   // スキップするかどうかのフラグ
	SkipReason     string // スキップ理由
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
	libraryInfoList := make([]GitHubRepoInfo, 0, len(repositoryUrls))
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), timeOutSec*time.Second)

	defer cancel()

	for _, repoURL := range repositoryUrls {
		utils.DebugPrintln("Fetching: " + repoURL)

		libraryInfo, err := g.getGitHubInfo(ctx, client, repoURL)
		if err != nil {
			libraryInfo = &GitHubRepoInfo{
				Skip:       true,
				SkipReason: "Failed fetching " + repoURL + " from GitHub",
			}

			utils.StdErrorPrintln("Failed fetching %s, error details: %v", repoURL, err)
		}

		libraryInfo.GithubRepoURL = repoURL
		libraryInfoList = append(libraryInfoList, *libraryInfo)
	}

	return libraryInfoList
}

const timeOutSec = 5

func (g *GitHubRepoAnalyzer) getGitHubInfo(
	ctx context.Context,
	client *http.Client,
	repoURL string,
) (*GitHubRepoInfo, error) {
	if g.githubToken == "" {
		return nil, ErrGitHubTokenNotSet
	}

	owner, repo := parseRepoURL(repoURL)

	headers := map[string]string{
		"Authorization": "token " + g.githubToken,
	}

	repoData, err := fetchRepoData(ctx, client, owner, repo, headers)
	if err != nil {
		return nil, err
	}

	lastCommitDate, err := fetchLastCommitDate(ctx, client, owner, repo, repoData, headers)
	if err != nil {
		return nil, err
	}

	repoInfo, err := createRepoInfo(repoData, lastCommitDate)
	if err != nil {
		return nil, err
	}

	calcScore(repoInfo, &g.weights)

	return repoInfo, nil
}

func parseRepoURL(repoURL string) (string, string) {
	repoURL = strings.TrimSuffix(repoURL, "/")
	parts := strings.Split(repoURL, "/")

	var owner, repo string

	if strings.Contains(repoURL, "/tree/") {
		baseIndex := indexOf(parts, "github.com") + 1
		owner, repo = parts[baseIndex], parts[baseIndex+1]
	} else {
		owner, repo = parts[len(parts)-2], parts[len(parts)-1]
	}

	repo = strings.TrimSuffix(repo, ".git")

	return owner, repo
}

func fetchRepoData(
	ctx context.Context,
	client *http.Client,
	owner, repo string,
	headers map[string]string,
) (map[string]interface{}, error) {
	return fetchJSON(ctx, client, fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), headers)
}

func fetchLastCommitDate(ctx context.Context, client *http.Client, owner, repo string,
	repoData map[string]interface{}, headers map[string]string) (string, error) {
	defaultBranch, isString := repoData["default_branch"].(string)
	if !isString {
		return "", ErrFailedToAssertDefaultBranch
	}

	commitURL := "https://api.github.com/repos/" + owner + "/" + repo + "/commits/" + defaultBranch
	commitData, err := fetchJSON(ctx, client, commitURL, headers)

	if err != nil {
		return "", err
	}

	lastCommitDate, ok := commitData["commit"].(map[string]interface{})["committer"].(map[string]interface{})["date"].(string)
	if !ok {
		return "", ErrFailedToAssertDate
	}

	return lastCommitDate, nil
}

func createRepoInfo(
	repoData map[string]interface{},
	lastCommitDate string,
) (*GitHubRepoInfo, error) {
	repoName, ok := repoData["name"].(string)
	if !ok {
		return nil, ErrFailedToAssertName
	}

	subscribersCount, ok := repoData["subscribers_count"].(float64)
	if !ok {
		return nil, ErrFailedToAssertSubscribersCount
	}

	stargazersCount, ok := repoData["stargazers_count"].(float64)
	if !ok {
		return nil, ErrFailedToAssertStargazersCount
	}

	forksCount, ok := repoData["forks_count"].(float64)
	if !ok {
		return nil, ErrFailedToAssertForksCount
	}

	openIssuesCount, ok := repoData["open_issues_count"].(float64)
	if !ok {
		return nil, ErrFailedToAssertOpenIssuesCount
	}

	archived, ok := repoData["archived"].(bool)
	if !ok {
		return nil, ErrFailedToAssertArchived
	}

	return &GitHubRepoInfo{
		RepositoryName: repoName,
		Watchers:       int(subscribersCount),
		Stars:          int(stargazersCount),
		Forks:          int(forksCount),
		OpenIssues:     int(openIssuesCount),
		LastCommitDate: lastCommitDate,
		Archived:       archived,
		Skip:           false,
		SkipReason:     "",
	}, nil
}

func calcScore(repoInfo *GitHubRepoInfo, weights *ParameterWeights) {
	days, err := daysSince(repoInfo.LastCommitDate)
	if err != nil {
		repoInfo.Skip = true

		repoInfo.SkipReason = "Date Format Error: " + repoInfo.LastCommitDate

		utils.StdErrorPrintln("Date Format Error: %v", err)
	}

	score := float64(repoInfo.Stars) * weights.Stars
	score += float64(repoInfo.Forks) * weights.Forks
	score += float64(repoInfo.OpenIssues) * weights.OpenIssues
	score += float64(days) * weights.LastCommitDate
	intArchived := map[bool]float64{true: 1.0, false: 0.0}[repoInfo.Archived]
	score += (intArchived) * weights.Archived

	repoInfo.Score = int(score)
}

const hoursOfDay = 24

// 日付文字列から現在日までの経過日数を返す関数
func daysSince(dateStr string) (int, error) {
	// 入力された日付文字列をパース（UTCフォーマット）
	layout := "2006-01-02T15:04:05Z"
	parsedTime, err := time.Parse(layout, dateStr)

	if err != nil {
		return 0, fmt.Errorf("failed to parse date '%s': %w", dateStr, err)
	}

	currentTime := time.Now()
	duration := currentTime.Sub(parsedTime)
	days := int(duration.Hours() / hoursOfDay)

	return days, nil
}

func fetchJSONData(
	ctx context.Context,
	client *http.Client,
	url string,
	headers map[string]string,
	result interface{},
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return fmt.Errorf("failed to create new HTTP request for URL %s: %w", url, err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request for URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode JSON response for URL %s: %w", url, err)
	}

	return nil
}

// fetchJSON sends a GET request and returns the parsed JSON object (map)
func fetchJSON(
	ctx context.Context,
	client *http.Client,
	url string,
	headers map[string]string,
) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := fetchJSONData(ctx, client, url, headers, &result)

	return result, err
}

// fetchJSONArray sends a GET request and returns the parsed JSON array (slice)
func fetchJSONArray(
	ctx context.Context,
	client *http.Client,
	url string,
	headers map[string]string,
) ([]interface{}, error) {
	var result []interface{}
	err := fetchJSONData(ctx, client, url, headers, &result)

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
