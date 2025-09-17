package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uzumaki-inc/stay_or_go/parser"
)

func TestSelectParser_ReturnsCorrectTypes(t *testing.T) {
	t.Parallel()

	goParser, err := parser.SelectParser("go")
	require.NoError(t, err)

	if _, ok := goParser.(parser.GoParser); !ok {
		t.Fatalf("expected GoParser, got %T", goParser)
	}

	rubyParser, err := parser.SelectParser("ruby")
	require.NoError(t, err)

	if _, ok := rubyParser.(parser.RubyParser); !ok {
		t.Fatalf("expected RubyParser, got %T", rubyParser)
	}
}

func TestSelectParser_UnsupportedLanguage(t *testing.T) {
	t.Parallel()

	pythonParser, err := parser.SelectParser("python")
	assert.Nil(t, pythonParser)
	require.Error(t, err)
	assert.ErrorIs(t, err, parser.ErrUnsupportedLanguage)
}

func TestNewLibInfo_DefaultsAndOptions(t *testing.T) {
	t.Parallel()

	// Defaults
	li := parser.NewLibInfo("libX")
	assert.Equal(t, "libX", li.Name)
	assert.False(t, li.Skip)
	assert.Empty(t, li.SkipReason)
	assert.Nil(t, li.Others)
	assert.Empty(t, li.RepositoryURL)

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
