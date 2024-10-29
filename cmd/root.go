package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/konyu/StayOrGo/analyzer"
	"github.com/konyu/StayOrGo/parser"
	"github.com/konyu/StayOrGo/presenter"
	"github.com/spf13/cobra"
)

// var greeting string
var (
	filePath           string
	outputFormat       string
	githubToken        string
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

// type AnalyzedLibInfo struct {
// 	LibInfo        *parser.LibInfo
// 	GitHubRepoInfo *analyzer.GitHubRepoInfo
// }

// func (ainfo AnalyzedLibInfo) Skip() bool {
// 	if ainfo.LibInfo.Skip == true {
// 		return true
// 	} else if ainfo.GitHubRepoInfo.Skip {
// 		return true
// 	}
// 	return false
// }

// func (ainfo AnalyzedLibInfo) SkipReason() string {
// 	if ainfo.LibInfo.Skip == true {
// 		return ainfo.LibInfo.SkipReason
// 	} else if ainfo.GitHubRepoInfo.Skip {
// 		return ainfo.GitHubRepoInfo.SkipReason
// 	}
// 	return "No skip reason"
// }

// 引数を全部設定するlintを回避
//
//nolint:exhaustruct
var rootCmd = &cobra.Command{
	Use:     "StayOrGo",
	Version: "0.1.0",
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please Enter specify a language (" +
				strings.Join(supportedLanguages, " or ") + ")")
			os.Exit(1)
		}

		language := args[0] // Get the language argument
		if !isSupportedLanguage(language) {
			fmt.Printf("Error: Unsupported language: %s. Supported languages are: %s\n",
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
			fmt.Printf("Error: Unsupported output format: %s. Supported output formats are: %s\n",
				outputFormat, strings.Join(keys, ", "))
			os.Exit(1)
		}

		// --github-tokenのデータがなければ、環境変数　GITHUB_TOKENをチェックし
		// それがなければ--github-tokenかGITHUB_TOKENを追加するよう促す
		if githubToken == "" {
			githubToken = os.Getenv("GITHUB_TOKEN")
			if githubToken == "" {
				fmt.Println(`Please provide a GitHub token using the --github-token flag
			 or set the GITHUB_TOKEN environment variable`)
				os.Exit(1)
			}
		}

		fmt.Println("Language", language)
		fmt.Println("Reading file:", filePath)
		fmt.Println("Output format:", outputFormat)

		// TODO: パラメータをファイルから読み込めるようにする
		weights := analyzer.NewParameterWeights()
		analyzer := analyzer.NewGitHubRepoAnalyzer(githubToken, weights)

		parser := parser.SelectParser(language) // 言語に合わせたパーサーを選択
		libInfoList := parser.Parse(filePath)   // パーサーでファイルをパース

		parser.GetRepositoryURL(libInfoList)
		fmt.Println("GetRepositoryURL result:")
		for _, info := range libInfoList {
			fmt.Println(info)
		}

		var repoUrls []string
		for _, info := range libInfoList {
			if !info.Skip {
				repoUrls = append(repoUrls, info.RepositoryUrl)
			}
		}

		gitHubRepoInfos := analyzer.FetchGithubInfo(repoUrls)
		presenter.MakeAnalyzedLibInfoList(libInfoList, gitHubRepoInfos)
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
}
