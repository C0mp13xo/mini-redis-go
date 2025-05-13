📦 Mini Redis-like In-Memory LRU Cache (Go)
This project is a simplified Redis-like in-memory key-value store written in Golang.

It supports:

✅ Basic Set, Get, Delete operations

✅ TTL (Time-To-Live) expiration for keys

✅ LRU Eviction Policy for limited capacity

✅ Periodic cleanup of expired keys

🗂️ Current Design
Data Structures:
Component	Type	Purpose
map[string]*list.Element	Map	O(1) lookup of keys
container/list.List	Doubly Linked List	Tracks usage order (front = MRU, back = LRU)
CacheItem struct	Struct	Holds Value, Key, and Expiry time

LRU Eviction:
On Get or Set: Move key to front of the list (MRU).

On overflow: Evict Back() element (LRU).

Expired keys cleaned in background goroutine.

TTL Cleanup:
Background goroutine using time.Ticker cleans expired keys every N seconds.

🛠️ Usage
Run Server:
bash
Copy
Edit
go run main.go
HTTP Endpoints:
Method	URL	Params	Description
GET	/set	key, value, ttl (seconds)	Set key with value and optional TTL
GET	/get	key	Get value of key
GET	/delete	key	Delete a key

❗ Currently, this is not following strict RESTful design. Will be refactored.

📊 ✅ Next Features (Planned):
Feature	Description
🟢 Cache Metrics	Add metrics: cache hits, misses, evictions, TTL expirations
🟢 High Concurrency Optimization	Sharded locks / sync.Map / concurrency-safe structures
🟢 Prometheus Metrics	Expose metrics endpoint for monitoring
🟢 RESTful API Refactor	Replace /set, /get, /delete with RESTful design (POST /cache, GET /cache/{key}, DELETE /cache/{key})