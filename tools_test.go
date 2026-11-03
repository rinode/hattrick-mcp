package main

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// expectedTools is the complete set of tools the server must register. It locks
// full CHPP coverage: adding or removing a tool without updating this list fails.
var expectedTools = strings.Fields(`
accept_challenge add_live_match are_challengeable challenge clear_live_matches
complete_login decline_challenge delete_finished_bids delete_live_match
get_achievements get_all_live_events get_alliance get_alliance_members
get_alliance_roles get_arena get_arena_stats get_arena_supporter_stats get_avatars
get_bookmarks get_challenges get_club get_country get_cup_matches get_current_bids
get_economy get_fans get_hof_players get_ladder_details get_ladders get_league
get_league_levels get_live_match get_live_matches get_manager get_match
get_match_lineup get_match_orders get_matches get_matches_archive get_my_alliances
get_national_team get_national_team_matches get_national_team_player_stats
get_national_team_players get_national_teams get_new_live_events get_old_coaches
get_old_players get_player get_player_events get_player_transfer_history get_players
get_region get_series get_series_fixtures get_staff get_supported_teams get_supporters
get_team get_teams get_tournament get_tournament_fixtures get_tournament_tables
get_tournaments get_training get_training_events get_training_stats get_transfers_team
get_translations get_world get_world_cup_groups get_world_cup_matches get_world_languages
get_youth_league get_youth_league_fixtures get_youth_matches get_youth_matches_archive
get_youth_player get_youth_players get_youth_team ignore_transfer is_challengeable
login place_bid predict_ratings search search_alliances search_transfers set_lineup
set_match_orders set_training set_youth_player_skills_visible
set_youth_team_skills_visible withdraw_challenge
`)

func listTools(t *testing.T) []*mcp.Tool {
	t.Helper()
	res, err := connect(t).ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	return res.Tools
}

func connect(t *testing.T) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()
	if _, err := newServer(&account{}, newCache()).Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func TestToolCoverage(t *testing.T) {
	got := map[string]bool{}
	for _, tl := range listTools(t) {
		got[tl.Name] = true
	}
	want := map[string]bool{}
	for _, n := range expectedTools {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		want[n] = true
		if !got[n] {
			t.Errorf("missing expected tool %q", n)
		}
	}
	for n := range got {
		if !want[n] {
			t.Errorf("unexpected tool %q (add it to expectedTools)", n)
		}
	}
	var names []string
	for n := range want {
		names = append(names, n)
	}
	sort.Strings(names)
	if len(got) != len(names) {
		t.Errorf("got %d tools, expected %d", len(got), len(names))
	}
}

// Every tool is read-only, a destructive write, or one of the two login tools.
func TestToolClassification(t *testing.T) {
	var reads, writes int
	for _, tool := range listTools(t) {
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
	if reads != 77 {
		t.Errorf("got %d read tools, want 77", reads)
	}
	if writes != 15 {
		t.Errorf("got %d write tools, want 15", writes)
	}
}

func TestEveryToolHasSchemaAndDescription(t *testing.T) {
	for _, tool := range listTools(t) {
		if strings.TrimSpace(tool.Description) == "" {
			t.Errorf("%s has no description", tool.Name)
		}
		if tool.InputSchema == nil {
			t.Errorf("%s has no input schema", tool.Name)
		}
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
