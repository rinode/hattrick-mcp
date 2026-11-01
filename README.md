# hattrick-mcp

An [MCP](https://modelcontextprotocol.io) server that gives AI assistants access to your [Hattrick](https://www.hattrick.org) team through the official CHPP API: the full read surface, plus guarded write actions (match orders, training, transfer bids, youth, challenges).

It runs locally over stdio and is built on [lucianoq/hattrick](https://github.com/lucianoq/hattrick). Pair it with [hattrick-context-pack](https://github.com/rmagasi/hattrick-context-pack) so the assistant knows both your data and the game rules.

## Install as a Claude Code plugin

This repo is a plugin marketplace. Installing `hattrick-mcp` also installs `hattrick-context-pack`, fetched from its upstream repo, so the rules stay current without a copy here. The server is a native binary from the latest release; no runtime is needed.

```sh
claude plugin marketplace add rinode/hattrick-mcp
claude plugin install hattrick-mcp@hattrick-mcp
```

Then open `/plugin`, select `hattrick-mcp` on the Installed tab, choose Configure, and enter your CHPP consumer key and secret. Requires a recent Claude Code; 2.1.193 does not start bundled servers.

## Install in Claude Desktop

Download `hattrick-mcp.mcpb` from the [latest release](https://github.com/rinode/hattrick-mcp/releases/latest), open it, and enter your CHPP consumer key and secret.

## Log in

Ask Claude to log you in to Hattrick. It calls `login`, gives you a URL where you grant access, and finishes with the code you get back. The access token is stored in your user config directory (`~/Library/Application Support/hattrick-mcp/credentials.json` on macOS).

The consumer key comes from an approved [CHPP product license](https://www.hattrick.org/Community/CHPP/). Never commit it.

## Tools

The server exposes the whole CHPP API — 94 tools, every endpoint, with every parameter each method accepts. Read tools collapse the "my" and "by ID" calls into one: omit the ID for your own data. Read results are cached for 5 minutes.

- Manager, team, club, region, achievements
- Players: squad, details, events, transfer history, old players, hall of fame, avatars
- Matches: fixtures, archive, details, lineup, orders, rating prediction, cup matches
- League and world: series, fixtures, division levels, countries, translations
- Economy: finances, fans, supporters
- Training, staff, arena, and their league stats
- Transfers: market search, current bids, team history
- Youth: players, team, league, fixtures
- Tournaments and national teams, including the world cup
- Alliances (federations), bookmarks, ladders, live matches and events, general search
- `search_transfers` takes the full filter set (age, age-days, up to four skill ranges, specialty, native country, TSI, price, paging)

### Write tools

These change the live account. Each is marked destructive and does nothing until called again with `"confirm": true`. The scopes are requested at login.

| Tool | Action |
|---|---|
| `set_lineup` | Build a lineup from player IDs per position, with orders, substitutions, man-marking, captain, set-pieces and penalty takers; without `confirm` it previews predicted ratings |
| `set_match_orders` | Set a lineup from a raw CHPP lineup JSON |
| `set_training` | Change training type, intensity and stamina |
| `place_bid` | Bid on a transfer-listed player (spends money) |
| `ignore_transfer`, `delete_finished_bids` | Manage transfer tracking |
| `set_youth_player_skills_visible`, `set_youth_team_skills_visible` | Reveal youth skill estimates |
| `challenge`, `accept_challenge`, `decline_challenge`, `withdraw_challenge` | Arrange friendlies |
| `add_live_match`, `delete_live_match`, `clear_live_matches` | Manage the live list |

## Development

Docker is only needed for development. Put `HATTRICK_CONSUMER_KEY` and `HATTRICK_CONSUMER_SECRET` in `.env`.

```sh
docker compose run --rm --build test   # gofmt, tidy, vet and tests, as CI runs them
docker compose run --rm --build hattrick-mcp  # the server over stdio
```

With Go 1.26 or later: `go test ./...` and `go build .`.

Pushing a `v*` tag builds the binaries, packs `hattrick-mcp.mcpb` and attaches it to a GitHub release, which the plugin installs from.

`go.mod` replaces `github.com/lucianoq/hattrick` with a fork (`github.com/rinode/hattrick`) that fixes five bugs against live CHPP: the match-orders POST signing and lineup JSON format, the match-orders response `NewPositionId`, transfer-search `ItemCount`, and case-sensitive file-name decoding. The fixes are proposed upstream in [lucianoq/hattrick#1](https://github.com/lucianoq/hattrick/pull/1); drop the `replace` once they land.

## License

Apache License 2.0
