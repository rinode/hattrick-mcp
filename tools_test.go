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

// Every tool is either read-only, a destructive write, or one of the two login tools.
func TestToolAnnotations(t *testing.T) {
	res, err := connect(t).ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	var reads, writes int
	for _, tool := range res.Tools {
		a := tool.Annotations
		switch {
		case tool.Name == "login" || tool.Name == "complete_login":
			if a != nil && a.ReadOnlyHint {
				t.Errorf("%s should not be read-only", tool.Name)
			}
		case a != nil && a.ReadOnlyHint:
			reads++
		case a != nil && a.DestructiveHint != nil && *a.DestructiveHint:
			writes++
		default:
			t.Errorf("%s is neither read-only nor destructive", tool.Name)
		}
	}
	if reads < 70 {
		t.Errorf("got %d read tools, want at least 70", reads)
	}
	if writes != 14 {
		t.Errorf("got %d write tools, want 14", writes)
	}
}

func TestNotLoggedIn(t *testing.T) {
	res := call(t, "get_team", map[string]any{})
	if !res.IsError || !strings.Contains(textOf(res), "not logged in") {
		t.Errorf("get_team: got %q", textOf(res))
	}
}

// A write tool does nothing until confirm is true, and the gate runs before the login check.
func TestWriteNeedsConfirm(t *testing.T) {
	res := call(t, "clear_live_matches", map[string]any{})
	if res.IsError || !strings.Contains(textOf(res), "Not applied") {
		t.Errorf("without confirm: got %q", textOf(res))
	}
	res = call(t, "clear_live_matches", map[string]any{"confirm": true})
	if !res.IsError || !strings.Contains(textOf(res), "not logged in") {
		t.Errorf("with confirm: got %q", textOf(res))
	}
}

func call(t *testing.T, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	res, err := connect(t).CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func textOf(res *mcp.CallToolResult) string {
	return res.Content[0].(*mcp.TextContent).Text
}
