# hattrick-mcp

An [MCP](https://modelcontextprotocol.io) server that gives AI assistants read-only access to your [Hattrick](https://www.hattrick.org) team through the official CHPP API.

It runs locally over stdio and is built on [lucianoq/hattrick](https://github.com/lucianoq/hattrick). Pair it with [hattrick-context-pack](https://github.com/rmagasi/hattrick-context-pack) so the assistant knows both your data and the game rules.

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

Omit an ID to use your own team. Results are cached for 5 minutes.

## Build

Requires Go 1.26 or later.

```sh
go build -ldflags "-X main.consumerKey=$HATTRICK_CONSUMER_KEY -X main.consumerSecret=$HATTRICK_CONSUMER_SECRET" -o hattrick-mcp .
```

Without `-ldflags`, set `HATTRICK_CONSUMER_KEY` and `HATTRICK_CONSUMER_SECRET` at runtime. The consumer key comes from an approved [CHPP product license](https://www.hattrick.org/Community/CHPP/). Never commit it.

## Log in

```sh
./hattrick-mcp login
```

Open the printed URL, grant access and paste the code. The access token is stored in your user config directory (`~/Library/Application Support/hattrick-mcp/credentials.json` on macOS).

## Connect to Claude

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "hattrick": {
      "command": "/absolute/path/to/hattrick-mcp",
      "args": ["serve"]
    }
  }
}
```

Claude Code: `claude mcp add hattrick -- /absolute/path/to/hattrick-mcp serve`

## Docker

Put `HATTRICK_CONSUMER_KEY` and `HATTRICK_CONSUMER_SECRET` in `.env`, then log in once. The access token is stored in the `credentials` volume.

```sh
docker compose run --rm hattrick-mcp login
```

Connect to Claude with `-T`, so stdin carries the MCP protocol without a TTY:

```json
{
  "mcpServers": {
    "hattrick": {
      "command": "docker",
      "args": ["compose", "-f", "/absolute/path/to/hattrick-mcp/compose.yaml", "run", "--rm", "-T", "hattrick-mcp"]
    }
  }
}
```

Pushing a `v*` tag publishes `ghcr.io/rinode/hattrick-mcp` for amd64 and arm64. The image contains no consumer key.

## Test

```sh
go test ./...
```

Or in Docker, as CI runs it: `docker compose run --rm test`.

## License

Apache License 2.0
