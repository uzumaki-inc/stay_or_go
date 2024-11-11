package presenter

import (
	"fmt"
	"os"
	"reflect"

	"github.com/konyu/StayOrGo/analyzer"
	"github.com/konyu/StayOrGo/parser"
	"github.com/konyu/StayOrGo/utils"
)

type AnalyzedLibInfo struct {
	LibInfo        *parser.LibInfo
	GitHubRepoInfo *analyzer.GitHubRepoInfo
}

func (ainfo AnalyzedLibInfo) Name() *string {
	if ainfo.LibInfo.Name != "" {
		return &ainfo.LibInfo.Name
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) RepositoryUrl() *string {
	if ainfo.LibInfo.RepositoryUrl != "" {
		return &ainfo.LibInfo.RepositoryUrl
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) Watchers() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Watchers
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) Stars() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Stars
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) Forks() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Forks
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) OpenPullRequests() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.OpenPullRequests
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) OpenIssues() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.OpenIssues
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) LastCommitDate() *string {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.LastCommitDate
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) GithubRepoUrl() *string {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.GithubRepoUrl
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) Archived() *bool {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Archived
	} else {
		return nil
	}
}

func (ainfo AnalyzedLibInfo) Score() *int {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.Score
	} else {
		return nil
	}
}
func (ainfo AnalyzedLibInfo) Skip() *bool {
	trueValue := true
	falseValue := false

	if ainfo.LibInfo.Skip == true {
		return &trueValue
	} else if ainfo.GitHubRepoInfo.Skip == true {
		return &trueValue
	}
	return &falseValue
}

func (ainfo AnalyzedLibInfo) SkipReason() *string {
	if ainfo.LibInfo.Skip == true {
		return &ainfo.LibInfo.SkipReason
	} else if ainfo.GitHubRepoInfo.Skip {
		return &ainfo.GitHubRepoInfo.SkipReason
	}
	return nil
}

func MakeAnalyzedLibInfoList(libInfoList []parser.LibInfo, gitHubRepoInfos []analyzer.GitHubRepoInfo) []AnalyzedLibInfo {
	var analyzedLibInfos []AnalyzedLibInfo
	var j = 0

	for _, info := range libInfoList {
		analyzedLibInfo := AnalyzedLibInfo{
			LibInfo: &info,
		}
		if j < len(gitHubRepoInfos) && info.RepositoryUrl == gitHubRepoInfos[j].GithubRepoUrl {
			analyzedLibInfo.GitHubRepoInfo = &gitHubRepoInfos[j]
			j++
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
		for i, header := range headerString {
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
				if i < len(headerString)-1 {
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
	"RepositoryUrl",
	"Watchers",
	"Stars",
	"Forks",
	"OpenPullRequests",
	"OpenIssues",
	"LastCommitDate",
	"Archived",
	"Score",
	"Skip",
	"SkipReason",
}

func SelectPresenter(format string, analyzedLibInfos []AnalyzedLibInfo) Presenter {
	var presenter Presenter
	switch format {
	case "tsv":
		presenter = TsvPresenter{analyzedLibInfos}
	case "csv":
		presenter = CsvPresenter{analyzedLibInfos}
	default:
		presenter = MarkdownPresenter{analyzedLibInfos}
	}
	return presenter
}
