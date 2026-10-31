# hattrick-mcp

An [MCP](https://modelcontextprotocol.io) server that gives AI assistants read-only access to your [Hattrick](https://www.hattrick.org) team through the official CHPP API.

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

| Tool | Data |
|---|---|
| `get_team` | Team details, spirit, confidence |
| `get_players` | Senior squad |
| `get_player` | One player in full |
| `get_matches` | Recent and upcoming matches |
| `get_match` | Ratings, chances, tactics, events |
| `get_match_lineup` | Lineup of a played match |
| `get_training` | Training type, intensity, history |
| `get_economy` | Finances |
| `get_series` | League table |
| `get_series_fixtures` | Series fixtures per season |
| `get_staff` | Trainer and staff |
| `get_arena` | Arena |
| `get_youth_players` | Youth academy |
| `search_transfers` | Transfer market search |
| `login`, `complete_login` | Log in to Hattrick |

Omit an ID to use your own team. Results are cached for 5 minutes.

## Development

Docker is only needed for development. Put `HATTRICK_CONSUMER_KEY` and `HATTRICK_CONSUMER_SECRET` in `.env`.

```sh
docker compose run --rm --build test   # gofmt, tidy, vet and tests, as CI runs them
docker compose run --rm --build hattrick-mcp  # the server over stdio
```

With Go 1.26 or later: `go test ./...` and `go build .`.

Pushing a `v*` tag builds the binaries, packs `hattrick-mcp.mcpb` and attaches it to a GitHub release, which the plugin installs from.

## License

Apache License 2.0
