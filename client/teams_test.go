package client

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gomicro/scribe"
	"github.com/google/go-github/v92/github"
)

func TestRemoveTeamMemberBySlug(t *testing.T) {
	var method, path string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	gh := github.NewClient(nil)
	gh.BaseURL, _ = url.Parse(srv.URL + "/")

	c := &Client{ghClient: gh}

	var out bytes.Buffer
	scrb := scribe.NewScribe(&out, &scribe.Theme{Describe: scribe.NoopDecorator, Print: scribe.NoopDecorator})

	c.RemoveTeamMemberBySlug(context.Background(), scrb, "my-org", "backend", "alice")

	if method != "" {
		t.Fatal("expected removal to be deferred until apply")
	}

	if len(c.stack) != 1 {
		t.Fatalf("expected 1 queued action, got %d", len(c.stack))
	}

	err := c.stack[0]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if method != http.MethodDelete {
		t.Errorf("method: got %s, want DELETE", method)
	}

	if path != "/orgs/my-org/teams/backend/memberships/alice" {
		t.Errorf("path: got %s", path)
	}
}
