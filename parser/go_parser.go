package parser

import (
	"fmt"
)

type GoParser struct{}

func (p GoParser) Parse(file string) []LibInfo {
	fmt.Println(file)

	result := make([]LibInfo, 0)

	return result
}

func (p GoParser) GetRepositoryURL(libInfoList []LibInfo) []LibInfo {
	return libInfoList
}
