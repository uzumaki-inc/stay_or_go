package analyzer

import (
	"os"

	"github.com/spf13/viper"

	"github.com/uzumaki-inc/stay_or_go/utils"
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
	Watchers       float64 `mapstructure:"watchers"`
	Stars          float64 `mapstructure:"stars"`
	Forks          float64 `mapstructure:"forks"`
	OpenIssues     float64 `mapstructure:"open_issues"`
	LastCommitDate float64 `mapstructure:"last_commit_date"`
	Archived       float64 `mapstructure:"archived"`
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

func NewParameterWeightsFromConfiFile(configFilePath string) ParameterWeights {
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		utils.StdErrorPrintln("Failed to read the configuration file: %v\n", err)
		os.Exit(1)
	}

	var weights ParameterWeights
	if err := viper.Unmarshal(&weights); err != nil {
		utils.StdErrorPrintln("Failed to unmarshal the configuration: %v\n", err)
		os.Exit(1)
	}

	return weights
}
