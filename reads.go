package main

import (
	"fmt"
	"time"

	"github.com/lucianoq/hattrick"
	"github.com/lucianoq/hattrick/api"
	"github.com/lucianoq/hattrick/chpp"
	"github.com/lucianoq/hattrick/chpp/id"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type teamArgs struct {
	TeamID id.Team `json:"team_id,omitempty" jsonschema:"team ID; omit for your own primary team"`
}

type userArgs struct {
	UserID id.User `json:"user_id,omitempty" jsonschema:"manager (user) ID; omit for yourself"`
}

type playerArgs struct {
	PlayerID id.Player `json:"player_id" jsonschema:"player ID"`
}

type playersArgs struct {
	TeamID  id.Team `json:"team_id,omitempty" jsonschema:"team ID; omit for your own team"`
	OrderBy string  `json:"order_by,omitempty" jsonschema:"sort order, e.g. age, salary, skill-specific; omit for default"`
}

type matchArgs struct {
	MatchID      id.Match          `json:"match_id" jsonschema:"match ID"`
	SourceSystem chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
	Events       bool              `json:"events,omitempty" jsonschema:"include the match event list"`
}

type lineupArgs struct {
	MatchID      id.Match          `json:"match_id" jsonschema:"match ID"`
	TeamID       id.Team           `json:"team_id,omitempty" jsonschema:"team ID; omit for your own team"`
	SourceSystem chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
}

type matchOrdersArgs struct {
	MatchID      id.Match          `json:"match_id" jsonschema:"match ID"`
	TeamID       id.Team           `json:"team_id" jsonschema:"your team ID for this match"`
	SourceSystem chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
}

type predictArgs struct {
	MatchID      id.Match          `json:"match_id" jsonschema:"match ID"`
	TeamID       id.Team           `json:"team_id" jsonschema:"your team ID for this match"`
	SourceSystem chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
	Lineup       string            `json:"lineup" jsonschema:"lineup JSON as accepted by the CHPP matchorders format"`
}

type matchesArgs struct {
	TeamID id.Team `json:"team_id,omitempty" jsonschema:"team ID; omit for your own team"`
	Until  string  `json:"until,omitempty" jsonschema:"only matches up to this date (YYYY-MM-DD); omit for the default window"`
}

type archiveArgs struct {
	TeamID     id.Team `json:"team_id,omitempty" jsonschema:"team ID; omit for your own team"`
	Start      string  `json:"start,omitempty" jsonschema:"start date YYYY-MM-DD (use with end)"`
	End        string  `json:"end,omitempty" jsonschema:"end date YYYY-MM-DD (use with start)"`
	Season     uint    `json:"season,omitempty" jsonschema:"season number instead of start/end"`
	IncludeHTO bool    `json:"include_hto,omitempty" jsonschema:"include Hattrick International tournament matches (date range only)"`
}

type seriesArgs struct {
	SeriesID id.Series `json:"series_id,omitempty" jsonschema:"series (LeagueLevelUnit) ID; omit for your own series"`
	Season   uint      `json:"season,omitempty" jsonschema:"season number; omit for the current season"`
}

type leagueLevelsArgs struct {
	LeagueID id.League `json:"league_id,omitempty" jsonschema:"league (country) ID; omit for your own country"`
}

type leagueArgs struct {
	LeagueID id.League `json:"league_id" jsonschema:"league (country) ID"`
}

type countryArgs struct {
	CountryID id.Country `json:"country_id" jsonschema:"country ID"`
}

type regionArgs struct {
	RegionID id.Region `json:"region_id,omitempty" jsonschema:"region ID; omit for your own region"`
}

type leagueStatArgs struct {
	LeagueID id.League `json:"league_id" jsonschema:"league (country) ID"`
}

type pagedTeamArgs struct {
	TeamID    id.Team `json:"team_id,omitempty" jsonschema:"team ID; omit for your own team"`
	PageIndex uint    `json:"page_index,omitempty" jsonschema:"result page, starting at 0"`
	PageSize  uint    `json:"page_size,omitempty" jsonschema:"results per page"`
}

type pagedUserArgs struct {
	UserID    id.User `json:"user_id,omitempty" jsonschema:"manager (user) ID; omit for yourself"`
	PageIndex uint    `json:"page_index,omitempty" jsonschema:"result page, starting at 0"`
	PageSize  uint    `json:"page_size,omitempty" jsonschema:"results per page"`
}

type transfersTeamArgs struct {
	TeamID    id.Team `json:"team_id,omitempty" jsonschema:"team ID; omit for your own team"`
	PageIndex uint    `json:"page_index,omitempty" jsonschema:"result page, starting at 0"`
}

type skillFilterArg struct {
	SkillID chpp.SkillID    `json:"skill_id" jsonschema:"1 keeper, 2 stamina, 3 set pieces, 4 defending, 5 scoring, 6 winger, 7 passing, 8 playmaking, 9 trainer, 10 leadership, 11 experience"`
	Min     chpp.SkillLevel `json:"min" jsonschema:"minimum level (0-20)"`
	Max     chpp.SkillLevel `json:"max" jsonschema:"maximum level (0-20), within 4 of min"`
}

type transferArgs struct {
	AgeMin   uint            `json:"age_min" jsonschema:"minimum age in years (17 or more)"`
	AgeMax   uint            `json:"age_max" jsonschema:"maximum age in years"`
	SkillID  chpp.SkillID    `json:"skill_id" jsonschema:"primary skill: 1 keeper, 2 stamina, 3 set pieces, 4 defending, 5 scoring, 6 winger, 7 passing, 8 playmaking, 9 trainer, 10 leadership, 11 experience"`
	SkillMin chpp.SkillLevel `json:"skill_min" jsonschema:"minimum level of skill_id (0-20)"`
	SkillMax chpp.SkillLevel `json:"skill_max" jsonschema:"maximum level of skill_id (0-20), within 4 of skill_min"`

	AgeDaysMin       uint             `json:"age_days_min,omitempty" jsonschema:"minimum age days (0-111), combined with age_min"`
	AgeDaysMax       uint             `json:"age_days_max,omitempty" jsonschema:"maximum age days (0-111), combined with age_max"`
	AdditionalSkills []skillFilterArg `json:"additional_skills,omitempty" jsonschema:"up to 3 further skill filters beyond the primary one"`
	Specialty        chpp.SpecialtyID `json:"specialty,omitempty" jsonschema:"0 none, 1 technical, 2 quick, 3 powerful, 4 unpredictable, 5 head, 6 resilient, 8 support"`
	NativeCountryID  id.Country       `json:"native_country_id,omitempty" jsonschema:"restrict to players native to this country"`
	TSIMin           uint             `json:"tsi_min,omitempty" jsonschema:"minimum TSI"`
	TSIMax           uint             `json:"tsi_max,omitempty" jsonschema:"maximum TSI"`
	PriceMin         chpp.Money       `json:"price_min,omitempty" jsonschema:"minimum current price in SEK"`
	PriceMax         chpp.Money       `json:"price_max,omitempty" jsonschema:"maximum current price in SEK"`
	PageSize         uint             `json:"page_size,omitempty" jsonschema:"results per page (default 25)"`
	PageIndex        uint             `json:"page_index,omitempty" jsonschema:"result page, starting at 0"`
}

type youthPlayersArgs struct {
	YouthTeamID   id.YouthTeam `json:"youth_team_id,omitempty" jsonschema:"youth team ID; omit for your own youth team"`
	ShowScoutCall bool         `json:"show_scout_call,omitempty" jsonschema:"include the scout's comments"`
	ShowLastMatch bool         `json:"show_last_match,omitempty" jsonschema:"include each player's last match"`
}

type youthPlayerArgs struct {
	YouthPlayerID id.YouthPlayer `json:"youth_player_id" jsonschema:"youth player ID"`
	ShowScoutCall bool           `json:"show_scout_call,omitempty" jsonschema:"include the scout's comments"`
	ShowLastMatch bool           `json:"show_last_match,omitempty" jsonschema:"include the player's last match"`
}

type youthTeamArgs struct {
	YouthTeamID id.YouthTeam `json:"youth_team_id,omitempty" jsonschema:"youth team ID; omit for your own youth team"`
	ShowScouts  bool         `json:"show_scouts,omitempty" jsonschema:"include the scouts"`
}

type youthLeagueArgs struct {
	YouthLeagueID id.YouthLeague `json:"youth_league_id,omitempty" jsonschema:"youth league ID; omit for your own youth league"`
	Season        uint           `json:"season,omitempty" jsonschema:"season number; omit for the current season"`
}

type tournamentArgs struct {
	TournamentID id.Tournament `json:"tournament_id" jsonschema:"tournament ID"`
	Season       uint          `json:"season,omitempty" jsonschema:"season number; omit for the current season"`
}

type tournamentFixturesArgs struct {
	TournamentID id.Tournament `json:"tournament_id" jsonschema:"tournament ID"`
	Season       uint          `json:"season,omitempty" jsonschema:"season number"`
	MatchRound   uint          `json:"match_round,omitempty" jsonschema:"match round instead of season"`
}

type tournamentTablesArgs struct {
	TournamentID  id.Tournament `json:"tournament_id" jsonschema:"tournament ID"`
	Season        uint          `json:"season,omitempty" jsonschema:"season number"`
	WorldCupRound uint          `json:"world_cup_round,omitempty" jsonschema:"world cup round instead of season"`
}

type nationalTeamsArgs struct {
	LeagueOfficeTypeID chpp.LeagueOfficeTypeID `json:"league_office_type_id,omitempty" jsonschema:"2 national teams (default), 4 U20 teams"`
}

type nationalTeamArgs struct {
	NationalTeamID id.NationalTeam `json:"national_team_id" jsonschema:"national team ID"`
}

type nationalPlayerStatsArgs struct {
	NationalTeamID    id.NationalTeam `json:"national_team_id" jsonschema:"national team ID"`
	MatchTypeCategory string          `json:"match_type_category,omitempty" jsonschema:"match type category filter"`
	ShowAll           bool            `json:"show_all,omitempty" jsonschema:"include all players"`
}

type worldCupGroupsArgs struct {
	CupID      chpp.WorldCupID `json:"cup_id" jsonschema:"39 or 137 world cup, 45 or 149 U20"`
	Season     uint            `json:"season" jsonschema:"season number"`
	MatchRound uint            `json:"match_round" jsonschema:"match round"`
}

type worldCupMatchesArgs struct {
	CupID           chpp.WorldCupID `json:"cup_id" jsonschema:"39 or 137 world cup, 45 or 149 U20"`
	Season          uint            `json:"season" jsonschema:"season number"`
	CupSeriesUnitID uint            `json:"cup_series_unit_id" jsonschema:"cup series unit ID"`
}

type allianceArgs struct {
	AllianceID id.Alliance `json:"alliance_id" jsonschema:"alliance (federation) ID"`
}

type allianceMembersArgs struct {
	AllianceID id.Alliance `json:"alliance_id" jsonschema:"alliance (federation) ID"`
	Subset     uint        `json:"subset,omitempty" jsonschema:"members subset index; omit for all members"`
}

type searchFederationsArgs struct {
	By         string      `json:"by,omitempty" jsonschema:"name (default), abbreviation or description"`
	Query      string      `json:"query" jsonschema:"text to search for"`
	LanguageID id.Language `json:"language_id,omitempty" jsonschema:"language ID to search in"`
	PageIndex  uint        `json:"page_index,omitempty" jsonschema:"result page, starting at 0"`
}

type bookmarksArgs struct {
	BookmarkType chpp.BookmarkType `json:"bookmark_type,omitempty" jsonschema:"0 all, 1 teams/users, 2 players, 3 matches, 4 forum users, 5 league units, 6 youth teams, 7 youth players, 8 youth matches"`
}

type laddersArgs struct {
	TeamID id.Team `json:"team_id,omitempty" jsonschema:"team ID; omit for your own ladders"`
}

type ladderDetailsArgs struct {
	LadderID  id.Ladder `json:"ladder_id" jsonschema:"ladder ID"`
	TeamID    id.Team   `json:"team_id,omitempty" jsonschema:"team ID to centre the page on"`
	PageIndex uint      `json:"page_index,omitempty" jsonschema:"result page, starting at 0"`
	PageSize  uint      `json:"page_size,omitempty" jsonschema:"results per page"`
}

type liveMatchesArgs struct {
	SourceSystem          chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
	IncludeStartingLineup bool              `json:"include_starting_lineup,omitempty" jsonschema:"include the starting lineup"`
	UseLiveEventsAndTexts bool              `json:"use_live_events_and_texts,omitempty" jsonschema:"include live events and texts"`
}

type liveMatchArgs struct {
	MatchID               id.Match          `json:"match_id" jsonschema:"match ID"`
	SourceSystem          chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
	IncludeStartingLineup bool              `json:"include_starting_lineup,omitempty" jsonschema:"include the starting lineup"`
	UseLiveEventsAndTexts bool              `json:"use_live_events_and_texts,omitempty" jsonschema:"include live events and texts"`
}

type searchArgs struct {
	SearchType     chpp.SearchType `json:"search_type" jsonschema:"0 players, 1 arenas, 2 managers, 3 series, 4 teams, 5 regions, 6 match"`
	Query          string          `json:"query" jsonschema:"text to search for"`
	SearchString2  string          `json:"search_string2,omitempty" jsonschema:"second search string (players only)"`
	SearchID       uint            `json:"search_id,omitempty" jsonschema:"ID to search for instead of text, for some types"`
	SearchLeagueID int             `json:"search_league_id,omitempty" jsonschema:"-1 all leagues, 0 your league, or a league ID"`
	PageIndex      uint            `json:"page_index,omitempty" jsonschema:"result page, starting at 0"`
}

type cupMatchesArgs struct {
	Cup    id.Cup `json:"cup_id" jsonschema:"cup ID"`
	Season uint   `json:"season,omitempty" jsonschema:"season number; omit for the latest matches"`
	Round  uint   `json:"round,omitempty" jsonschema:"cup round; omit for the latest matches"`
}

type translationsArgs struct {
	LanguageID id.Language `json:"language_id" jsonschema:"language ID"`
}

type avatarsArgs struct {
	Kind string `json:"kind,omitempty" jsonschema:"players (default), hof, staff or youth"`
	ID   uint   `json:"id,omitempty" jsonschema:"team ID (players/hof/staff) or youth team ID (youth); omit for your own"`
}

func parseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date %q, use YYYY-MM-DD", s)
}

func registerReads(s *mcp.Server, acct *account, c *cache) {
	// Manager, team, club, region, achievements.
	add(s, acct, c, "get_manager", "Manager profile: username, country, supporter tier, age in Hattrick and more.",
		func(ht hattrick.Client, a userArgs) (any, error) {
			if a.UserID == 0 {
				return ht.GetMe()
			}
			return ht.GetManager(a.UserID)
		})
	add(s, acct, c, "get_team", "Team details: name, league, series, arena, trainer, fan club, form, team spirit and confidence.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyPrimaryTeam()
			}
			return ht.GetTeam(a.TeamID)
		})
	add(s, acct, c, "get_teams", "All teams of a manager (a manager may have several).",
		func(ht hattrick.Client, a userArgs) (any, error) {
			if a.UserID == 0 {
				return ht.GetMyTeams()
			}
			return ht.GetTeamsByUser(a.UserID)
		})
	add(s, acct, c, "get_club", "Club-wide items: staff count, youth squad, trainers, specialists and facilities.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyClub()
			}
			return ht.GetClubByID(a.TeamID)
		})
	add(s, acct, c, "get_region", "Region details: name, number of users and teams, weather.",
		func(ht hattrick.Client, a regionArgs) (any, error) {
			if a.RegionID == 0 {
				return ht.GetMyRegion()
			}
			return ht.GetRegion(a.RegionID)
		})
	add(s, acct, c, "get_achievements", "A manager's achievements and the points earned.",
		func(ht hattrick.Client, a userArgs) (any, error) {
			if a.UserID == 0 {
				return ht.GetMyAchievements()
			}
			return ht.GetAchievements(a.UserID)
		})

	// Players.
	add(s, acct, c, "get_players", "Senior squad list with age, form, stamina, skills (own team), TSI, salary, specialty, injuries and cards.",
		func(ht hattrick.Client, a playersArgs) (any, error) {
			switch {
			case a.TeamID == 0 && a.OrderBy != "":
				return ht.GetMyPlayersOrderedBy(a.OrderBy)
			case a.TeamID == 0:
				return ht.GetMyPlayers()
			case a.OrderBy != "":
				return ht.GetPlayersOrderedBy(a.TeamID, a.OrderBy)
			}
			return ht.GetPlayers(a.TeamID)
		})
	add(s, acct, c, "get_player", "Full details of one player, including skills, transfer status and last match.",
		func(ht hattrick.Client, a playerArgs) (any, error) { return ht.GetPlayerDetails(a.PlayerID, true) })
	add(s, acct, c, "get_player_events", "Career events of a player (milestones, transfers, bookings).",
		func(ht hattrick.Client, a playerArgs) (any, error) { return ht.GetPlayerEvents(a.PlayerID) })
	add(s, acct, c, "get_player_transfer_history", "All transfers of a player with dates and prices.",
		func(ht hattrick.Client, a playerArgs) (any, error) { return ht.GetPlayerTransferHistory(a.PlayerID) })
	add(s, acct, c, "get_old_players", "Former senior players of a team.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyOldPlayers()
			}
			return ht.GetOldPlayers(a.TeamID)
		})
	add(s, acct, c, "get_old_coaches", "Former trainers of a team.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyOldCoaches()
			}
			return ht.GetOldCoaches(a.TeamID)
		})
	add(s, acct, c, "get_hof_players", "Hall of fame players of a team.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetHOFPlayersForPrimaryTeam()
			}
			return ht.GetHOFPlayers(a.TeamID)
		})
	add(s, acct, c, "get_avatars", "Layered avatar image URLs for players, hall of fame, staff or youth.",
		func(ht hattrick.Client, a avatarsArgs) (any, error) {
			switch a.Kind {
			case "hof":
				if a.ID == 0 {
					return ht.GetAvatarsMyHallOfFame()
				}
				return ht.GetAvatarsHallOfFame(id.Team(a.ID))
			case "staff":
				if a.ID == 0 {
					return ht.GetMyStaffAvatars()
				}
				return ht.GetStaffAvatars(id.Team(a.ID))
			case "youth":
				if a.ID == 0 {
					return ht.GetMyYouthAvatars()
				}
				return ht.GetYouthAvatars(id.YouthTeam(a.ID))
			default:
				if a.ID == 0 {
					return ht.GetAvatarsMyPlayers()
				}
				return ht.GetAvatarsPlayers(id.Team(a.ID))
			}
		})

	// Matches.
	add(s, acct, c, "get_matches", "Recent and upcoming senior matches of a team; pass until to extend the window.",
		func(ht hattrick.Client, a matchesArgs) (any, error) {
			until, err := parseDate(a.Until)
			if err != nil {
				return nil, err
			}
			switch {
			case a.TeamID == 0 && until.IsZero():
				return ht.GetMyMatches()
			case a.TeamID == 0:
				return ht.GetMyMatchesUntil(until)
			case until.IsZero():
				return ht.GetMatches(a.TeamID)
			}
			return ht.GetMatchesUntil(a.TeamID, until)
		})
	add(s, acct, c, "get_youth_matches", "Recent and upcoming youth matches of a team; pass until to extend the window.",
		func(ht hattrick.Client, a matchesArgs) (any, error) {
			until, err := parseDate(a.Until)
			if err != nil {
				return nil, err
			}
			switch {
			case a.TeamID == 0 && until.IsZero():
				return ht.GetMyYouthMatches()
			case a.TeamID == 0:
				return ht.GetMyYouthMatchesUntil(until)
			case until.IsZero():
				return ht.GetYouthMatches(a.TeamID)
			}
			return ht.GetYouthMatchesUntil(a.TeamID, until)
		})
	add(s, acct, c, "get_matches_archive", "Played senior matches in a date range or a season.",
		func(ht hattrick.Client, a archiveArgs) (any, error) {
			if a.Season != 0 {
				if a.TeamID == 0 {
					return ht.GetMyMatchesArchiveBySeason(a.Season)
				}
				return ht.GetMatchesArchiveBySeason(a.TeamID, a.Season)
			}
			start, err := parseDate(a.Start)
			if err != nil {
				return nil, err
			}
			end, err := parseDate(a.End)
			if err != nil {
				return nil, err
			}
			switch {
			case a.IncludeHTO && a.TeamID == 0:
				return ht.GetMyMatchesArchiveIncludeHTO(start, end)
			case a.IncludeHTO:
				return ht.GetMatchesArchiveIncludeHTO(a.TeamID, start, end)
			case a.TeamID == 0:
				return ht.GetMyMatchesArchive(start, end)
			}
			return ht.GetMatchesArchive(a.TeamID, start, end)
		})
	add(s, acct, c, "get_youth_matches_archive", "Played youth matches in a date range.",
		func(ht hattrick.Client, a archiveArgs) (any, error) {
			start, err := parseDate(a.Start)
			if err != nil {
				return nil, err
			}
			end, err := parseDate(a.End)
			if err != nil {
				return nil, err
			}
			if a.TeamID == 0 {
				return ht.GetMyYouthMatchesArchive(start, end)
			}
			return ht.GetYouthMatchesArchive(a.TeamID, start, end)
		})
	add(s, acct, c, "get_match", "Match details: score, ratings per sector, possession, chances, tactics and optionally events.",
		func(ht hattrick.Client, a matchArgs) (any, error) {
			return ht.GetMatchDetails(a.MatchID, srcOr(a.SourceSystem), a.Events)
		})
	add(s, acct, c, "get_match_lineup", "Lineup a team used in a played match, with positions, behaviours and ratings.",
		func(ht hattrick.Client, a lineupArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyMatchLineup(a.MatchID)
			}
			return ht.GetMatchLineupBySourceSystem(a.MatchID, a.TeamID, srcOr(a.SourceSystem))
		})
	add(s, acct, c, "get_match_orders", "The orders currently set for an upcoming match: lineup, positions, orders and substitutions.",
		func(ht hattrick.Client, a matchOrdersArgs) (any, error) {
			return ht.GetMatchOrders(a.MatchID, a.TeamID, srcOr(a.SourceSystem))
		})
	add(s, acct, c, "predict_ratings", "Predict the ratings a lineup would give for a match, without saving it.",
		func(ht hattrick.Client, a predictArgs) (any, error) {
			return ht.PredictRatings(a.MatchID, a.TeamID, srcOr(a.SourceSystem), a.Lineup)
		})
	add(s, acct, c, "get_cup_matches", "Matches of a cup, for a season and round or the latest.",
		func(ht hattrick.Client, a cupMatchesArgs) (any, error) {
			if a.Season == 0 && a.Round == 0 {
				return ht.GetCupMatchesLast(a.Cup)
			}
			return ht.GetCupMatches(a.Cup, a.Season, a.Round)
		})

	// League, series, world.
	add(s, acct, c, "get_series", "League table of a series.",
		func(ht hattrick.Client, a seriesArgs) (any, error) {
			if a.SeriesID == 0 {
				return ht.GetMySeries()
			}
			return ht.GetSeries(a.SeriesID)
		})
	add(s, acct, c, "get_series_fixtures", "All fixtures and results of a series for a season.",
		func(ht hattrick.Client, a seriesArgs) (any, error) {
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
	add(s, acct, c, "get_league_levels", "The division levels of a country and their series counts.",
		func(ht hattrick.Client, a leagueLevelsArgs) (any, error) {
			if a.LeagueID == 0 {
				return ht.GetMyLeagueLevels()
			}
			return ht.GetLeagueLevels(a.LeagueID)
		})
	add(s, acct, c, "get_world", "All Hattrick leagues (countries) with season and economy summary.",
		func(ht hattrick.Client, _ noArgs) (any, error) { return ht.GetWorld() })
	add(s, acct, c, "get_league", "One league (country): season, round, dates, economy and more.",
		func(ht hattrick.Client, a leagueArgs) (any, error) { return ht.GetLeague(a.LeagueID) })
	add(s, acct, c, "get_country", "A country's league data by country ID.",
		func(ht hattrick.Client, a countryArgs) (any, error) { return ht.GetCountry(a.CountryID) })
	add(s, acct, c, "get_world_languages", "The languages available in Hattrick.",
		func(ht hattrick.Client, _ noArgs) (any, error) { return ht.GetWorldLanguages() })
	add(s, acct, c, "get_translations", "CHPP term translations for a language.",
		func(ht hattrick.Client, a translationsArgs) (any, error) { return ht.GetTranslations(a.LanguageID) })

	// Economy, fans, supporters.
	add(s, acct, c, "get_economy", "Club finances: cash, weekly income and expenses, last week, fan mood and sponsors.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetFinancesForPrimaryTeam()
			}
			return ht.GetFinances(a.TeamID)
		})
	add(s, acct, c, "get_fans", "Fan club size and mood for a team.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetFansForPrimaryTeam()
			}
			return ht.GetFans(a.TeamID)
		})
	add(s, acct, c, "get_supported_teams", "Teams a manager supports.",
		func(ht hattrick.Client, a pagedUserArgs) (any, error) {
			if a.UserID == 0 {
				return ht.GetMySupportedTeams(a.PageIndex, a.PageSize)
			}
			return ht.GetSupportedTeams(a.UserID, a.PageIndex, a.PageSize)
		})
	add(s, acct, c, "get_supporters", "Managers who support a team.",
		func(ht hattrick.Client, a pagedTeamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMySupporters(a.PageIndex, a.PageSize)
			}
			return ht.GetSupportersForTeam(a.TeamID, a.PageIndex, a.PageSize)
		})

	// Training, staff, arena.
	add(s, acct, c, "get_training", "Training type, intensity, stamina share, trainer and recent training history.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyTraining()
			}
			return ht.GetTraining(a.TeamID)
		})
	add(s, acct, c, "get_training_stats", "Training choices across a country's top division.",
		func(ht hattrick.Client, a leagueStatArgs) (any, error) { return ht.GetTrainingStats(a.LeagueID) })
	add(s, acct, c, "get_training_events", "Skill-up history of one player.",
		func(ht hattrick.Client, a playerArgs) (any, error) { return ht.GetTrainingEvents(a.PlayerID) })
	add(s, acct, c, "get_staff", "Trainer and staff members of a team with levels and salaries.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyStaff()
			}
			return ht.GetStaff(a.TeamID)
		})
	add(s, acct, c, "get_arena", "Arena capacity per seat type, expansions and income.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyArena()
			}
			return ht.GetArenaByTeamID(a.TeamID)
		})
	add(s, acct, c, "get_arena_stats", "Arena attendance stats across a country's top division.",
		func(ht hattrick.Client, a leagueStatArgs) (any, error) { return ht.GetArenaStats(a.LeagueID) })
	add(s, acct, c, "get_arena_supporter_stats", "Your arena's supporter attendance statistics.",
		func(ht hattrick.Client, _ noArgs) (any, error) { return ht.GetMyArenaSupporterStats() })

	// Transfers.
	add(s, acct, c, "search_transfers", "Search the transfer market by age, a primary skill range and optional filters.",
		func(ht hattrick.Client, a transferArgs) (any, error) {
			extra := make([]api.TransferSearchSkillFilter, 0, len(a.AdditionalSkills))
			for _, f := range a.AdditionalSkills {
				extra = append(extra, api.TransferSearchSkillFilter{SkillType: f.SkillID, Min: f.Min, Max: f.Max})
			}
			return ht.GetTransferSearch(a.AgeMin, a.AgeMax, a.SkillID, a.SkillMin, a.SkillMax,
				api.TransferSearchFilters{
					AgeDaysMin:       a.AgeDaysMin,
					AgeDaysMax:       a.AgeDaysMax,
					AdditionalSkills: extra,
					Specialty:        a.Specialty,
					NativeCountryID:  a.NativeCountryID,
					TSIMin:           a.TSIMin,
					TSIMax:           a.TSIMax,
					PriceMin:         a.PriceMin,
					PriceMax:         a.PriceMax,
					PageSize:         a.PageSize,
					PageIndex:        a.PageIndex,
				})
		})
	add(s, acct, c, "get_current_bids", "Players you are currently bidding on.",
		func(ht hattrick.Client, _ noArgs) (any, error) { return ht.GetCurrentBids() })
	add(s, acct, c, "get_transfers_team", "A team's transfer history (players bought and sold) by page.",
		func(ht hattrick.Client, a transfersTeamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyTransfersTeam(a.PageIndex)
			}
			return ht.GetTransfersTeam(a.TeamID, a.PageIndex)
		})

	// Youth.
	add(s, acct, c, "get_youth_players", "Youth academy players with known skills, potentials and last match.",
		func(ht hattrick.Client, a youthPlayersArgs) (any, error) {
			if a.YouthTeamID == 0 {
				return ht.GetMyYouthPlayersDetails(a.ShowScoutCall, a.ShowLastMatch)
			}
			return ht.GetYouthPlayersDetails(a.YouthTeamID, a.ShowScoutCall, a.ShowLastMatch)
		})
	add(s, acct, c, "get_youth_player", "Full details of one youth player.",
		func(ht hattrick.Client, a youthPlayerArgs) (any, error) {
			return ht.GetYouthPlayer(a.YouthPlayerID, a.ShowScoutCall, a.ShowLastMatch)
		})
	add(s, acct, c, "get_youth_team", "Youth team details and scouts.",
		func(ht hattrick.Client, a youthTeamArgs) (any, error) {
			if a.YouthTeamID == 0 {
				d, scouts, err := ht.GetMyYouthTeam(a.ShowScouts)
				return map[string]any{"team": d, "scouts": scouts}, err
			}
			d, scouts, err := ht.GetYouthTeam(a.YouthTeamID, a.ShowScouts)
			return map[string]any{"team": d, "scouts": scouts}, err
		})
	add(s, acct, c, "get_youth_league", "Youth league table.",
		func(ht hattrick.Client, a youthLeagueArgs) (any, error) {
			if a.YouthLeagueID == 0 {
				return ht.GetMyYouthLeague()
			}
			return ht.GetYouthLeague(a.YouthLeagueID)
		})
	add(s, acct, c, "get_youth_league_fixtures", "Youth league fixtures for a season.",
		func(ht hattrick.Client, a youthLeagueArgs) (any, error) {
			if a.YouthLeagueID == 0 {
				return ht.GetMyYouthLeagueFixtures()
			}
			return ht.GetYouthLeagueFixtures(a.YouthLeagueID, a.Season)
		})

	// Tournaments.
	add(s, acct, c, "get_tournaments", "Tournaments a team takes part in.",
		func(ht hattrick.Client, a teamArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyTournaments()
			}
			return ht.GetTournaments(a.TeamID)
		})
	add(s, acct, c, "get_tournament", "Details of one tournament.",
		func(ht hattrick.Client, a tournamentArgs) (any, error) {
			if a.Season == 0 {
				return ht.GetTournament(a.TournamentID)
			}
			return ht.GetTournamentBySeason(a.TournamentID, a.Season)
		})
	add(s, acct, c, "get_tournament_fixtures", "Fixtures of a tournament, by season or match round.",
		func(ht hattrick.Client, a tournamentFixturesArgs) (any, error) {
			switch {
			case a.MatchRound != 0:
				return ht.GetTournamentFixturesByMatchRound(a.TournamentID, a.MatchRound)
			case a.Season != 0:
				return ht.GetTournamentFixturesBySeason(a.TournamentID, a.Season)
			}
			return ht.GetTournamentFixtures(a.TournamentID)
		})
	add(s, acct, c, "get_tournament_tables", "League tables of a tournament, by season or world cup round.",
		func(ht hattrick.Client, a tournamentTablesArgs) (any, error) {
			switch {
			case a.WorldCupRound != 0:
				return ht.GetTournamentLeagueTablesByWorldCupRound(a.TournamentID, a.WorldCupRound)
			case a.Season != 0:
				return ht.GetTournamentLeagueTablesBySeason(a.TournamentID, a.Season)
			}
			return ht.GetTournamentLeagueTables(a.TournamentID)
		})

	// National teams.
	add(s, acct, c, "get_national_teams", "List of national or U20 teams.",
		func(ht hattrick.Client, a nationalTeamsArgs) (any, error) {
			if a.LeagueOfficeTypeID == 0 {
				a.LeagueOfficeTypeID = chpp.LeagueOfficeTypeNationalTeams
			}
			return ht.GetNationalTeams(a.LeagueOfficeTypeID)
		})
	add(s, acct, c, "get_national_team", "Details of one national team.",
		func(ht hattrick.Client, a nationalTeamArgs) (any, error) {
			return ht.GetNationalTeamDetails(a.NationalTeamID)
		})
	add(s, acct, c, "get_national_team_players", "Players of a national team.",
		func(ht hattrick.Client, a nationalTeamArgs) (any, error) {
			return ht.GetNationalTeamPlayers(a.NationalTeamID)
		})
	add(s, acct, c, "get_national_team_player_stats", "Player appearance stats for a national team.",
		func(ht hattrick.Client, a nationalPlayerStatsArgs) (any, error) {
			return ht.GetNationalTeamPlayerStats(a.NationalTeamID, a.MatchTypeCategory, a.ShowAll)
		})
	add(s, acct, c, "get_national_team_matches", "Matches of national or U20 teams.",
		func(ht hattrick.Client, a nationalTeamsArgs) (any, error) {
			if a.LeagueOfficeTypeID == 0 {
				a.LeagueOfficeTypeID = chpp.LeagueOfficeTypeNationalTeams
			}
			return ht.GetNationalTeamMatches(a.LeagueOfficeTypeID)
		})
	add(s, acct, c, "get_world_cup_groups", "World cup group standings for a season and round.",
		func(ht hattrick.Client, a worldCupGroupsArgs) (any, error) {
			return ht.GetWorldCupGroups(a.CupID, a.Season, a.MatchRound)
		})
	add(s, acct, c, "get_world_cup_matches", "World cup matches for a season and series unit.",
		func(ht hattrick.Client, a worldCupMatchesArgs) (any, error) {
			return ht.GetWorldCupMatches(a.CupID, a.Season, a.CupSeriesUnitID)
		})

	// Alliances (federations).
	add(s, acct, c, "get_my_alliances", "Alliances (federations) you belong to.",
		func(ht hattrick.Client, _ noArgs) (any, error) { return ht.GetMyAlliances() })
	add(s, acct, c, "get_alliance", "Details of one alliance (federation).",
		func(ht hattrick.Client, a allianceArgs) (any, error) { return ht.GetAllianceDetails(a.AllianceID) })
	add(s, acct, c, "get_alliance_members", "Members of an alliance.",
		func(ht hattrick.Client, a allianceMembersArgs) (any, error) {
			if a.Subset == 0 {
				return ht.GetAllianceDetailsMembers(a.AllianceID)
			}
			return ht.GetAllianceDetailsMembersSubset(a.AllianceID, a.Subset)
		})
	add(s, acct, c, "get_alliance_roles", "Roles within an alliance.",
		func(ht hattrick.Client, a allianceArgs) (any, error) { return ht.GetAllianceDetailsRoles(a.AllianceID) })
	add(s, acct, c, "search_alliances", "Search alliances by name, abbreviation or description.",
		func(ht hattrick.Client, a searchFederationsArgs) (any, error) {
			switch a.By {
			case "abbreviation":
				return ht.GetAlliancesAbbreviationIncludes(a.Query, a.LanguageID, a.PageIndex)
			case "description":
				return ht.GetAlliancesDescriptionIncludes(a.Query, a.LanguageID, a.PageIndex)
			default:
				return ht.GetAlliancesNameStartsWith(a.Query, a.LanguageID, a.PageIndex)
			}
		})

	// Bookmarks, ladders, live, search.
	add(s, acct, c, "get_bookmarks", "Your bookmarks of a given type.",
		func(ht hattrick.Client, a bookmarksArgs) (any, error) { return ht.GetBookmarks(a.BookmarkType) })
	add(s, acct, c, "get_ladders", "Ladder challenges a team is in.",
		func(ht hattrick.Client, a laddersArgs) (any, error) {
			if a.TeamID == 0 {
				return ht.GetMyLadders()
			}
			return ht.GetLadders(a.TeamID)
		})
	add(s, acct, c, "get_ladder_details", "Standings of a ladder, optionally centred on a team or page.",
		func(ht hattrick.Client, a ladderDetailsArgs) (any, error) {
			switch {
			case a.TeamID != 0:
				return ht.GetLadderDetailsForTeam(a.LadderID, a.TeamID)
			case a.PageSize != 0:
				return ht.GetLadderDetailsPage(a.LadderID, a.PageIndex, a.PageSize)
			}
			return ht.GetLadderDetails(a.LadderID)
		})
	add(s, acct, c, "get_live_matches", "Live match data for the matches on your live list.",
		func(ht hattrick.Client, a liveMatchesArgs) (any, error) {
			return ht.GetLiveMatches(srcOr(a.SourceSystem), a.IncludeStartingLineup, a.UseLiveEventsAndTexts)
		})
	add(s, acct, c, "get_live_match", "Live match data for one match.",
		func(ht hattrick.Client, a liveMatchArgs) (any, error) {
			return ht.GetLiveMatch(a.MatchID, srcOr(a.SourceSystem), a.IncludeStartingLineup, a.UseLiveEventsAndTexts)
		})
	add(s, acct, c, "search", "General Hattrick search for players, teams, managers, series, arenas, regions or matches.",
		func(ht hattrick.Client, a searchArgs) (any, error) {
			return ht.Search(a.SearchType, a.Query, api.SearchFilters{
				SearchString2:  a.SearchString2,
				SearchID:       a.SearchID,
				SearchLeagueID: a.SearchLeagueID,
				PageIndex:      a.PageIndex,
			})
		})

	// Challenges (read side; arranging a friendly is a write tool).
	add(s, acct, c, "get_challenges", "Friendly challenges you have made and offers from others.",
		func(ht hattrick.Client, a weekendArgs) (any, error) {
			mine, others, err := ht.GetChallenges(a.Weekend)
			return map[string]any{"by_me": mine, "offers": others}, err
		})
	add(s, acct, c, "is_challengeable", "Whether a team can be challenged to a friendly.",
		func(ht hattrick.Client, a challengeableArgs) (any, error) {
			return ht.IsChallengeable(a.Weekend, a.TeamID)
		})
	add(s, acct, c, "are_challengeable", "Whether each of several teams can be challenged to a friendly.",
		func(ht hattrick.Client, a areChallengeableArgs) (any, error) {
			return ht.AreChallengeable(a.Weekend, a.TeamIDs...)
		})
	add(s, acct, c, "get_all_live_events", "All live events for the matches on your live list.",
		func(ht hattrick.Client, a liveEventsArgs) (any, error) {
			return ht.GetAllLiveEvents(a.IncludeStartingLineup, a.UseLiveEventsAndTexts)
		})
	add(s, acct, c, "get_new_live_events", "Live events since the given last-shown indexes.",
		func(ht hattrick.Client, a liveEventsArgs) (any, error) {
			return ht.GetNewLiveEvents(a.LastShownIndexes, a.IncludeStartingLineup, a.UseLiveEventsAndTexts)
		})
}

type areChallengeableArgs struct {
	Weekend bool      `json:"weekend,omitempty" jsonschema:"true for the weekend friendly slot"`
	TeamIDs []id.Team `json:"team_ids" jsonschema:"team IDs to check"`
}

type liveEventsArgs struct {
	LastShownIndexes      string `json:"last_shown_indexes,omitempty" jsonschema:"opaque cursor from a previous call (get_new_live_events only)"`
	IncludeStartingLineup bool   `json:"include_starting_lineup,omitempty" jsonschema:"include the starting lineup"`
	UseLiveEventsAndTexts bool   `json:"use_live_events_and_texts,omitempty" jsonschema:"include live events and texts"`
}

type weekendArgs struct {
	Weekend bool `json:"weekend,omitempty" jsonschema:"true for the weekend friendly slot"`
}

type challengeableArgs struct {
	Weekend bool    `json:"weekend,omitempty" jsonschema:"true for the weekend friendly slot"`
	TeamID  id.Team `json:"team_id" jsonschema:"team ID to check"`
}
