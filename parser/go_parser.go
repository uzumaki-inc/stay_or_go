package parser

import "github.com/konyu/StayOrGo/common"

type GoParser struct{}

func (p GoParser) Parse(file string) []common.LibInfo {
	result := make([]common.LibInfo, 0)
	return result
}

func (p GoParser) GetRepositoryUrl(libInfoList []common.LibInfo) []common.LibInfo {
	return libInfoList
}
