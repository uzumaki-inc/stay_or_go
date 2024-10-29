package analyzer

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

// stars * 0.1
// forks * 0.1
// open_pull_requests * 0.01 //OpenしているPRの数 startsやwatcherに比べると影響は少ない
// open_issues* 0.01 //OpenしているIssueの数
// 実行日からlast_commit_dateまでの日数  * 0.2) //かつて人気があってもメンテナンスされていないものはスコアが下がるように調整
func NewParameterWeights() ParameterWeights {
	return ParameterWeights{
		Watchers:         0.1,
		Stars:            0.1,
		Forks:            0.1,
		OpenPullRequests: 0.01,
		OpenIssues:       0.01,
		LastCommitDate:   -0.2,
		Archived:         -100000000,
	}
}
