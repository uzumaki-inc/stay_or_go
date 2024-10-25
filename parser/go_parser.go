package parser

import (
	"fmt"

	"github.com/konyu/StayOrGo/common"
)

type GoParser struct{}

func (p GoParser) Parse(file string) []common.LibInfo {
	fmt.Println(file)

	result := make([]common.LibInfo, 0)

	return result
}

func (p GoParser) GetRepositoryURL(libInfoList []common.LibInfo) []common.LibInfo {
	return libInfoList
}
