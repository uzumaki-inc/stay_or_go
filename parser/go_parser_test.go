package parser_test

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/uzumaki-inc/stay_or_go/parser"
)

//nolint:funlen // Complex test logic requires many assertions
func TestGoParser_Parse_RequireReplaceAndIndirect(t *testing.T) {
	t.Parallel()

	content := `module example.com/demo

require (
    github.com/user/libone v1.2.3
    golang.org/x/sys v0.1.0 // indirect
    github.com/user/libtwo v0.9.0
    code.gitea.io/sdk v1.0.0
    github.com/replaced/mod v1.0.0
)

replace (
    github.com/replaced/mod v1.0.0 => ./local/mod
)
`

	tmpFile, err := os.CreateTemp(t.TempDir(), "go.mod-*.tmp")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	_ = tmpFile.Close()

	p := parser.GoParser{}

	libs, err := p.Parse(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Expect: libone, libtwo, sdk, mod(replaced) => 4 entries (indirect excluded)
	assert.Len(t, libs, 4)

	// helper to find by name
	byName := func(name string) *parser.LibInfo {
		for i := range libs {
			if libs[i].Name == name {
				return &libs[i]
			}
		}

		return nil
	}

	libone := byName("libone")
	libtwo := byName("libtwo")
	sdk := byName("sdk")
	replaced := byName("mod")

	if libone == nil || libtwo == nil || sdk == nil || replaced == nil {
		t.Fatalf("expected all libs to be found, got: %+v", libs)
	}

	assert.False(t, libone.Skip)
	assert.Equal(t, []string{"github.com/user/libone", "v1.2.3"}, libone.Others)

	assert.False(t, libtwo.Skip)
	assert.Equal(t, []string{"github.com/user/libtwo", "v0.9.0"}, libtwo.Others)

	assert.False(t, sdk.Skip)
	assert.Equal(t, []string{"code.gitea.io/sdk", "v1.0.0"}, sdk.Others)

	assert.True(t, replaced.Skip)
	assert.Equal(t, "replaced module", replaced.SkipReason)
}

//nolint:paralleltest,funlen // Uses httpmock which doesn't support parallel tests, complex setup
func TestGoParser_GetRepositoryURL_SetsURLAndSkips(t *testing.T) {
	// Prepare initial lib list as if parsed
	libs := []parser.LibInfo{
		parser.NewLibInfo("libone", parser.WithOthers([]string{"github.com/user/libone", "v1.2.3"})),
		parser.NewLibInfo("libtwo", parser.WithOthers([]string{"github.com/user/libtwo", "v0.9.0"})),
		parser.NewLibInfo("sdk", parser.WithOthers([]string{"code.gitea.io/sdk", "v1.0.0"})),
		parser.NewLibInfo("mod", parser.WithSkip(true), parser.WithSkipReason("replaced module")),
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Success: libone -> explicit GitHub URL in origin
	httpmock.RegisterResponder(
		"GET",
		"https://proxy.golang.org/github.com/user/libone/@v/v1.2.3.info",
		httpmock.NewStringResponder(200,
			`{"version":"v1.2.3","time":"2024-01-01T00:00:00Z",`+
				`"origin":{"vcs":"git","url":"https://github.com/user/libone","ref":"main","hash":"deadbeef"}}`),
	)

	// Success via fallback: libtwo -> origin.url empty but module path contains github.com
	httpmock.RegisterResponder(
		"GET",
		"https://proxy.golang.org/github.com/user/libtwo/@v/v0.9.0.info",
		httpmock.NewStringResponder(200,
			`{"version":"v0.9.0","time":"2024-01-02T00:00:00Z",`+
				`"origin":{"vcs":"git","url":"","ref":"main","hash":"deadbeef"}}`),
	)

	// Non-GitHub or error -> mark skip
	httpmock.RegisterResponder(
		"GET",
		"https://proxy.golang.org/code.gitea.io/sdk/@v/v1.0.0.info",
		httpmock.NewStringResponder(404, `not found`),
	)

	p := parser.GoParser{}
	updated := p.GetRepositoryURL(libs)

	// Map by name for assertions
	get := func(name string) *parser.LibInfo {
		for i := range updated {
			if updated[i].Name == name {
				return &updated[i]
			}
		}

		return nil
	}

	libone := get("libone")
	libtwo := get("libtwo")
	sdk := get("sdk")
	replaced := get("mod")

	if libone == nil || libtwo == nil || sdk == nil || replaced == nil {
		t.Fatalf("expected all libs to be found, got: %+v", updated)
	}

	assert.Equal(t, "https://github.com/user/libone", libone.RepositoryURL)
	assert.Equal(t, "https://github.com/user/libtwo", libtwo.RepositoryURL)

	assert.True(t, sdk.Skip)
	assert.Equal(t, "Does not support libraries hosted outside of Github", sdk.SkipReason)
	assert.Empty(t, sdk.RepositoryURL)

	// replaced item should remain skipped and untouched
	assert.True(t, replaced.Skip)
	assert.Equal(t, "replaced module", replaced.SkipReason)
}
