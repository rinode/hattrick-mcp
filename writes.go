package main

import (
	"context"
	"encoding/json"

	"github.com/lucianoq/hattrick"
	"github.com/lucianoq/hattrick/chpp"
	"github.com/lucianoq/hattrick/chpp/id"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type setMatchOrdersArgs struct {
	confirmField
	MatchID      id.Match          `json:"match_id" jsonschema:"match ID"`
	TeamID       id.Team           `json:"team_id" jsonschema:"your team ID for this match"`
	SourceSystem chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
	Lineup       string            `json:"lineup" jsonschema:"lineup JSON in the CHPP matchorders format"`
}

type setTrainingArgs struct {
	confirmField
	TeamID       id.Team           `json:"team_id" jsonschema:"your team ID"`
	TrainingType chpp.TrainingType `json:"training_type" jsonschema:"2 set pieces, 3 defending, 4 scoring, 5 winger, 6 scoring+set pieces, 7 passing, 8 playmaking, 9 keeper, 10 passing (def+mid), 11 defending (def+mid), 12 winger+attackers"`
	Level        uint              `json:"level" jsonschema:"training intensity percentage (0-100)"`
	StaminaLevel uint              `json:"stamina_level" jsonschema:"stamina share percentage (0-100)"`
}

type placeBidArgs struct {
	confirmField
	PlayerID id.Player  `json:"player_id" jsonschema:"player ID to bid on"`
	TeamID   id.Team    `json:"team_id" jsonschema:"your team ID placing the bid"`
	Bid      chpp.Money `json:"bid" jsonschema:"bid amount in SEK"`
	MaxBid   chpp.Money `json:"max_bid,omitempty" jsonschema:"automatic maximum bid in SEK; omit to bid exactly bid"`
}

type ignoreTransferArgs struct {
	confirmField
	TransferID id.Transfer         `json:"transfer_id" jsonschema:"transfer ID to stop tracking"`
	Category   chpp.TrackingTypeID `json:"category" jsonschema:"1 selling, 2 buying, 3 mother club, 4 previous team, 5 hot-listed, 8 losing bids, 9 finished, 10 prospects"`
}

type challengeArgs struct {
	confirmField
	OpponentTeam id.Team           `json:"opponent_team_id" jsonschema:"team ID to challenge"`
	FriendlyType chpp.FriendlyType `json:"friendly_type,omitempty" jsonschema:"0 normal rules (default), 1 cup rules, 12 national team"`
	MatchPlace   chpp.MatchPlace   `json:"match_place,omitempty" jsonschema:"0 home (default), 1 away, 2 neutral"`
	OtherArena   id.Arena          `json:"other_arena_id,omitempty" jsonschema:"arena ID for a neutral match"`
	Weekend      bool              `json:"weekend,omitempty" jsonschema:"true for the weekend friendly slot"`
}

type friendlyMatchArgs struct {
	confirmField
	FriendlyMatchID id.FriendlyMatch `json:"friendly_match_id" jsonschema:"friendly match (challenge) ID"`
}

type unlockYouthPlayerArgs struct {
	confirmField
	YouthPlayerID id.YouthPlayer `json:"youth_player_id" jsonschema:"youth player ID"`
	ShowScoutCall bool           `json:"show_scout_call,omitempty" jsonschema:"include the scout's comments in the result"`
	ShowLastMatch bool           `json:"show_last_match,omitempty" jsonschema:"include the player's last match in the result"`
}

type unlockYouthTeamArgs struct {
	confirmField
	YouthTeamID   id.YouthTeam `json:"youth_team_id,omitempty" jsonschema:"youth team ID; omit for your own"`
	ShowScoutCall bool         `json:"show_scout_call,omitempty" jsonschema:"include the scout's comments in the result"`
	ShowLastMatch bool         `json:"show_last_match,omitempty" jsonschema:"include each player's last match in the result"`
}

type liveAddArgs struct {
	confirmField
	MatchID               id.Match          `json:"match_id" jsonschema:"match ID to add to your live list"`
	SourceSystem          chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
	IncludeStartingLineup bool              `json:"include_starting_lineup,omitempty" jsonschema:"include the starting lineup in the result"`
	UseLiveEventsAndTexts bool              `json:"use_live_events_and_texts,omitempty" jsonschema:"include live events and texts in the result"`
}

type liveDeleteArgs struct {
	confirmField
	MatchID      id.Match          `json:"match_id" jsonschema:"match ID to remove from your live list"`
	SourceSystem chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`
}

type confirmOnly struct {
	confirmField
}

type setLineupArgs struct {
	confirmField
	MatchID      id.Match          `json:"match_id" jsonschema:"match ID"`
	TeamID       id.Team           `json:"team_id" jsonschema:"your team ID for this match"`
	SourceSystem chpp.SourceSystem `json:"source_system,omitempty" jsonschema:"hattrick (default), youth or htointegrated"`

	Goalkeeper           id.Player `json:"goalkeeper,omitempty" jsonschema:"player ID in goal"`
	RightBack            id.Player `json:"right_back,omitempty" jsonschema:"player ID at right wing back"`
	RightCentralDefender id.Player `json:"right_central_defender,omitempty" jsonschema:"player ID at right central defender"`
	CentralDefender      id.Player `json:"central_defender,omitempty" jsonschema:"player ID at middle central defender"`
	LeftCentralDefender  id.Player `json:"left_central_defender,omitempty" jsonschema:"player ID at left central defender"`
	LeftBack             id.Player `json:"left_back,omitempty" jsonschema:"player ID at left wing back"`
	RightWinger          id.Player `json:"right_winger,omitempty" jsonschema:"player ID at right winger"`
	RightInnerMid        id.Player `json:"right_inner_mid,omitempty" jsonschema:"player ID at right inner midfielder"`
	CentralInnerMid      id.Player `json:"central_inner_mid,omitempty" jsonschema:"player ID at middle inner midfielder"`
	LeftInnerMid         id.Player `json:"left_inner_mid,omitempty" jsonschema:"player ID at left inner midfielder"`
	LeftWinger           id.Player `json:"left_winger,omitempty" jsonschema:"player ID at left winger"`
	RightForward         id.Player `json:"right_forward,omitempty" jsonschema:"player ID at right forward"`
	CentralForward       id.Player `json:"central_forward,omitempty" jsonschema:"player ID at middle forward"`
	LeftForward          id.Player `json:"left_forward,omitempty" jsonschema:"player ID at left forward"`

	Behaviours map[string]int `json:"behaviours,omitempty" jsonschema:"optional per-position behaviour by position name: 0 normal, 1 offensive, 2 defensive, 3 towards middle, 4 towards wing, 5 extra forward, 6 extra inner mid, 7 extra defender"`

	BenchGoalkeeper      id.Player `json:"bench_goalkeeper,omitempty" jsonschema:"substitute goalkeeper"`
	BenchCentralDefender id.Player `json:"bench_central_defender,omitempty" jsonschema:"substitute central defender"`
	BenchWingBack        id.Player `json:"bench_wing_back,omitempty" jsonschema:"substitute wing back"`
	BenchInnerMid        id.Player `json:"bench_inner_mid,omitempty" jsonschema:"substitute inner midfielder"`
	BenchForward         id.Player `json:"bench_forward,omitempty" jsonschema:"substitute forward"`
	BenchWinger          id.Player `json:"bench_winger,omitempty" jsonschema:"substitute winger"`
	BenchExtra           id.Player `json:"bench_extra,omitempty" jsonschema:"extra substitute"`

	BackupGoalkeeper      id.Player `json:"backup_goalkeeper,omitempty" jsonschema:"backup substitute goalkeeper"`
	BackupCentralDefender id.Player `json:"backup_central_defender,omitempty" jsonschema:"backup substitute central defender"`
	BackupWingBack        id.Player `json:"backup_wing_back,omitempty" jsonschema:"backup substitute wing back"`
	BackupInnerMid        id.Player `json:"backup_inner_mid,omitempty" jsonschema:"backup substitute inner midfielder"`
	BackupForward         id.Player `json:"backup_forward,omitempty" jsonschema:"backup substitute forward"`
	BackupWinger          id.Player `json:"backup_winger,omitempty" jsonschema:"backup substitute winger"`
	BackupExtra           id.Player `json:"backup_extra,omitempty" jsonschema:"backup extra substitute"`

	Captain        id.Player              `json:"captain,omitempty" jsonschema:"captain player ID"`
	SetPiecesTaker id.Player              `json:"set_pieces_taker,omitempty" jsonschema:"set-pieces taker player ID"`
	PenaltyTakers  []id.Player            `json:"penalty_takers,omitempty" jsonschema:"ordered player IDs who take penalties, best first (up to 11)"`
	Tactic         chpp.MatchTacticType   `json:"tactic,omitempty" jsonschema:"0 normal, 1 pressing, 2 counter-attacks, 3 attack in the middle, 4 attack in wings, 7 play creatively, 8 long shots"`
	Attitude       chpp.MatchTeamAttitude `json:"attitude,omitempty" jsonschema:"-1 play it cool, 0 normal, 1 match of the season"`
	CoachModifier  chpp.CoachModifier     `json:"coach_modifier,omitempty" jsonschema:"team attitude on the defensive(-10)..offensive(10) axis, 0 neutral"`
	ManMarker      id.Player              `json:"man_marker,omitempty" jsonschema:"your player who man-marks an opponent"`
	ManMarking     id.Player              `json:"man_marking,omitempty" jsonschema:"the opponent player being man-marked"`

	Substitutions []subArg `json:"substitutions,omitempty" jsonschema:"up to 5 orders (substitutions, behaviour changes, position swaps)"`
}

type subArg struct {
	PlayerIn  id.Player `json:"player_in" jsonschema:"player coming on (or the affected player for a behaviour/position change)"`
	PlayerOut id.Player `json:"player_out,omitempty" jsonschema:"player going off; equal to player_in for a behaviour or position change"`
	OrderType int       `json:"order_type,omitempty" jsonschema:"1 substitution or behaviour change (default), 3 player swap"`
	Minute    int       `json:"minute,omitempty" jsonschema:"trigger minute: -1 anytime, 46 half-time, 91 before extra time, or a specific minute"`
	Position  int       `json:"position,omitempty" jsonschema:"new position 0-13 (GK..left forward), -1 no change"`
	Behaviour int       `json:"behaviour,omitempty" jsonschema:"new behaviour: -1 no change, 0 normal, 1 offensive, 2 defensive, 3 towards middle, 4 towards wing, 5 extra forward, 6 extra inner mid, 7 extra defender"`
	Card      int       `json:"card,omitempty" jsonschema:"red-card condition: -1 ignore, 1 my player, 2 opponent, 11-15 my CD/mid/fwd/WB/winger, 21-25 opponent variants"`
	Standing  int       `json:"standing,omitempty" jsonschema:"score condition: -1 any, 0 tied, 1 lead, 2 down, 3 lead>1, 4 down>1, 5 not down, 6 not lead, 7 lead>2, 8 down>2, 9 not tied"`
}

func buildLineup(a setLineupArgs) (string, error) {
	b := a.Behaviours
	pb := func(p id.Player, key string) chpp.PlayerIDBehaviour {
		return chpp.PlayerIDBehaviour{Player: p, Behaviour: chpp.MatchBehaviourID(b[key])}
	}
	fs := chpp.LineupFirstStringPlayers{
		Goalkeeper:             pb(a.Goalkeeper, "goalkeeper"),
		RightDefender:          pb(a.RightBack, "right_back"),
		RightCentralDefender:   pb(a.RightCentralDefender, "right_central_defender"),
		CentralDefender:        pb(a.CentralDefender, "central_defender"),
		LeftCentralDefender:    pb(a.LeftCentralDefender, "left_central_defender"),
		LeftDefender:           pb(a.LeftBack, "left_back"),
		RightWinger:            pb(a.RightWinger, "right_winger"),
		RightInnerMidfielder:   pb(a.RightInnerMid, "right_inner_mid"),
		CentralInnerMidfielder: pb(a.CentralInnerMid, "central_inner_mid"),
		LeftInnerMidfielder:    pb(a.LeftInnerMid, "left_inner_mid"),
		LeftWinger:             pb(a.LeftWinger, "left_winger"),
		RightForward:           pb(a.RightForward, "right_forward"),
		CentralForward:         pb(a.CentralForward, "central_forward"),
		LeftForward:            pb(a.LeftForward, "left_forward"),
	}
	bench := chpp.LineupBenchPlayers{
		Primary: chpp.LineupBenchSlots{
			Goalkeeper:      a.BenchGoalkeeper,
			CentralDefender: a.BenchCentralDefender,
			WingBack:        a.BenchWingBack,
			InnerMidfielder: a.BenchInnerMid,
			Forward:         a.BenchForward,
			Winger:          a.BenchWinger,
			Extra:           a.BenchExtra,
		},
		Backup: chpp.LineupBenchSlots{
			Goalkeeper:      a.BackupGoalkeeper,
			CentralDefender: a.BackupCentralDefender,
			WingBack:        a.BackupWingBack,
			InnerMidfielder: a.BackupInnerMid,
			Forward:         a.BackupForward,
			Winger:          a.BackupWinger,
			Extra:           a.BackupExtra,
		},
	}
	var kickers [11]id.Player
	for i := 0; i < len(a.PenaltyTakers) && i < len(kickers); i++ {
		kickers[i] = a.PenaltyTakers[i]
	}
	subs := make([]chpp.LineupSubstitutionOrder, 0, len(a.Substitutions))
	for _, s := range a.Substitutions {
		ot := chpp.MatchOrderType(s.OrderType)
		if s.OrderType == 0 {
			ot = chpp.MatchOrderTypeSubOrBehaviouralChange
		}
		subs = append(subs, chpp.LineupSubstitutionOrder{
			PlayerIn:  s.PlayerIn,
			PlayerOut: s.PlayerOut,
			OrderType: ot,
			Minute:    chpp.LineupMatchMinuteCriteria(s.Minute),
			Position:  chpp.LineupPosition(s.Position),
			Behaviour: chpp.MatchBehaviourID(s.Behaviour),
			Card:      chpp.LineupRedCardCriteria(s.Card),
			Standing:  chpp.LineupGoalDiffCriteria(s.Standing),
		})
	}
	return chpp.NewLineup(fs, bench, kickers, a.Captain, a.SetPiecesTaker,
		a.Tactic, a.Attitude, a.CoachModifier, a.ManMarker, a.ManMarking, subs)
}

func registerWrites(s *mcp.Server, acct *account) {
	destructive := true
	mcp.AddTool(s, &mcp.Tool{
		Name:        "set_lineup",
		Description: `Set the match lineup from player IDs per position. Without "confirm": true it previews the predicted ratings without saving; with it, it saves the orders.`,
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructive},
	}, func(_ context.Context, _ *mcp.CallToolRequest, a setLineupArgs) (*mcp.CallToolResult, any, error) {
		ht, err := acct.client()
		if err != nil {
			return nil, nil, err
		}
		lineup, err := buildLineup(a)
		if err != nil {
			return nil, nil, err
		}
		var res any
		if a.Confirm {
			res, err = ht.SetMatchOrders(a.MatchID, a.TeamID, srcOr(a.SourceSystem), lineup)
		} else {
			res, err = ht.PredictRatings(a.MatchID, a.TeamID, srcOr(a.SourceSystem), lineup)
		}
		if err != nil {
			return nil, nil, err
		}
		out := map[string]any{"saved": a.Confirm, "result": plain(res)}
		jb, err := json.Marshal(out)
		if err != nil {
			return nil, nil, err
		}
		return text(string(jb)), nil, nil
	})

	addWrite(s, acct, "set_match_orders", "Set a match lineup from a raw CHPP matchorders lineup JSON. Prefer set_lineup unless you already have the JSON.",
		func(ht hattrick.Client, a setMatchOrdersArgs) (any, error) {
			return ht.SetMatchOrders(a.MatchID, a.TeamID, srcOr(a.SourceSystem), a.Lineup)
		})
	addWrite(s, acct, "set_training", "Change the training type, intensity and stamina share of your team.",
		func(ht hattrick.Client, a setTrainingArgs) (any, error) {
			changed, tt, err := ht.SetTraining(a.TeamID, a.TrainingType, a.Level, a.StaminaLevel)
			return map[string]any{"changed": changed, "training": tt}, err
		})
	addWrite(s, acct, "place_bid", "Place a bid on a player on the transfer market. Spends real in-game money and is irreversible.",
		func(ht hattrick.Client, a placeBidArgs) (any, error) {
			max := a.MaxBid
			if max == 0 {
				max = a.Bid
			}
			return ht.PlaceBidOnPlayer(a.PlayerID, a.TeamID, a.Bid, max)
		})
	addWrite(s, acct, "ignore_transfer", "Stop tracking a transfer in one of your transfer lists.",
		func(ht hattrick.Client, a ignoreTransferArgs) (any, error) {
			return "ignored", ht.IgnoreTransfer(a.TransferID, a.Category)
		})
	addWrite(s, acct, "delete_finished_bids", "Clear all finished transfers from your bids list.",
		func(ht hattrick.Client, _ confirmOnly) (any, error) {
			return "cleared", ht.DeleteAllFinishedBids()
		})
	addWrite(s, acct, "set_youth_player_skills_visible", "Unlock (reveal) a youth player's hidden skill estimates.",
		func(ht hattrick.Client, a unlockYouthPlayerArgs) (any, error) {
			return ht.UnlockYouthPlayerSkills(a.YouthPlayerID, a.ShowScoutCall, a.ShowLastMatch)
		})
	addWrite(s, acct, "set_youth_team_skills_visible", "Unlock (reveal) hidden skill estimates for your whole youth team.",
		func(ht hattrick.Client, a unlockYouthTeamArgs) (any, error) {
			if a.YouthTeamID == 0 {
				d, _, err := ht.GetMyYouthTeam(false)
				if err != nil {
					return nil, err
				}
				a.YouthTeamID = d.ID
			}
			return ht.UnlockYouthTeamSkills(a.YouthTeamID, a.ShowScoutCall, a.ShowLastMatch)
		})
	addWrite(s, acct, "challenge", "Arrange a friendly match against another team.",
		func(ht hattrick.Client, a challengeArgs) (any, error) {
			return "challenged", ht.Challenge(a.OpponentTeam, a.FriendlyType, a.MatchPlace, a.OtherArena, a.Weekend)
		})
	addWrite(s, acct, "accept_challenge", "Accept a friendly challenge offered to you.",
		func(ht hattrick.Client, a friendlyMatchArgs) (any, error) {
			return "accepted", ht.AcceptChallenge(a.FriendlyMatchID)
		})
	addWrite(s, acct, "decline_challenge", "Decline a friendly challenge offered to you.",
		func(ht hattrick.Client, a friendlyMatchArgs) (any, error) {
			return "declined", ht.DeclineChallenge(a.FriendlyMatchID)
		})
	addWrite(s, acct, "withdraw_challenge", "Withdraw a friendly challenge you arranged.",
		func(ht hattrick.Client, a friendlyMatchArgs) (any, error) {
			return "withdrawn", ht.WithdrawChallenge(a.FriendlyMatchID)
		})
	addWrite(s, acct, "add_live_match", "Add a match to your live-tracking list.",
		func(ht hattrick.Client, a liveAddArgs) (any, error) {
			return ht.AddLiveMatch(a.MatchID, srcOr(a.SourceSystem), a.IncludeStartingLineup, a.UseLiveEventsAndTexts)
		})
	addWrite(s, acct, "delete_live_match", "Remove a match from your live-tracking list.",
		func(ht hattrick.Client, a liveDeleteArgs) (any, error) {
			return ht.DeleteLiveMatch(a.MatchID, srcOr(a.SourceSystem))
		})
	addWrite(s, acct, "clear_live_matches", "Remove all matches from your live-tracking list.",
		func(ht hattrick.Client, _ confirmOnly) (any, error) {
			return "cleared", ht.ClearLiveMatches()
		})
}
