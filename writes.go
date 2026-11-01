package main

import (
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
	TrainingType chpp.TrainingType `json:"training_type" jsonschema:"2 set pieces, 3 defending, 4 scoring, 5 winger, 6 scoring+set pieces, 7 playmaking, 8 goalkeeping, 9 passing, 10 crossing, 11 shooting, 12 through passes, 13 defensive positions"`
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

func registerWrites(s *mcp.Server, acct *account) {
	addWrite(s, acct, "set_match_orders", "Set the lineup, positions, orders and substitutions for an upcoming match.",
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
