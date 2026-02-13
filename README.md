# ccut

A terminal UI for visualizing your Claude Code usage stats. Reads directly from `~/.claude/` -- no auth required.

<!-- TODO: add screenshot -->

## Features

**Dashboard** -- Today's stats (messages, sessions, tool calls, tokens), lifetime totals, model usage breakdown, daily activity bar chart, and hourly heatmap.

**Sessions** -- Browse recent sessions by project. Select one to see full token breakdown, tools used, git branch, and duration. Filter with `/`, group by project with `g`.

**Live Monitor** -- Watches `stats-cache.json` in real time via fsnotify. Shows delta indicators (e.g. `Messages: 5745 (+12)`) and a "last updated" timer.

## Installation

```
go install github.com/prs/ccut@latest
```

Or build from source:

```
git clone https://github.com/prs/ccut.git
cd ccut
go build -o ccut
```

## Usage

```
ccut
```

## Key Bindings

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Next / previous tab |
| `1` `2` `3` | Jump to Dashboard, Sessions, Live |
| `j` / `k` or arrows | Navigate / scroll |
| `Enter` | Select session |
| `Esc` | Back to list |
| `/` | Filter sessions |
| `g` | Group by project |
| `r` | Force refresh (live) |
| `q` / `Ctrl+C` | Quit |

## Data Sources

All data is read from `~/.claude/`:

| File | Purpose |
|------|---------|
| `stats-cache.json` | Aggregated daily activity, model tokens, hourly counts, totals |
| `projects/<project>/<session>.jsonl` | Per-session message logs with token usage |
| `.credentials.json` | Subscription type and rate limit tier |

## Requirements

- Go 1.21+
