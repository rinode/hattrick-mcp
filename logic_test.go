package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/lucianoq/hattrick/chpp"
	"github.com/lucianoq/hattrick/chpp/id"
)

// buildLineup must emit the exact JSON shape the live CHPP setmatchorder
// endpoint accepts: numeric ids/enums, newLineup and orderType as strings,
// 14 positions, 14 bench slots, 11 kickers.
func TestBuildLineupJSON(t *testing.T) {
	args := setLineupArgs{
		Goalkeeper:      1,
		CentralDefender: 2,
		CentralInnerMid: 3,
		Captain:         3,
		SetPiecesTaker:  2,
		PenaltyTakers:   []id.Player{3, 2},
		Tactic:          chpp.MatchTacticTypePressing,     // 1
		Attitude:        chpp.MatchTeamAttitudePlayItCool, // -1
		CoachModifier:   chpp.CoachModifier(5),            // 5
		ManMarker:       1,
		ManMarking:      99,
		Substitutions: []subArg{{
			PlayerIn: 4, PlayerOut: 2, OrderType: 0, Minute: 46, Position: 3, Behaviour: 2, Card: -1, Standing: -1,
		}},
	}
	raw, err := buildLineup(args)
	if err != nil {
		t.Fatal(err)
	}
	// Regression: the numeric fields must not be quoted strings (the bug CHPP rejected).
	if strings.Contains(raw, `"captain":"`) || strings.Contains(raw, `"tactic":"`) {
		t.Fatalf("numeric field serialized as a string: %s", raw)
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	for key, n := range map[string]int{"positions": 14, "bench": 14, "kickers": 11} {
		if got := len(m[key].([]any)); got != n {
			t.Errorf("%s: got %d, want %d", key, got, n)
		}
	}
	if m["captain"].(float64) != 3 {
		t.Errorf("captain = %v, want number 3", m["captain"])
	}
	settings := m["settings"].(map[string]any)
	if settings["newLineup"] != "" {
		t.Errorf("newLineup = %v, want empty string", settings["newLineup"])
	}
	if settings["tactic"].(float64) != 1 || settings["speechLevel"].(float64) != -1 || settings["coachModifier"].(float64) != 5 {
		t.Errorf("settings numbers wrong: %v", settings)
	}
	if settings["manMarkerPlayerId"].(float64) != 1 || settings["manMarkingPlayerId"].(float64) != 99 {
		t.Errorf("man-marking wrong: %v", settings)
	}
	pos := m["positions"].([]any)
	if pos[0].(map[string]any)["id"].(float64) != 1 || pos[3].(map[string]any)["id"].(float64) != 2 || pos[8].(map[string]any)["id"].(float64) != 3 {
		t.Errorf("position ids wrong")
	}
	k := m["kickers"].([]any)
	if k[0].(map[string]any)["id"].(float64) != 3 || k[1].(map[string]any)["id"].(float64) != 2 {
		t.Errorf("kicker order wrong")
	}
	subs := m["substitutions"].([]any)
	if len(subs) != 1 {
		t.Fatalf("subs len = %d", len(subs))
	}
	sub := subs[0].(map[string]any)
	if sub["orderType"] != "1" { // string, defaulted from 0
		t.Errorf("orderType = %v, want string \"1\"", sub["orderType"])
	}
	if sub["playerin"].(float64) != 4 || sub["playerout"].(float64) != 2 || sub["min"].(float64) != 46 || sub["pos"].(float64) != 3 {
		t.Errorf("sub fields wrong: %v", sub)
	}
}

func TestPlainConversion(t *testing.T) {
	// Skill enum → "name (n)".
	skill := plain(chpp.SkillLevel(8))
	if s, ok := skill.(string); !ok || !strings.Contains(s, "(8)") {
		t.Errorf("skill = %v, want string containing (8)", skill)
	}
	// Zero value dropped.
	if plain(chpp.SkillLevel(0)) != nil {
		t.Errorf("zero skill should be dropped")
	}
	// Booleans are kept, including false (a false flag is real information).
	if plain(false) != false {
		t.Errorf("false bool should be kept, got %v", plain(false))
	}
	if plain(true) != true {
		t.Errorf("true bool should be kept")
	}
	// Money passes through as a number.
	if m, ok := plain(chpp.Money(2500)).(chpp.Money); !ok || m != 2500 {
		t.Errorf("money = %v, want chpp.Money(2500)", plain(chpp.Money(2500)))
	}
	// Struct: exported non-zero fields kept, zero fields dropped.
	type sample struct {
		Skill chpp.SkillLevel
		Empty int
	}
	out, ok := plain(sample{Skill: 8}).(map[string]any)
	if !ok {
		t.Fatalf("struct plain = %T", plain(sample{Skill: 8}))
	}
	if _, has := out["Skill"]; !has {
		t.Errorf("Skill field dropped")
	}
	if _, has := out["Empty"]; has {
		t.Errorf("zero Empty field should be dropped")
	}
}

func TestParseDate(t *testing.T) {
	if tm, err := parseDate(""); err != nil || !tm.IsZero() {
		t.Errorf("empty: %v %v", tm, err)
	}
	if tm, err := parseDate("2026-01-02"); err != nil || tm.Year() != 2026 {
		t.Errorf("date: %v %v", tm, err)
	}
	if _, err := parseDate("not-a-date"); err == nil {
		t.Errorf("expected error for bad date")
	}
}

func TestSrcOr(t *testing.T) {
	if srcOr("") != chpp.SourceSystemHattrickMainSystem {
		t.Errorf("empty should default to main system")
	}
	if srcOr(chpp.SourceSystemYouthSystem) != chpp.SourceSystemYouthSystem {
		t.Errorf("explicit value should pass through")
	}
}

func TestConfirmed(t *testing.T) {
	if !confirmed(confirmOnly{confirmField{Confirm: true}}) {
		t.Errorf("confirm true not detected")
	}
	if confirmed(confirmOnly{confirmField{Confirm: false}}) {
		t.Errorf("confirm false reported as true")
	}
	if confirmed(struct{ X int }{X: 1}) {
		t.Errorf("struct without Confirm should be false")
	}
}

func TestCache(t *testing.T) {
	c := &cache{m: map[string]cached{}}
	calls := 0
	load := func() (string, error) { calls++; return "v", nil }

	for i := 0; i < 3; i++ {
		if v, err := c.get("tool", map[string]int{"a": 1}, load); err != nil || v != "v" {
			t.Fatalf("get: %v %v", v, err)
		}
	}
	if calls != 1 {
		t.Errorf("same args loaded %d times, want 1 (cached)", calls)
	}
	if _, err := c.get("tool", map[string]int{"a": 2}, load); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("different args: calls = %d, want 2", calls)
	}
	// Expiry: an entry older than the TTL reloads.
	c.mu.Lock()
	for k, e := range c.m {
		e.at = time.Now().Add(-2 * cacheTTL)
		c.m[k] = e
	}
	c.mu.Unlock()
	if _, err := c.get("tool", map[string]int{"a": 1}, load); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Errorf("expired entry not reloaded: calls = %d, want 3", calls)
	}
}
