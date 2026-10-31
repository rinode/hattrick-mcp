package main

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func connect(t *testing.T) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()
	if _, err := newServer(&account{}).Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func TestToolsList(t *testing.T) {
	res, err := connect(t).ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Tools) != 16 {
		t.Errorf("got %d tools, want 16", len(res.Tools))
	}
	for _, tool := range res.Tools {
		readOnly := tool.Annotations != nil && tool.Annotations.ReadOnlyHint
		if login := strings.Contains(tool.Name, "login"); readOnly == login {
			t.Errorf("%s: read-only is %v", tool.Name, readOnly)
		}
	}
}

func TestNotLoggedIn(t *testing.T) {
	cs := connect(t)
	for name, want := range map[string]string{
		"get_team":       "not logged in",
		"complete_login": "no login in progress",
	} {
		res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: map[string]any{"code": "x"}})
		if err != nil {
			t.Fatal(err)
		}
		got := res.Content[0].(*mcp.TextContent).Text
		if !res.IsError || !strings.Contains(got, want) {
			t.Errorf("%s: got %q, want an error containing %q", name, got, want)
		}
	}
}
