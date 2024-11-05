package parser

type GoParser struct{}

func (p GoParser) Parse(file string) []LibInfo {
	result := make([]LibInfo, 0)

	return result
}

func (p GoParser) GetRepositoryURL(libInfoList []LibInfo) []LibInfo {
	return libInfoList
}
