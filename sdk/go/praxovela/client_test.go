package praxovela

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionRunLifecycle(t *testing.T) {
	var paths []string
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.Path)
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		bodies = append(bodies, body)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/sessions":
			_ = json.NewEncoder(w).Encode(Session{ID: "s1", WorkspaceID: "ws1"})
		case "/v1/sessions/s1/runs":
			_ = json.NewEncoder(w).Encode(Run{ID: "r1", SessionID: "s1", Status: "running"})
		case "/v1/runcenter/runs/r1":
			_ = json.NewEncoder(w).Encode(Run{ID: "r1", SessionID: "s1", Status: "completed"})
		case "/v1/runs/r1/cancel":
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c, err := New(Config{BaseURL: srv.URL, Token: "tok"})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	sess, err := c.CreateSession(ctx, "ws1")
	if err != nil {
		t.Fatal(err)
	}
	if sess.ID != "s1" {
		t.Fatalf("session = %+v", sess)
	}
	run, err := c.StartRun(ctx, "s1", "do work")
	if err != nil {
		t.Fatal(err)
	}
	if run.ID != "r1" {
		t.Fatalf("run = %+v", run)
	}
	if bodies[0]["workspace"] != "ws1" || bodies[1]["message"] != "do work" {
		t.Fatalf("bodies = %#v", bodies)
	}
	got, err := c.GetRun(ctx, "r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "completed" {
		t.Fatalf("run = %+v", got)
	}
	if err := c.CancelRun(ctx, "r1"); err != nil {
		t.Fatal(err)
	}
	want := []string{"POST /v1/sessions", "POST /v1/sessions/s1/runs", "GET /v1/runcenter/runs/r1", "POST /v1/runs/r1/cancel"}
	if len(paths) != len(want) {
		t.Fatalf("paths = %v", paths)
	}
	for i := range want {
		if paths[i] != want[i] {
			t.Fatalf("paths = %v, want %v", paths, want)
		}
	}
}

func TestNewRejectsInvalidConfig(t *testing.T) {
	if _, err := New(Config{}); err == nil {
		t.Fatal("expected error for empty BaseURL")
	}
	if _, err := New(Config{BaseURL: "ftp://x"}); err == nil {
		t.Fatal("expected error for non-http URL")
	}
}
