# IsPlexUpToDate

Compares the running Plex Media Server version against the latest version published by Plex, and reports the result to a Gatus external endpoint. Meant to run as a Kubernetes `CronJob`.

## How it works

1. `GET {PLEX_URL}/identity` (with `Accept: application/json`) to get the local version.
2. `GET https://plex.tv/api/downloads/5.json` to get the latest public version (Linux, first item in `computer.Linux.releases`).
3. Compares both versions.
4. `POST {GATUS_URL}?success=...&error=...` with `Authorization: Bearer {GATUS_APIKEY}`.

If either version can't be fetched, the check fails immediately (`log.Fatal`, exit code 1) after attempting to report the failure to Gatus. An available update is **not** treated as a process error: the program exits normally (exit 0), the `success=false` status and message are carried by Gatus.

## Configuration

| Variable | Example | Description |
|---|---|---|
| `PLEX_URL` | `http://plex.media.svc.cluster.local:32400` | Plex server base URL (without `/identity`) |
| `GATUS_URL` | `http://gatus.monitoring.svc.cluster.local:8080/api/v1/endpoints/xxxxx/external` | Full URL of the Gatus external endpoint |
| `GATUS_APIKEY` | `PZUFRZ7Zozryfgou` | Token sent in the `Authorization: Bearer` header |

All three are required — the program exits early otherwise (`Missing Env variables.`).

## Build

```bash
go build -o isplexuptodate .
```

Docker image: see `Dockerfile` (multi-stage, multi-arch amd64/arm64, code expected under `./src`).

## Run

```bash
PLEX_URL=http://plex.media.svc:32400 \
GATUS_URL=http://gatus.monitoring.svc:8080/api/v1/endpoints/xxxxx/external \
GATUS_APIKEY=PZUFRZ7Zozryfgou \
./isplexuptodate
```

## Deployment

Meant to run as a Kubernetes `CronJob`, with `GATUS_APIKEY` stored in a `Secret` rather than hardcoded in the manifest.

## Known limitations

- No custom TLS handling: if `PLEX_URL` is `https://` with a self-signed cert, the default Go client will reject the connection.
- Version comparison is a strict string match — a server running a build with a different suffix (e.g. LinuxServer.io's `-lsXXX`) will never match the official Plex feed, even when up to date.
