package cmd

import (
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
    "github.com/uzumaki-inc/stay_or_go/analyzer"
    "github.com/uzumaki-inc/stay_or_go/parser"
    "github.com/uzumaki-inc/stay_or_go/presenter"
    "github.com/uzumaki-inc/stay_or_go/utils"
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

// AnalyzerPort is a minimal adapter for analyzer used by cmd to enable testing with stubs.
type AnalyzerPort interface {
    FetchGithubInfo(repositoryUrls []string) []analyzer.GitHubRepoInfo
}

// PresenterPort narrows the presenter to only what's used here.
type PresenterPort interface {
    Display()
}

// Deps bundles injectable constructors/selectors for testability.
type Deps struct {
    NewAnalyzer     func(token string, weights analyzer.ParameterWeights) AnalyzerPort
    SelectParser    func(language string) (parser.Parser, error)
    SelectPresenter func(format string, analyzedLibInfos []presenter.AnalyzedLibInfo) PresenterPort
}

var defaultDeps = Deps{
    NewAnalyzer: func(token string, weights analyzer.ParameterWeights) AnalyzerPort {
        return analyzer.NewGitHubRepoAnalyzer(token, weights)
    },
    SelectParser:    parser.SelectParser,
    SelectPresenter: func(format string, analyzedLibInfos []presenter.AnalyzedLibInfo) PresenterPort {
        return presenter.SelectPresenter(format, analyzedLibInfos)
    },
}

// 引数を全部設定するlintを回避
//
//nolint:exhaustruct, lll
var rootCmd = &cobra.Command{
	Use:     "stay_or_go",
	Version: "0.1.2",
	Short:   "Analyze and score your Go and Ruby dependencies for popularity and maintenance",
	Long: `stay_or_go scans your Go (go.mod) and Ruby (Gemfile) dependency files to evaluate each library's popularity and maintenance status.
It generates scores to help you decide whether to keep (‘Stay’) or replace (‘Go’) your dependencies.
Output the results in Markdown, CSV, or TSV formats.`,
    Run: func(_ *cobra.Command, args []string) {
        if len(args) == 0 {
            fmt.Fprintln(os.Stderr, "Please Enter specify a language ("+
                strings.Join(supportedLanguages, " or ")+")")
            os.Exit(1)
        }

        language := args[0]
        // Delegate to testable runner
        if err := run(language, filePath, outputFormat, githubToken, configFilePath, utils.Verbose, defaultDeps); err != nil {
            os.Exit(1)
        }
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

// run executes the core logic with injectable dependencies. Returns error instead of exiting.
func run(language, inFile, format, token, config string, verbose bool, deps Deps) error {
    if !isSupportedLanguage(language) {
        utils.StdErrorPrintln("Error: Unsupported language: %s. Supported languages are: %s\n",
            language, strings.Join(supportedLanguages, ", "))
        return fmt.Errorf("unsupported language: %s", language)
    }

    file := inFile
    if file == "" {
        file = languageConfigMap[language]
    }

    if !supportedOutputFormats[format] {
        var keys []string
        for key := range supportedOutputFormats {
            keys = append(keys, key)
        }
        utils.StdErrorPrintln("Error: Unsupported output format: %s. Supported output formats are: %s\n",
            format, strings.Join(keys, ", "))
        return fmt.Errorf("unsupported format: %s", format)
    }

    if token == "" {
        token = os.Getenv("GITHUB_TOKEN")
        if token == "" {
            fmt.Fprintln(os.Stderr, "Please provide a GitHub token using the --github-token flag or set the GITHUB_TOKEN environment variable")
            return fmt.Errorf("missing github token")
        }
    }

    utils.DebugPrintln("Selected Language: " + language)
    utils.DebugPrintln("Reading file: " + file)
    utils.DebugPrintln("Output format: " + format)

    var weights analyzer.ParameterWeights
    if config != "" {
        utils.DebugPrintln("Config file: " + config)
        weights = analyzer.NewParameterWeightsFromConfiFile(config)
    } else {
        weights = analyzer.NewParameterWeights()
    }
    az := deps.NewAnalyzer(token, weights)

    utils.StdErrorPrintln("Selecting language... ")
    pr, err := deps.SelectParser(language)
    if err != nil {
        utils.StdErrorPrintln("Error selecting parser: %v", err)
        return err
    }
    utils.StdErrorPrintln("Parsing file...")
    libInfoList, err := pr.Parse(file)
    if err != nil {
        utils.StdErrorPrintln("Error parsing file: %v", err)
        return err
    }
    utils.StdErrorPrintln("Getting repository URLs...")
    pr.GetRepositoryURL(libInfoList)

    var repoURLs []string
    for _, info := range libInfoList {
        if !info.Skip {
            repoURLs = append(repoURLs, info.RepositoryURL)
        }
    }

    utils.StdErrorPrintln("Analyzing libraries with Github...")
    var gitHubRepoInfos []analyzer.GitHubRepoInfo
    if len(repoURLs) > 0 {
        gitHubRepoInfos = az.FetchGithubInfo(repoURLs)
    } else {
        gitHubRepoInfos = []analyzer.GitHubRepoInfo{}
    }

    utils.StdErrorPrintln("Making dataset...")
    analyzedLibInfos := presenter.MakeAnalyzedLibInfoList(libInfoList, gitHubRepoInfos)
    pz := deps.SelectPresenter(format, analyzedLibInfos)

    utils.StdErrorPrintln("Displaying result...\n")
    pz.Display()

    return nil
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
