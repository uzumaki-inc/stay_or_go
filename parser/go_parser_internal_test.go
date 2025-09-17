package parser

import (
	"errors"
	"testing"
)

func TestExtractRepoURL_Errors(t *testing.T) {
	t.Parallel()

	// invalid JSON
	_, err := extractRepoURL([]byte("not-json"), "github.com/user/lib")
	if err == nil || !errors.Is(err, ErrFailedToUnmarshalJSON) {
		t.Fatalf("expected ErrFailedToUnmarshalJSON, got %v", err)
	}

	// no github in name and empty origin.url
	body := []byte(`{"origin":{"url":""}}`)

	_, err = extractRepoURL(body, "code.gitea.io/sdk")
	if err == nil || !errors.Is(err, ErrNotAGitHubRepository) {
		t.Fatalf("expected ErrNotAGitHubRepository, got %v", err)
	}

	// non-github URL in origin
	body2 := []byte(`{"origin":{"url":"https://example.com/foo"}}`)

	_, err = extractRepoURL(body2, "example.com/foo")
	if err == nil || !errors.Is(err, ErrNotAGitHubRepository) {
		t.Fatalf("expected ErrNotAGitHubRepository, got %v", err)
	}
}
