# IsPlexUpToDate

Compares the running Plex Media Server version against the latest version published by Plex, and if the Plex instance need to be updated the result is pushed to Discord. Meant to run as a Kubernetes `CronJob`.

## How it works

1. `GET {PLEX_URL}/identity` (with `Accept: application/json`) to get the local version.
2. `GET https://plex.tv/api/downloads/5.json` to get the latest public version (Linux, first item in `computer.Linux.releases`).
3. Compares both versions.
4. `POST {DISCORD_WEBHOOK_URL}`.

## Configuration

| Variable | Example | Description |
|---|---|---|
| `PLEX_URL` | `http://plex.media.svc.cluster.local:32400` | Plex server base URL (without `/identity`) |
| `DISCORD_WEBHOOK_URL` | `https://discord.com/api/webhooks/...` | Discord webhook URL for notifications (optional) |

All required variables are required — the program exits early otherwise (`Missing Env variables.`).

## Build

```bash
go build -o isplexuptodate ./cmd/isplexuptodate
```

Docker image: see `Dockerfile` (multi-stage, multi-arch amd64/arm64, code expected under `./cmd` and `./internal`).

## Run

```bash
PLEX_URL=http://plex.media.svc:32400 \
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/... \
./isplexuptodate
```

## Deployment

Meant to run as a Kubernetes `CronJob`, with `DISCORD_WEBHOOK_URL` stored in a `Secret` rather than hardcoded in the manifest.

## Known limitations

- No custom TLS handling: if `PLEX_URL` is `https://` with a self-signed cert, the default Go client will reject the connection.
- Version comparison is a strict string match — a server running a build with a different suffix (e.g. LinuxServer.io's `-lsXXX`) will never match the official Plex feed, even when up to date.
