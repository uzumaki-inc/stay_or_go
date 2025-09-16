package parser_test

import (
    "os"
    "testing"

    "github.com/jarcoal/httpmock"
    "github.com/stretchr/testify/assert"

    "github.com/uzumaki-inc/stay_or_go/parser"
)

func TestRubyParser_Parse_SkipsGemsInBlocks(t *testing.T) {

    content := `source "https://rubygems.org" do
  gem 'rails'
end

platforms :jruby do
  gem 'jruby-openssl'
end

install_if -> { true } do
  gem 'pg'
end

gem 'puma'
`

    f, err := os.CreateTemp("", "Gemfile-*.tmp")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(f.Name())
    if _, err := f.WriteString(content); err != nil {
        t.Fatal(err)
    }
    _ = f.Close()

    p := parser.RubyParser{}
    libs, err := p.Parse(f.Name())
    if err != nil {
        t.Fatal(err)
    }

    // Expect 4 gems listed in order encountered
    assert.Len(t, libs, 4)

    // First three are inside blocks → skipped
    assert.Equal(t, "rails", libs[0].Name)
    assert.True(t, libs[0].Skip)
    assert.Equal(t, "Not hosted on Github", libs[0].SkipReason)

    assert.Equal(t, "jruby-openssl", libs[1].Name)
    assert.True(t, libs[1].Skip)
    assert.Equal(t, "Not hosted on Github", libs[1].SkipReason)

    assert.Equal(t, "pg", libs[2].Name)
    assert.True(t, libs[2].Skip)
    assert.Equal(t, "Not hosted on Github", libs[2].SkipReason)

    // Outside blocks → not skipped
    assert.Equal(t, "puma", libs[3].Name)
    assert.False(t, libs[3].Skip)
}

func TestRubyParser_GetRepositoryURL_NonGitHubHomepageSkips(t *testing.T) {

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // homepage_uri points to non-GitHub → should skip
    httpmock.RegisterResponder(
        "GET",
        "https://rubygems.org/api/v1/gems/foo.json",
        httpmock.NewStringResponder(200, `{"homepage_uri": "https://example.com/foo", "source_code_uri": ""}`),
    )

    libs := []parser.LibInfo{{Name: "foo"}}

    p := parser.RubyParser{}
    updated := p.GetRepositoryURL(libs)

    assert.True(t, updated[0].Skip)
    assert.Equal(t, "Does not support libraries hosted outside of Github", updated[0].SkipReason)
    assert.Equal(t, "", updated[0].RepositoryURL)
}
