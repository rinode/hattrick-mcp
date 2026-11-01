package main

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/lucianoq/hattrick"
	"github.com/lucianoq/hattrick/chpp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// cacheTTL keeps repeated questions from hitting the CHPP servers.
const cacheTTL = 5 * time.Minute

const instructions = `Access to the authorised manager's Hattrick.org data via CHPP.
If a tool reports "not logged in", call login and then complete_login.
Omit an ID argument to get the manager's own team, series or match data.
Fields that are zero, false or empty are left out of results.
Skill values read "name (n)" on the 0 (non-existent) to 20 (divine) scale. Money is in SEK; divide by 10 for EUR display.
Tools whose name starts with set_, place_, unlock_, challenge, accept_, decline_, withdraw_, ignore_, delete_ or clear_ change the live account. They do nothing until called again with "confirm": true.`

type loginArgs struct {
	Code string `json:"code" jsonschema:"verification code Hattrick shows after granting access"`
}

type noArgs struct{}

func newServer(acct *account) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "hattrick-mcp", Version: version}, &mcp.ServerOptions{Instructions: instructions})
	c := &cache{m: map[string]cached{}}

	mcp.AddTool(s, &mcp.Tool{Name: "login", Description: "Start the Hattrick login. Returns a URL where the manager grants access and gets a verification code; then call complete_login with that code."},
		func(context.Context, *mcp.CallToolRequest, noArgs) (*mcp.CallToolResult, any, error) {
			url, err := acct.login()
			if err != nil {
				return nil, nil, err
			}
			return text("Ask the manager to open this URL, grant access and give you the verification code:\n" + url), nil, nil
		})
	mcp.AddTool(s, &mcp.Tool{Name: "complete_login", Description: "Finish the Hattrick login with the verification code from the login URL."},
		func(_ context.Context, _ *mcp.CallToolRequest, in loginArgs) (*mcp.CallToolResult, any, error) {
			if err := acct.completeLogin(in.Code); err != nil {
				return nil, nil, err
			}
			return text("Logged in."), nil, nil
		})

	registerReads(s, acct, c)
	registerWrites(s, acct)
	return s
}

// add registers a read-only tool whose result is returned as compact JSON text and cached.
func add[In any](s *mcp.Server, acct *account, c *cache, name, desc string, fn func(hattrick.Client, In) (any, error)) {
	t := &mcp.Tool{Name: name, Description: desc, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true}}
	mcp.AddTool(s, t, func(_ context.Context, _ *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
		ht, err := acct.client()
		if err != nil {
			return nil, nil, err
		}
		out, err := c.get(name, in, func() (string, error) {
			v, err := fn(ht, in)
			if err != nil {
				return "", err
			}
			b, err := json.Marshal(plain(v))
			return string(b), err
		})
		if err != nil {
			return nil, nil, err
		}
		return text(out), nil, nil
	})
}

// addWrite registers a tool that changes the live account. It is marked destructive
// and does nothing until the arguments carry "confirm": true.
func addWrite[In any](s *mcp.Server, acct *account, name, desc string, fn func(hattrick.Client, In) (any, error)) {
	destructive := true
	t := &mcp.Tool{Name: name, Description: desc + ` Set "confirm": true to apply; without it the tool only reports what it would do.`,
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructive}}
	mcp.AddTool(s, t, func(_ context.Context, _ *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
		if !confirmed(in) {
			return text("Not applied. This changes the live Hattrick account. Call again with \"confirm\": true to apply."), nil, nil
		}
		ht, err := acct.client()
		if err != nil {
			return nil, nil, err
		}
		v, err := fn(ht, in)
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(plain(v))
		if err != nil {
			return nil, nil, err
		}
		return text(string(b)), nil, nil
	})
}

// confirmField is embedded in every write tool's arguments.
type confirmField struct {
	Confirm bool `json:"confirm,omitempty" jsonschema:"must be true to apply this change to the live account"`
}

func confirmed(in any) bool {
	f := reflect.ValueOf(in).FieldByName("Confirm")
	return f.IsValid() && f.Kind() == reflect.Bool && f.Bool()
}

// srcOr defaults an unset source system to the main Hattrick system.
func srcOr(s chpp.SourceSystem) chpp.SourceSystem {
	if s == "" {
		return chpp.SourceSystemHattrickMainSystem
	}
	return s
}

func text(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: s}}}
}

type cached struct {
	text string
	at   time.Time
}

type cache struct {
	mu sync.Mutex
	m  map[string]cached
}

func (c *cache) get(name string, in any, load func() (string, error)) (string, error) {
	k, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	key := name + string(k)

	c.mu.Lock()
	e, ok := c.m[key]
	c.mu.Unlock()
	if ok && time.Since(e.at) < cacheTTL {
		return e.text, nil
	}

	txt, err := load()
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.m[key] = cached{txt, time.Now()}
	c.mu.Unlock()
	return txt, nil
}
