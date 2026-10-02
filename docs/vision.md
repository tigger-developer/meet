---
title: Vision
version: 0.2
last-updated: 2026-10-02
---

# Vision

meet is a self-hosted, branded wrapper around 8x8 JaaS (Jitsi as a Service)
for small-scale video meetings. The goal is a simple, scriptable meeting
tool with full control over branding, authentication, and data retention.

## Goals

- **Minimal surface area.** A single Go binary serves scheduled meeting rooms,
  handles webhooks, manages room and timer state, and generates moderator
  tokens. It uses append-only files rather than a database.
- **Own your data.** Recordings, transcriptions, and chat logs are
  automatically downloaded before the 8x8 link expires. Video is retained in
  Cloudflare Stream, while transcripts, chat logs, and recording-link
  notifications are retained in a self-hosted Nextcloud instance.
- **CLI-first administration.** Moderator access is generated via CLI
  (`meet token`) or delivered to a preapproved room moderator through a
  room-scoped magic link. There is no web administration panel.
- **Scriptable and composable.** Room names are URL paths. Token generation
  and room scheduling are CLI operations. The webhook and timer interfaces use
  standard HTTP. Everything integrates with shell scripts and automation.

## Non-goals

- Multi-tenant or multi-user admin (single operator assumed)
- Custom video infrastructure (delegates to 8x8 JaaS)
- Mobile apps (the web UI is responsive via JaaS)
- General user accounts or a multi-user administration system

## Changelog

- **0.2, 2026-10-02:** Reconciled the vision with scheduled rooms,
  room-scoped moderator access, the shared timer, and Cloudflare Stream video
  archival.
