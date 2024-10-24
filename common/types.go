package common

type AnalyzedLibInfo struct {
	Skip           bool   // スキップするかどうかのフラグ
	SkipReason     string // スキップ理由
	LibInfo        LibInfo
	GitHubRepoInfo GitHubRepoInfo
}

type LibInfo struct {
	Name          string   // ライブラリの名前
	Others        []string // その他のライブラリの設定値
	RepositoryUrl string   // githubのりポトリのURL
}

type GitHubRepoInfo struct {
	RepositoryName   string `json:"repository_name"`
	Watchers         int    `json:"watchers"`
	Stars            int    `json:"stars"`
	Forks            int    `json:"forks"`
	OpenPullRequests int    `json:"open_pull_requests"`
	OpenIssues       int    `json:"open_issues"`
	LastCommitDate   string `json:"last_commit_date"`
	LibraryName      string `json:"library_name"`
	GithubRepoUrl    string `json:"github_repo_url"`
	Archived         bool   `json:"archived"`
}
