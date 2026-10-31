package main

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/lucianoq/hattrick"
	"github.com/lucianoq/hattrick/api"
	"github.com/lucianoq/hattrick/chpp"
	"github.com/lucianoq/hattrick/chpp/id"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// cacheTTL keeps repeated questions from hitting the CHPP servers.
const cacheTTL = 5 * time.Minute

const instructions = `Read-only access to the authorised manager's Hattrick.org data via CHPP.
Omit an ID argument to get the manager's own team, series or match data.
Fields that are zero, false or empty are left out of results.
Skill values read "name (n)" on the 0 (non-existent) to 20 (divine) scale. Money is in SEK; divide by 10 for EUR display.`

type teamArgs struct {
	TeamID id.Team `json:"team_id,omitempty" jsonschema:"team ID; omit for your own primary team"`
}

type playerArgs struct {
	PlayerID id.Player `json:"player_id" jsonschema:"player ID"`
}

type matchArgs struct {
	MatchID id.Match `json:"match_id" jsonschema:"match ID"`
	Events  bool     `json:"events,omitempty" jsonschema:"include the match event list"`
}

type lineupArgs struct {
	MatchID id.Match `json:"match_id" jsonschema:"match ID"`
	TeamID  id.Team  `json:"team_id,omitempty" jsonschema:"team ID; omit for your own team"`
}

type seriesArgs struct {
	SeriesID id.Series `json:"series_id,omitempty" jsonschema:"series (LeagueLevelUnit) ID; omit for your own series"`
	Season   uint      `json:"season,omitempty" jsonschema:"season number; omit for the current season"`
}

type transferArgs struct {
	AgeMin    uint            `json:"age_min" jsonschema:"minimum age in years (17 or more)"`
	AgeMax    uint            `json:"age_max" jsonschema:"maximum age in years"`
	SkillID   chpp.SkillID    `json:"skill_id" jsonschema:"1 keeper, 2 stamina, 3 set pieces, 4 defending, 5 scoring, 6 winger, 7 passing, 8 playmaking, 9 trainer, 10 leadership, 11 experience"`
	SkillMin  chpp.SkillLevel `json:"skill_min" jsonschema:"minimum level of skill_id (0-20)"`
	SkillMax  chpp.SkillLevel `json:"skill_max" jsonschema:"maximum level of skill_id (0-20)"`
	PriceMax  chpp.Money      `json:"price_max,omitempty" jsonschema:"maximum current price in SEK"`
	PageIndex uint            `json:"page_index,omitempty" jsonschema:"result page, starting at 0"`
}

type noArgs struct{}

func newServer(ht hattrick.Client) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "hattrick-mcp", Version: version}, &mcp.ServerOptions{Instructions: instructions})
	c := &cache{m: map[string]cached{}}

	add(s, c, "get_team", "Team details: name, league, series, arena, trainer, fan club, form, team spirit and confidence.",
		func(a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyPrimaryTeam()
			}
			return ht.GetTeam(a.TeamID)
		})
	add(s, c, "get_players", "Senior squad list with age, form, stamina, skills (own team), TSI, salary, specialty, injuries and cards.",
		func(a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyPlayers()
			}
			return ht.GetPlayers(a.TeamID)
		})
	add(s, c, "get_player", "Full details of one player, including skills, transfer status and last match.",
		func(a playerArgs) (any, error) { return ht.GetPlayerDetails(a.PlayerID, true) })
	add(s, c, "get_matches", "Recent and upcoming matches of a team.",
		func(a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyMatches()
			}
			return ht.GetMatches(a.TeamID)
		})
	add(s, c, "get_match", "Match details: score, ratings per sector, possession, chances, tactics and optionally events.",
		func(a matchArgs) (any, error) {
			return ht.GetMatchDetails(a.MatchID, chpp.SourceSystemHattrickMainSystem, a.Events)
		})
	add(s, c, "get_match_lineup", "Lineup a team used in a played match, with positions, behaviours and ratings.",
		func(a lineupArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyMatchLineup(a.MatchID)
			}
			return ht.GetMatchLineup(a.MatchID, a.TeamID)
		})
	add(s, c, "get_training", "Training type, intensity, stamina share, trainer and recent training history.",
		func(a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyTraining()
			}
			return ht.GetTraining(a.TeamID)
		})
	add(s, c, "get_economy", "Your club finances: cash, weekly income and expenses, last week, fan mood and sponsors.",
		func(noArgs) (any, error) { return ht.GetFinancesForPrimaryTeam() })
	add(s, c, "get_series", "League table of a series.",
		func(a seriesArgs) (any, error) {
			if a.SeriesID == 0 {
				return ht.GetMySeries()
			}
			return ht.GetSeries(a.SeriesID)
		})
	add(s, c, "get_series_fixtures", "All fixtures and results of a series for a season.",
		func(a seriesArgs) (any, error) {
			switch {
			case a.SeriesID == 0 && a.Season == 0:
				return ht.GetMySeriesFixtures()
			case a.SeriesID == 0:
				return ht.GetMySeriesFixturesBySeason(a.Season)
			case a.Season == 0:
				return ht.GetSeriesFixtures(a.SeriesID)
			}
			return ht.GetSeriesFixturesBySeason(a.SeriesID, a.Season)
		})
	add(s, c, "get_staff", "Your trainer and staff members with levels and salaries.",
		func(noArgs) (any, error) { return ht.GetMyStaff() })
	add(s, c, "get_arena", "Arena capacity per seat type, expansions and income.",
		func(a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyArena()
			}
			return ht.GetArenaByTeamID(a.TeamID)
		})
	add(s, c, "get_youth_players", "Your youth academy players with known skills, potentials and last match.",
		func(noArgs) (any, error) { return ht.GetMyYouthPlayersDetails(false, true) })
	add(s, c, "search_transfers", "Search the transfer market by age and one skill range.",
		func(a transferArgs) (any, error) {
			return ht.GetTransferSearch(a.AgeMin, a.AgeMax, a.SkillID, a.SkillMin, a.SkillMax,
				api.TransferSearchFilters{PriceMax: a.PriceMax, PageIndex: a.PageIndex})
		})
	return s
}

// add registers a read-only tool whose result is returned as compact JSON text.
func add[In any](s *mcp.Server, c *cache, name, desc string, fn func(In) (any, error)) {
	t := &mcp.Tool{Name: name, Description: desc, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true}}
	mcp.AddTool(s, t, func(_ context.Context, _ *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
		text, err := c.get(name, in, func() (string, error) {
			v, err := fn(in)
			if err != nil {
				return "", err
			}
			b, err := json.Marshal(plain(v))
			return string(b), err
		})
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}, nil, nil
	})
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

	text, err := load()
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.m[key] = cached{text, time.Now()}
	c.mu.Unlock()
	return text, nil
}
