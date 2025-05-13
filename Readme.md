This is a simplified Redis-like in-memory cache service implemented in Golang, aimed to help you understand the internals of caching systems like Redis.
We're building this from scratch, step-by-step, towards an enterprise-level scalable architecture.

✅ Current Features
 In-Memory Key-Value Cache

 Expiration Support (TTL)

 Background Cleanup of Expired Keys using time.Ticker

 Graceful Shutdown Handling via OS signals (SIGINT, SIGTERM)

 HTTP Server exposing cache endpoints (planned)

✅ Architecture Overview (So Far)
scss
Copy
Edit
main.go
│
├── Starts Cache & Cleanup goroutine
├── Runs HTTP Server (ListenAndServe)
├── Handles graceful shutdown signals (Ctrl+C / SIGTERM)
│
└── cleanupExpiredKeys()
    └── Runs every N seconds using time.Ticker
        └── Iterates keys, deletes expired ones
✅ Example: Graceful Shutdown Flow
User presses Ctrl+C.

Signal handler goroutine triggers shutdown:

Stops cleanupExpiredKeys ticker.

Gracefully shuts down HTTP server.

Application exits cleanly.

✅ Key Concepts Covered
Concept	Why it's Important
time.Ticker	Efficient periodic tasks without tight loops.
Goroutines	Concurrency handling for cleanup & signal handling.
OS Signal Handling	For graceful termination of services.
Channel-based Coordination	Clean shutdown signals between routines.

✅ Why This Approach?
This mimics how Redis periodically evicts expired keys (lazy & active expiry mechanisms).

Helps understand background tasks vs user-facing request handling.

Teaches you how to design long-running services with graceful shutdown.

🗺️ Roadmap (Planned Features)
Feature	Description
✅ TTL-based Expiry (Done)	Per-key expiration timeouts.
🟡 Basic HTTP API (GET, SET, DELETE)	Expose cache endpoints via REST.
🟡 LRU Eviction Policy	Evict least recently used keys when memory limit is reached.
🟡 In-memory Size Limits	Configurable max keys or memory usage.
🟡 Persistent Storage (AOF/RDB Simulation)	Snapshot and append-only persistence like Redis.
🟡 Statistics & Metrics Endpoint	Expose ops/sec, hits/misses, memory usage via /metrics.
🟡 Sharding / Partitioning Simulation	For scaling beyond 10M keys across instances.
🟡 Pub/Sub Skeleton	Simulate Redis pub/sub messaging.
🟡 Cluster Mode / Replica Simulation	Design for high-availability awareness.

🚀 Goal
➡ Educational, but realistic
➡ Designed to help you learn:

How Redis-like caches are built.

How to scale them theoretically to 10M+ keys.

How caches integrate into real-world products.

🏗️ Tech Stack
Go 1.21+

Standard Library Only (for now)

✅ Running the Project
bash
Copy
Edit
go run main.go
You'll see cleanup ticks in the console.
HTTP API endpoints will come next.

✅ To Be Done Tomorrow:
Add HTTP API endpoints (GET, SET, DELETE).

Add LRU eviction (core data structure design).

Metrics (hit/miss counters, key count, memory stats).

Scalability design notes (how Redis scales to millions of keys).

✅ Final Goal: Redis Simplified in Go
By end of this project, you'll have:

An actual Redis-like cache service in Go.

Theoretical understanding of production-scale caching.

Hands-on code for expiry, eviction, persistence, clustering basics.

📂 Folder Structure (Soon)
bash
Copy
Edit
cache/
    cache.go      # Core cache logic (map, expiry, eviction)
    eviction.go   # LRU / LFU implementations
    persistence.go# AOF / RDB simulation
server/
    http.go       # HTTP server, endpoints
    metrics.go    # Stats & metrics
main.go           # Entry point, signal handling, cleanup routines
README.md         # This file
✅ Contributions (Your Learning Journey)
This is designed for personal learning. Fork, modify, and experiment!