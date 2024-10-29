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

func (ainfo AnalyzedLibInfo) Skip() bool {
	if ainfo.LibInfo.Skip == true {
		return true
	} else if ainfo.GitHubRepoInfo.Skip {
		return true
	}
	return false
}

func (ainfo AnalyzedLibInfo) SkipReason() string {
	if ainfo.LibInfo.Skip == true {
		return ainfo.LibInfo.SkipReason
	} else if ainfo.GitHubRepoInfo.Skip {
		return ainfo.GitHubRepoInfo.SkipReason
	}
	return "No skip reason"
}

type Presenter interface {
	display([]AnalyzedLibInfo)
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
