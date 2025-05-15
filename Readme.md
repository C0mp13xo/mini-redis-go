# Mini-Go-Redis Cache Server

A simple in-memory LRU cache with TTL expiration implemented in Go, supporting persistence via JSON snapshot files.

## Features

- Thread-safe LRU cache with configurable capacity.
- TTL-based expiration of cache entries.
- Background cleanup of expired keys.
- HTTP API for setting, getting, deleting keys.
- Metrics endpoint showing hits, misses, evictions, and expired counts.
- Persistence by saving cache snapshot to disk on shutdown and loading on startup.

## Installation

Make sure you have Go installed (>= 1.22).

Clone the repo and build:

```bash
go build -o mini-go-redis


API Endpoints
POST /cache?key=KEY&value=VALUE&ttl=SECONDS

Set a key with optional TTL in seconds.

GET /cache/KEY

Retrieve value for a key.

DELETE /cache/KEY

Delete a key.

GET /metrics

Show cache statistics.

