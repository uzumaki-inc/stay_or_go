package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var greeting string
var fileName string
var supportedLanguages = []string{"ruby", "go"}
var languageConfigMap = map[string]string{
	"ruby": "Gemfile",
	"go":   "go.mod",
}

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
			fmt.Println("Please Enter specify a language (" + strings.Join(supportedLanguages, " or ") + ")")
			os.Exit(1)
		}

		language := args[0] // Get the language argument
		if !isSupportedLanguage(language) {
			fmt.Println("Error: Unsupported language:", language)
			os.Exit(1)
		}

		if fileName == "" {
			fileName = languageConfigMap[language]
		}
		fmt.Println("Reading file:", fileName)

		fmt.Println(greeting, "World!")
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
	rootCmd.Flags().StringVarP(&fileName, "file", "f", "", "Specify the file to read")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&greeting, "greeting", "g", "Hello", "Greeting message to display before 'World!'")
}
