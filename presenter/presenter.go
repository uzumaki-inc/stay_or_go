package presenter

import (
	"fmt"

	"github.com/konyu/StayOrGo/analyzer"
	"github.com/konyu/StayOrGo/parser"
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

func (ainfo AnalyzedLibInfo) RepositoryName() *string {
	if ainfo.GitHubRepoInfo != nil && ainfo.GitHubRepoInfo.RepositoryName != "" {
		return &ainfo.GitHubRepoInfo.RepositoryName
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

func (ainfo AnalyzedLibInfo) LibraryName() *string {
	if ainfo.GitHubRepoInfo != nil {
		return &ainfo.GitHubRepoInfo.LibraryName
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

	for i, info := range libInfoList {
		analyzedLibInfo := AnalyzedLibInfo{
			LibInfo: &info,
		}
		if i <= len(gitHubRepoInfos) && info.RepositoryUrl == gitHubRepoInfos[j].GithubRepoUrl {
			analyzedLibInfo.GitHubRepoInfo = &gitHubRepoInfos[j]
			j++
		}
		analyzedLibInfos = append(analyzedLibInfos, analyzedLibInfo)
	}

	for _, info := range analyzedLibInfos {
		if info.GitHubRepoInfo != nil {
			fmt.Printf("Repo: %s, Stars: %d, Forks: %d, Last Commit: %s, Archived: %t, Score: %d \n",
				info.GitHubRepoInfo.RepositoryName, info.GitHubRepoInfo.Stars, info.GitHubRepoInfo.Forks, info.GitHubRepoInfo.LastCommitDate, info.GitHubRepoInfo.Archived, info.GitHubRepoInfo.Score)
		}
	}
	return analyzedLibInfos
}

type Presenter interface {
	Display()
	makeHeader() []string
}

var headerString = []string{
	"Name",
	"RepositoryUrl",
	// "RepositoryName",
	"Watchers",
	"Stars",
	"Forks",
	"OpenPullRequests",
	"OpenIssues",
	"LastCommitDate",
	"LibraryName",
	"GithubRepoUrl",
	"Archived",
	"Score",
	"Skip",
	"SkipReason",
}

func SelectPresenter(format string, analyzedLibInfos []AnalyzedLibInfo) Presenter {
	var presenter Presenter
	switch format {
	case "tsv":
		// presenter = TsvPresenter{analyzedLibInfos}
	case "csv":
		// presenter = CsvPresenter{analyzedLibInfos}
	default:
		presenter = MarkdownPresenter{analyzedLibInfos}
	}
	return presenter
}
