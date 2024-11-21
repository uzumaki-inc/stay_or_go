package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uzumaki-inc/StayOrGo/analyzer"
	"github.com/uzumaki-inc/StayOrGo/parser"
	"github.com/uzumaki-inc/StayOrGo/presenter"
	"github.com/uzumaki-inc/StayOrGo/utils"
)

// var greeting string
var (
	filePath       string
	outputFormat   string
	githubToken    string
	configFilePath string

	supportedLanguages = []string{"ruby", "go"}
	languageConfigMap  = map[string]string{
		"ruby": "Gemfile",
		"go":   "go.mod",
	}
	supportedOutputFormats = map[string]bool{
		"csv":      true,
		"tsv":      true,
		"markdown": true,
	}
)

// 引数を全部設定するlintを回避
//
//nolint:exhaustruct, lll
var rootCmd = &cobra.Command{
	Use:     "StayOrGo",
	Version: "0.1.0",
	Short:   "Analyze and score your Go and Ruby dependencies for popularity and maintenance",
	Long: `StayOrGo scans your Go (go.mod) and Ruby (Gemfile) dependency files to evaluate each library's popularity and maintenance status.
It generates scores to help you decide whether to keep (‘Stay’) or replace (‘Go’) your dependencies.
Output the results in Markdown, CSV, or TSV formats.`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Please Enter specify a language ("+
				strings.Join(supportedLanguages, " or ")+")")
			os.Exit(1)
		}

		language := args[0] // Get the language argument
		if !isSupportedLanguage(language) {
			utils.StdErrorPrintln("Error: Unsupported language: %s. Supported languages are: %s\n",
				language, strings.Join(supportedLanguages, ", "))
			os.Exit(1)
		}

		if filePath == "" {
			filePath = languageConfigMap[language]
		}

		if !supportedOutputFormats[outputFormat] {
			var keys []string
			for key := range supportedOutputFormats {
				keys = append(keys, key)
			}

			utils.StdErrorPrintln("Error: Unsupported output format: %s. Supported output formats are: %s\n",
				outputFormat, strings.Join(keys, ", "))
			os.Exit(1)
		}

		// --github-tokenのデータがなければ、環境変数　GITHUB_TOKENをチェックし
		// それがなければ--github-tokenかGITHUB_TOKENを追加するよう促す
		if githubToken == "" {
			githubToken = os.Getenv("GITHUB_TOKEN")
			if githubToken == "" {
				fmt.Fprintln(os.Stderr, `Please provide a GitHub token using the --github-token flag
			 or set the GITHUB_TOKEN environment variable`)
				os.Exit(1)
			}
		}

		utils.DebugPrintln("Selected Language: " + language)
		utils.DebugPrintln("Reading file: " + filePath)
		utils.DebugPrintln("Output format: " + outputFormat)

		var weights analyzer.ParameterWeights
		if configFilePath != "" {
			utils.DebugPrintln("Config file: " + configFilePath)
			weights = analyzer.NewParameterWeightsFromConfiFile(configFilePath)
		} else {
			weights = analyzer.NewParameterWeights()
		}
		analyzer := analyzer.NewGitHubRepoAnalyzer(githubToken, weights)

		utils.StdErrorPrintln("Selecting language... ")
		parser, err := parser.SelectParser(language)
		if err != nil {
			utils.StdErrorPrintln("Error selecting parser: %v", err)
			os.Exit(1)
		}
		utils.StdErrorPrintln("Parsing file...")
		libInfoList, err := parser.Parse(filePath)
		if err != nil {
			utils.StdErrorPrintln("Error parsing file: %v", err)
			os.Exit(1)
		}
		utils.StdErrorPrintln("Getting repository URLs...")
		parser.GetRepositoryURL(libInfoList)

		var repoURLs []string
		for _, info := range libInfoList {
			if !info.Skip {
				repoURLs = append(repoURLs, info.RepositoryURL)
			}
		}

		utils.StdErrorPrintln("Analyzing libraries with Github...")
		gitHubRepoInfos := analyzer.FetchGithubInfo(repoURLs)

		utils.StdErrorPrintln("Making dataset...")
		analyzedLibInfos := presenter.MakeAnalyzedLibInfoList(libInfoList, gitHubRepoInfos)
		presenter := presenter.SelectPresenter(outputFormat, analyzedLibInfos)

		utils.StdErrorPrintln("Displaying result...\n")
		presenter.Display()
	},
}

func isSupportedLanguage(language string) bool {
	for _, l := range supportedLanguages {
		if l == language {
			return true
		}
	}

	return false
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&filePath, "input", "i", "", "Specify the file to read")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "markdown", "Specify the output format (csv, tsv, markdown)")
	rootCmd.Flags().StringVarP(&githubToken, "github-token", "g", "", "GitHub token for authentication")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolVarP(&utils.Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "Modify evaluate parameters")
}
