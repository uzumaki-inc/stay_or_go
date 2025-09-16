package presenter

import (
    "testing"

    "github.com/uzumaki-inc/stay_or_go/analyzer"
    "github.com/uzumaki-inc/stay_or_go/parser"
)

func TestMakeAnalyzedLibInfoList_Mapping(t *testing.T) {
    t.Parallel()

    li1 := parser.LibInfo{Name: "a", RepositoryURL: "https://github.com/u/a"}
    li2 := parser.LibInfo{Name: "b", RepositoryURL: "https://github.com/u/b"}

    gi1 := analyzer.GitHubRepoInfo{GithubRepoURL: "https://github.com/u/a", Stars: 10}
    infos := MakeAnalyzedLibInfoList([]parser.LibInfo{li1, li2}, []analyzer.GitHubRepoInfo{gi1})

    if len(infos) != 2 { t.Fatalf("want 2") }
    if infos[0].GitHubRepoInfo == nil || infos[0].GitHubRepoInfo.Stars != 10 {
        t.Fatalf("first should be mapped with stars 10")
    }
    if infos[1].GitHubRepoInfo != nil {
        t.Fatalf("second should be nil GitHubRepoInfo")
    }
}

