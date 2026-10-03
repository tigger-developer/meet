---
title: Architecture
version: 0.5
last-updated: 2026-10-03
---

# Architecture

## Overview

meet is a single Go binary with several subcommands (`serve`, `token`,
`moderator-link`, `create`, `cancel`, `list`) and a companion Secure Shell
(SSH) wrapper (`meet-helper`) for
invoking any subcommand on a remote deploy host. The server embeds all
static assets (HTML template, font) and requires no external runtime
assets.

## Components

```
                     ┌─────────────┐
                     │   Caddy     │ TLS termination, reverse proxy
                     │   :443      │
                     └──────┬──────┘
                            │
                     ┌──────▼──────┐
                     │  meet serve │ Go HTTP server
                     │  :18085     │
                     └──────┬──────┘
                            │
              ┌─────────────┼─────────────┐
              │             │             │
      ┌───────▼──────┐ ┌───▼────┐  ┌─────▼────────┐
      │ Room, login, │ │ Timer  │  │ Recording    │
      │ static pages │ │ HTTP/SSE│  │ webhook      │
      └──────────────┘ └────────┘  └─────┬────────┘
                                         │
                         ┌───────────────┴──────────────┐
                         │                              │
                   ┌─────▼─────┐                 ┌──────▼─────┐
                   │Cloudflare │                 │ Nextcloud  │
                   │ Stream    │                 │ WebDAV     │
                   └───────────┘                 └────────────┘
```

## Request flow

### Meeting page (`GET /{room}`)

1. Accept only a single-segment room path.
2. Admit the request when the latest append-only room-registry row has an
   active occurrence, or when a valid moderator JWT authorizes the room.
3. Return the same non-storable inactive-room HTTP 404 for every blocked slug.
4. Render the embedded meeting template with the room, application ID, and
   branded domain.
5. Pass an optional moderator JWT to the 8x8 client.

`/{room}/moderator` provides a generic login form only while the room is
active. A preapproved email and room pair receives a signed, expiring,
single-use verification link. Verification produces a short-lived JWT bound to
that room. The CLI-generated operator token remains a separate wildcard path.

### Shared timer

Each room has independent in-memory runtime state. Timer settings persist in an
append-only state file, but active runs reset on service restart. Participants
subscribe through server-sent events (SSE); moderator-scoped or wildcard JWTs
authorize control requests. The server emits time-based cues once per run and
re-broadcasts authoritative state every ten seconds.

The browser owns cue playback and temporary local microphone muting. It waits
for microphone mute confirmation before sound, and restores on playback
completion without fixed padding. Concurrent cues share mute ownership until
the last playback finishes. Already-muted microphones and later microphone
state changes are respected. Playback errors release ownership; missing mute
confirmation skips playback and logs an error instead of guessing the state.
The external API cannot distinguish an unchanged mute state reasserted by a
moderator. Sound mapping and live validation are defined in
[W022 - Timer sounds](../specs/W022-timer-sounds/spec.org).

### Webhook (`POST /webhook/recording`)

1. Validate `Authorization` header against configured token
2. Parse JSON payload, log event metadata (all authenticated events are logged)
3. For download events (`RECORDING_UPLOADED`, `TRANSCRIPTION_UPLOADED`,
   `CHAT_UPLOADED`):
   - Check idempotency key against dedup map
   - Respond HTTP 200 immediately
   - Spawn a goroutine for the download/upload pipeline.
   - Download the file to `{STATE_DIRECTORY}/download/{filename}`.
   - Upload video to Cloudflare Stream. Record the playback URL in
     `recordings.csv` and write a Markdown notification to Nextcloud.
   - Upload transcripts and chat logs directly to Nextcloud through WebDAV.
   - Move successful local files to `{STATE_DIRECTORY}/uploaded/`.
   - Retry failures with exponential backoff for up to 24 hours; an exhausted
     file remains in `download/` for recovery.
4. For all other events: log and respond HTTP 200

### Token generation (`meet token`)

1. Load config (defaults + host + secrets)
2. Parse RSA private key from PEM
3. Build JWT with moderator claims and recording feature flag
4. Print `{base_url}/{room}?jwt={signed_token}` to stdout

## Config layers

Config is merged left-to-right from three sources:

1. `config/defaults.yaml` - always loaded, committed, universal baseline
2. `CONFIG_PATH` env var - host-specific overrides (addr, base_url)
3. `SECRETS_PATH` env var - age-encrypted secrets (keys, credentials)

The `--config` flag overrides this when passed explicitly. This convention
is shared with writeback and golink.

Later files override earlier ones. YAML unmarshalling is additive - fields
not present in a later file retain their earlier values.

## State

The server uses two forms of state:

**In-memory:** the dedup map (bounded at 1000 entries, evicts oldest) prevents
reprocessing of duplicate webhook deliveries within a single run. A restart
clears it. Per-room active timer runs also live in memory and reset on restart.

**On-disk:** the `STATE_DIRECTORY` (systemd's `/var/lib/meet/`) contains two
subdirectories that form a simple state machine:

- `rooms.csv` - append-only room lifecycle and schedule definitions.
- `recordings.csv` - successful Cloudflare Stream uploads and playback URLs.
- `timer-settings.csv` - persisted per-room timer configuration.
- `download/` - files downloaded from 8x8 but not yet uploaded.
- `uploaded/` - successfully uploaded files retained for the configured local
  retention period and purged by a daily ticker.

On startup, the server scans `download/` and retries any pending uploads.
This recovers from crashes, restarts, and deployment-induced service restarts.

## Security

- The webhook endpoint validates a bearer token configured in secrets.
- Moderator JWTs are RS256-signed with a private key stored in secrets.
- Magic-link requests use room-specific email allowlists, generic responses,
  expiring signatures, and replay prevention.
- Blocked room slugs share one byte-identical response to avoid room-name
  enumeration.
- No secrets are stored in committed config files; deployments use
  age-encrypted YAML.
- Caddy handles TLS termination; the app binds to loopback only
- systemd runs the service with `DynamicUser=yes` and aggressive sandboxing

## Changelog

- **0.5, 2026-10-03:** Described playback-bound microphone ownership and its
  external API limitation.

- **0.4, 2026-10-02:** Reconciled routing, moderator authentication, timer
  state, and Cloudflare Stream recording archival with the migrated legacy
  acceptance criteria.
