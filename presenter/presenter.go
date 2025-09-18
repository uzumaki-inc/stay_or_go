package presenter

import (
	"fmt"
	"reflect"

	"github.com/uzumaki-inc/stay_or_go/analyzer"
	"github.com/uzumaki-inc/stay_or_go/parser"
	"github.com/uzumaki-inc/stay_or_go/utils"
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

	for idx := range libInfoList {
		analyzedLibInfo := AnalyzedLibInfo{
			LibInfo:        &libInfoList[idx],
			GitHubRepoInfo: nil,
		}

		if repoIndex < len(gitHubRepoInfos) && libInfoList[idx].RepositoryURL == gitHubRepoInfos[repoIndex].GithubRepoURL {
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
		row := makeRow(info, separator)
		rows = append(rows, row)
	}

	return rows
}

func makeRow(info AnalyzedLibInfo, separator string) string {
	row := ""
	val := reflect.ValueOf(info)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for index, header := range headerString {
		cellValue := getCellValue(val, header, info)
		row += cellValue

		// 最後の要素でない場合にのみseparatorを追加
		if index < len(headerString)-1 {
			row += separator
		}
	}

	if separator == "|" {
		row = "|" + row + "|"
	}

	return row
}

func getCellValue(val reflect.Value, header string, info AnalyzedLibInfo) string {
	method := val.MethodByName(header)

	if !method.IsValid() {
		utils.StdErrorPrintln("method %s not found in %v", header, info)

		return "N/A"
	}

	result := method.Call(nil)

	if len(result) == 0 || !result[0].IsValid() || result[0].IsNil() {
		return "N/A"
	}

	return fmt.Sprintf("%v", result[0].Elem().Interface())
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
