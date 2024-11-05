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
// The number of open PRs and issues has less impact compared to stars and watchers
// open_pull_requests * 0.01
// open_issues * 0.01
// Adjusts the score to decrease for projects that were once popular but are no longer maintained
// Days from the execution date to the last commit date * 0.2
// If archived is true, it is not maintained, so it is heavily penalized
func NewParameterWeights() ParameterWeights {
	return ParameterWeights{
		Watchers:         0.1,
		Stars:            0.1,
		Forks:            0.1,
		OpenPullRequests: 0.01,
		OpenIssues:       0.01,
		LastCommitDate:   -0.2,
		Archived:         -1000000,
	}
}
