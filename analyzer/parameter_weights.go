package analyzer

import (
	"os"

	"github.com/konyu/StayOrGo/utils"
	"github.com/spf13/viper"
)

type ParameterWeights struct {
	Watchers         float64 `mapstructure:"watchers"`
	Stars            float64 `mapstructure:"stars"`
	Forks            float64 `mapstructure:"forks"`
	OpenPullRequests float64 `mapstructure:"open_pull_requests"`
	OpenIssues       float64 `mapstructure:"open_issues"`
	LastCommitDate   float64 `mapstructure:"last_commit_date"`
	Archived         float64 `mapstructure:"archived"`
}

func NewParameterWeights() ParameterWeights {
	return ParameterWeights{
		Watchers:         0.1,
		Stars:            0.1,
		Forks:            0.1,
		OpenPullRequests: 0.01,
		OpenIssues:       0.01,
		LastCommitDate:   -0.05,
		Archived:         -1000000,
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
