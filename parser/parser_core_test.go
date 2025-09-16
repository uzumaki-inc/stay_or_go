package parser_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzumaki-inc/stay_or_go/parser"
)

func TestSelectParser_ReturnsCorrectTypes(t *testing.T) {
	t.Parallel()

	p, err := parser.SelectParser("go")
	assert.NoError(t, err)
	if _, ok := p.(parser.GoParser); !ok {
		t.Fatalf("expected GoParser, got %T", p)
	}

	p, err = parser.SelectParser("ruby")
	assert.NoError(t, err)
	if _, ok := p.(parser.RubyParser); !ok {
		t.Fatalf("expected RubyParser, got %T", p)
	}
}

func TestSelectParser_UnsupportedLanguage(t *testing.T) {
	t.Parallel()

	p, err := parser.SelectParser("python")
	assert.Nil(t, p)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, parser.ErrUnsupportedLanguage))
}

func TestNewLibInfo_DefaultsAndOptions(t *testing.T) {
	t.Parallel()

	// Defaults
	li := parser.NewLibInfo("libX")
	assert.Equal(t, "libX", li.Name)
	assert.False(t, li.Skip)
	assert.Equal(t, "", li.SkipReason)
	assert.Nil(t, li.Others)
	assert.Equal(t, "", li.RepositoryURL)

	// With options
	li2 := parser.NewLibInfo(
		"libY",
		parser.WithSkip(true),
		parser.WithSkipReason("reason"),
		parser.WithOthers([]string{"a", "b"}),
	)
	assert.Equal(t, "libY", li2.Name)
	assert.True(t, li2.Skip)
	assert.Equal(t, "reason", li2.SkipReason)
	assert.Equal(t, []string{"a", "b"}, li2.Others)
}
