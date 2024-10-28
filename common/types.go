package common

// type AnalyzedLibInfo struct {
// 	Skip           bool   // スキップするかどうかのフラグ
// 	SkipReason     string // スキップ理由
// 	LibInfo        LibInfo
// 	GitHubRepoInfo GitHubRepoInfo
// }

// type LibInfo struct {
// 	Skip          bool     // スキップするかどうかのフラグ
// 	SkipReason    string   // スキップ理由
// 	Name          string   // ライブラリの名前
// 	Others        []string // その他のライブラリの設定値
// 	RepositoryUrl string   // githubのりポトリのURL
// }

type ParameterWeights struct {
	Watchers         float64
	Stars            float64
	Forks            float64
	OpenPullRequests float64
	OpenIssues       float64
	LastCommitDate   float64
	Archived         float64
	Score            float64
}

func NewParameterWeights() ParameterWeights {
	return ParameterWeights{
		Watchers:         0.1,
		Stars:            0.1,
		Forks:            0.1,
		OpenPullRequests: 0.1,
		OpenIssues:       0.1,
		LastCommitDate:   0.1,
		Archived:         0.1,
	}
}
