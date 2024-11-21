package presenter

import (
	"fmt"
	"os"
	"reflect"

	"github.com/uzumaki-inc/StayOrGo/analyzer"
	"github.com/uzumaki-inc/StayOrGo/parser"
	"github.com/uzumaki-inc/StayOrGo/utils"
)

type AnalyzedLibInfo struct {
	LibInfo        *parser.LibInfo
	GitHubRepoInfo *analyzer.GitHubRepoInfo
}

func (ainfo AnalyzedLibInfo) Name() *string {
	if ainfo.LibInfo.Name != "" {
		return &ainfo.LibInfo.Name
	}

	return nil
}

func (ainfo AnalyzedLibInfo) RepositoryURL() *string {
	if ainfo.LibInfo.RepositoryURL != "" {
		return &ainfo.LibInfo.RepositoryURL
	}

	return nil
}

func (ainfo AnalyzedLibInfo) Watchers() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Watchers
	}

	return nil
}

func (ainfo AnalyzedLibInfo) Stars() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Stars
	}

	return nil
}

func (ainfo AnalyzedLibInfo) Forks() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Forks
	}

	return nil
}

func (ainfo AnalyzedLibInfo) OpenIssues() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.OpenIssues
	}

	return nil
}

func (ainfo AnalyzedLibInfo) LastCommitDate() *string {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.LastCommitDate
	}

	return nil
}

func (ainfo AnalyzedLibInfo) GithubRepoURL() *string {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.GithubRepoURL
	}

	return nil
}

func (ainfo AnalyzedLibInfo) Archived() *bool {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Archived
	}

	return nil
}

func (ainfo AnalyzedLibInfo) Score() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Score
	}

	return nil
}

func (ainfo AnalyzedLibInfo) Skip() *bool {
	trueValue := true
	falseValue := false

	if ainfo.LibInfo.Skip {
		return &trueValue
	} else if ainfo.GitHubRepoInfo.Skip {
		return &trueValue
	}

	return &falseValue
}

func (ainfo AnalyzedLibInfo) SkipReason() *string {
	if ainfo.LibInfo.Skip {
		return &ainfo.LibInfo.SkipReason
	} else if ainfo.GitHubRepoInfo.Skip {
		return &ainfo.GitHubRepoInfo.SkipReason
	}

	return nil
}

func MakeAnalyzedLibInfoList(
	libInfoList []parser.LibInfo,
	gitHubRepoInfos []analyzer.GitHubRepoInfo,
) []AnalyzedLibInfo {
	analyzedLibInfos := make([]AnalyzedLibInfo, 0, len(libInfoList))

	repoIndex := 0

	for _, info := range libInfoList {
		analyzedLibInfo := AnalyzedLibInfo{
			LibInfo:        &info,
			GitHubRepoInfo: nil,
		}

		if repoIndex < len(gitHubRepoInfos) && info.RepositoryURL == gitHubRepoInfos[repoIndex].GithubRepoURL {
			analyzedLibInfo.GitHubRepoInfo = &gitHubRepoInfos[repoIndex]
			repoIndex++
		}

		analyzedLibInfos = append(analyzedLibInfos, analyzedLibInfo)
	}

	return analyzedLibInfos
}

type Presenter interface {
	Display()
	makeHeader() []string
	makeBody() []string
}

func Display(p Presenter) {
	header := p.makeHeader()
	body := p.makeBody()

	for _, line := range header {
		fmt.Println(line)
	}

	for _, line := range body {
		fmt.Println(line)
	}
}

func makeBody(analyzedLibInfos []AnalyzedLibInfo, separator string) []string {
	rows := []string{}

	for _, info := range analyzedLibInfos {
		row := ""
		val := reflect.ValueOf(info)

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		for index, header := range headerString {
			method := val.MethodByName(header)

			if method.IsValid() {
				result := method.Call(nil)

				var resultStr interface{}

				if len(result) > 0 && result[0].IsValid() && !result[0].IsNil() {
					resultStr = result[0].Elem().Interface()
				} else {
					resultStr = "N/A"
				}

				row += fmt.Sprintf("%v", resultStr)
				// 最後の要素でない場合にのみseparatorを追加
				if index < len(headerString)-1 {
					row += separator
				}
			} else {
				utils.StdErrorPrintln("method %s not found in %v", header, info)
				os.Exit(1)
			}
		}

		if separator == "|" {
			row = "|" + row + "|"
		}

		rows = append(rows, row)
	}

	return rows
}

var headerString = []string{
	"Name",
	"RepositoryURL",
	"Watchers",
	"Stars",
	"Forks",
	"OpenIssues",
	"LastCommitDate",
	"Archived",
	"Score",
	"Skip",
	"SkipReason",
}

func SelectPresenter(format string, analyzedLibInfos []AnalyzedLibInfo) Presenter {
	switch format {
	case "tsv":
		return TsvPresenter{analyzedLibInfos}
	case "csv":
		return CsvPresenter{analyzedLibInfos}
	default:
		return MarkdownPresenter{analyzedLibInfos}
	}
}
