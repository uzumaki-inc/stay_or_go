package analyzer

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultWatcherWeight         = 0.1
	defaultStarWeight            = 0.1
	defaultForkWeight            = 0.1
	defaultOpenPullRequestWeight = 0.01
	defaultOpenIssueWeight       = 0.01
	defaultLastCommitDateWeight  = -0.05
	defaultArchivedWeight        = -1000000
)

type ParameterWeights struct {
	Watchers       float64 `yaml:"watchers"`
	Stars          float64 `yaml:"stars"`
	Forks          float64 `yaml:"forks"`
	OpenIssues     float64 `yaml:"open_issues"`
	LastCommitDate float64 `yaml:"last_commit_date"`
	Archived       float64 `yaml:"archived"`
}

func NewParameterWeights() ParameterWeights {
	return ParameterWeights{
		Watchers:       defaultWatcherWeight,
		Stars:          defaultStarWeight,
		Forks:          defaultForkWeight,
		OpenIssues:     defaultOpenIssueWeight,
		LastCommitDate: defaultLastCommitDateWeight,
		Archived:       defaultArchivedWeight,
	}
}

// NewParameterWeightsFromReader creates ParameterWeights from an io.Reader
// This function replaces Viper with direct YAML parsing and proper error handling
func NewParameterWeightsFromReader(reader io.Reader) (ParameterWeights, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return ParameterWeights{}, fmt.Errorf("failed to read config data: %w", err)
	}

	var weights ParameterWeights

	err = yaml.Unmarshal(data, &weights)
	if err != nil {
		return ParameterWeights{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return weights, nil
}

// NewParameterWeightsFromFile creates ParameterWeights from a file path
// This function provides proper error handling without os.Exit
func NewParameterWeightsFromFile(configFilePath string) (ParameterWeights, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return ParameterWeights{}, fmt.Errorf("failed to read config file: %w", err)
	}
	defer file.Close()

	return NewParameterWeightsFromReader(file)
}
