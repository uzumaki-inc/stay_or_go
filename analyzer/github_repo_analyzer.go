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
	ErrUnexpectedStatusCode           = errors.New("unexpected status code")
)

const (
	hoursOfDay = 24
	timeOutSec = 5
)

type RepoData struct {
	Name             string `json:"name"`
	SubscribersCount int    `json:"subscribers_count"`
	StargazersCount  int    `json:"stargazers_count"`
	ForksCount       int    `json:"forks_count"`
	OpenIssuesCount  int    `json:"open_issues_count"`
	Archived         bool   `json:"archived"`
	DefaultBranch    string `json:"default_branch"`
}

type CommitData struct {
	Commit struct {
		Committer struct {
			Date string `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
}

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

	for _, repoURL := range repositoryUrls {
		utils.DebugPrintln("Fetching: " + repoURL)

		libraryInfo, err := g.getGitHubInfo(client, repoURL)
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

func (g *GitHubRepoAnalyzer) getGitHubInfo(
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

	repoData, err := fetchRepoData(client, owner, repo, headers)
	if err != nil {
		return nil, err
	}

	lastCommitDate, err := fetchLastCommitDate(client, owner, repo, repoData, headers)
	if err != nil {
		return nil, err
	}

	repoInfo := createRepoInfo(repoData, lastCommitDate)

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
	client *http.Client,
	owner, repo string,
	headers map[string]string,
) (*RepoData, error) {
	var repoData RepoData

	err := fetchJSONData(client, fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo), headers, &repoData)
	if err != nil {
		return nil, err
	}

	return &repoData, nil
}

func fetchLastCommitDate(client *http.Client, owner, repo string,
	repoData *RepoData, headers map[string]string,
) (string, error) {
	commitURL := "https://api.github.com/repos/" + owner + "/" + repo + "/commits/" + repoData.DefaultBranch

	var commitData CommitData

	err := fetchJSONData(client, commitURL, headers, &commitData)
	if err != nil {
		return "", err
	}

	return commitData.Commit.Committer.Date, nil
}

func createRepoInfo(
	repoData *RepoData,
	lastCommitDate string,
) *GitHubRepoInfo {
	return &GitHubRepoInfo{
		RepositoryName: repoData.Name,
		Watchers:       repoData.SubscribersCount,
		Stars:          repoData.StargazersCount,
		Forks:          repoData.ForksCount,
		OpenIssues:     repoData.OpenIssuesCount,
		LastCommitDate: lastCommitDate,
		Archived:       repoData.Archived,
		Skip:           false,
		SkipReason:     "",
	}
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
	client *http.Client,
	url string,
	headers map[string]string,
	result interface{},
) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOutSec*time.Second)
	defer cancel()

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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d for URL %s", ErrUnexpectedStatusCode, resp.StatusCode, url)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode JSON response for URL %s: %w", url, err)
	}

	return nil
}

func indexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}

	return -1
}
