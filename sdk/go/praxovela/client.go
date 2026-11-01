// Package praxovela is the public Go client for the PRAXOVELA / AXON Core HTTP
// API. It wraps sessions, agent runs, run-center queries, cancellation, and
// recovery behind typed methods so callers do not hand-roll HTTP.
package praxovela

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type Client struct {
	base   string
	token  string
	client *http.Client
}

type Session struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type Run struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id,omitempty"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Category    string    `json:"category,omitempty"`
	Status      string    `json:"status,omitempty"`
	ActorID     string    `json:"actor_id,omitempty"`
	NeedsReview bool      `json:"needs_review,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type Event struct {
	Type        string          `json:"type"`
	TenantID    string          `json:"tenant_id,omitempty"`
	WorkspaceID string          `json:"workspace_id,omitempty"`
	SessionID   string          `json:"session_id,omitempty"`
	RunID       string          `json:"run_id,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
	Timestamp   time.Time       `json:"timestamp,omitempty"`
}

type Recovery struct {
	CheckpointID   string `json:"checkpoint_id,omitempty"`
	SessionID      string `json:"session_id,omitempty"`
	RunID          string `json:"run_id,omitempty"`
	Restored       bool   `json:"restored"`
	MessageCount   int    `json:"message_count,omitempty"`
	CursorVerified bool   `json:"cursor_verified,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return "praxovela: unexpected status " + http.StatusText(e.StatusCode)
}

func New(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, errors.New("praxovela: BaseURL is required")
	}
	u, err := url.Parse(cfg.BaseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, errors.New("praxovela: BaseURL must be http or https")
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{base: strings.TrimRight(cfg.BaseURL, "/"), token: cfg.Token, client: cfg.HTTPClient}, nil
}

func (c *Client) Health(ctx context.Context) error {
	return c.do(ctx, http.MethodGet, "/health", nil, nil)
}

func (c *Client) CreateSession(ctx context.Context, workspace string) (*Session, error) {
	body := map[string]any{}
	if workspace != "" {
		body["workspace"] = workspace
	}
	var out Session
	if err := c.do(ctx, http.MethodPost, "/v1/sessions", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetSession(ctx context.Context, id string) (*Session, error) {
	var out Session
	if err := c.do(ctx, http.MethodGet, "/v1/sessions/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListSessions(ctx context.Context) ([]Session, error) {
	var out []Session
	if err := c.do(ctx, http.MethodGet, "/v1/sessions", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) StartRun(ctx context.Context, sessionID, message string) (*Run, error) {
	var out Run
	if err := c.do(ctx, http.MethodPost, "/v1/sessions/"+url.PathEscape(sessionID)+"/runs",
		map[string]any{"message": message}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetRun(ctx context.Context, runID string) (*Run, error) {
	var out Run
	if err := c.do(ctx, http.MethodGet, "/v1/runcenter/runs/"+url.PathEscape(runID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListRuns(ctx context.Context) ([]Run, error) {
	var out struct {
		Runs []Run `json:"runs"`
	}
	if err := c.do(ctx, http.MethodGet, "/v1/runcenter/runs", nil, &out); err != nil {
		return nil, err
	}
	return out.Runs, nil
}

func (c *Client) ListRunEvents(ctx context.Context, runID string) ([]Event, error) {
	var out []Event
	if err := c.do(ctx, http.MethodGet, "/v1/runcenter/runs/"+url.PathEscape(runID)+"/events", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CancelRun(ctx context.Context, runID string) error {
	return c.do(ctx, http.MethodPost, "/v1/runs/"+url.PathEscape(runID)+"/cancel", map[string]any{}, nil)
}

func (c *Client) RecoverRun(ctx context.Context, runID string) (*Recovery, error) {
	var out Recovery
	if err := c.do(ctx, http.MethodPost, "/v1/runs/"+url.PathEscape(runID)+"/recover", map[string]any{}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) do(ctx context.Context, method, path string, payload, out any) error {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, body)
	if err != nil {
		return err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return &HTTPError{StatusCode: resp.StatusCode, Body: string(snippet)}
	}
	if out == nil {
		return nil
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}
