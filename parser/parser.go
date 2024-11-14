package parser

import (
	"errors"
	"fmt"
)

var (
	ErrMethodNotFound           = errors.New("method not found in struct")
	ErrFailedToReadFile         = errors.New("failed to read file")
	ErrFailedToResetFilePointer = errors.New("failed to reset file pointer")
	ErrFailedToScanFile         = errors.New("failed to scan file")
	ErrFailedToGetRepository    = errors.New("can't get the gem repository, skipping")
	ErrNotAGitHubRepository     = errors.New("not a GitHub repository, skipping")
	ErrFailedToReadResponseBody = errors.New("failed to read response body")
	ErrFailedToUnmarshalJSON    = errors.New("failed to unmarshal JSON response")
	ErrInvalidLineFormat        = errors.New("invalid line format")
	ErrMissingGemName           = errors.New("missing gem name")
	ErrUnsupportedLanguage      = errors.New("unsupported language")
)

const timeOutSec = 30

type LibInfo struct {
	Skip          bool     // スキップするかどうかのフラグ
	SkipReason    string   // スキップ理由
	Name          string   // ライブラリの名前
	Others        []string // その他のライブラリの設定値
	RepositoryURL string   // githubのりポトリのURL
}

type LibInfoOption func(*LibInfo)

func WithSkip(skip bool) LibInfoOption {
	return func(l *LibInfo) {
		l.Skip = skip
	}
}

func WithSkipReason(reason string) LibInfoOption {
	return func(l *LibInfo) {
		l.SkipReason = reason
	}
}

func WithOthers(others []string) LibInfoOption {
	return func(l *LibInfo) {
		l.Others = others
	}
}

func NewLibInfo(name string, options ...LibInfoOption) LibInfo {
	libInfo := LibInfo{
		Name:          name,
		Skip:          false,
		SkipReason:    "",
		Others:        nil,
		RepositoryURL: "",
	}

	for _, option := range options {
		option(&libInfo)
	}

	return libInfo
}

type Parser interface {
	Parse(file string) ([]LibInfo, error)
	GetRepositoryURL(AnalyzedLibInfoList []LibInfo) []LibInfo
}

func SelectParser(language string) (Parser, error) {
	switch language {
	case "ruby":
		return RubyParser{}, nil
	case "go":
		return GoParser{}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedLanguage, language)
	}
}
