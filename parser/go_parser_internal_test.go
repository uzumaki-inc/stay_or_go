package parser

import (
    "testing"
)

func TestExtractRepoURL_Errors(t *testing.T) {
    t.Parallel()

    // invalid JSON
    if _, err := extractRepoURL([]byte("not-json"), "github.com/user/lib"); err == nil || err != ErrFailedToUnmarshalJSON {
        t.Fatalf("expected ErrFailedToUnmarshalJSON, got %v", err)
    }

    // no github in name and empty origin.url
    body := []byte(`{"origin":{"url":""}}`)
    if _, err := extractRepoURL(body, "code.gitea.io/sdk"); err == nil || err != ErrNotAGitHubRepository {
        t.Fatalf("expected ErrNotAGitHubRepository, got %v", err)
    }

    // non-github URL in origin
    body2 := []byte(`{"origin":{"url":"https://example.com/foo"}}`)
    if _, err := extractRepoURL(body2, "example.com/foo"); err == nil || err != ErrNotAGitHubRepository {
        t.Fatalf("expected ErrNotAGitHubRepository, got %v", err)
    }
}

