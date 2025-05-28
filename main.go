package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	hashRing *HashRing
	caches   map[NodeID]*Cache
	// thisNode NodeID = "NodeA" // Simulate that this instance is NodeA
)

func main() {
	caches = make(map[NodeID]*Cache)
	hashRing = NewHashRing(100)

	nodes := []NodeID{"NodeA", "NodeB", "NodeC"}
	for _, node := range nodes {
		hashRing.AddNode(node)
		caches[node] = NewCache(100)
	}

	// if err := cache.LoadSnapshot("cache_snapshot.json"); err != nil {
	// 	fmt.Println("No Snapshot Loaded !!!!")
	// }
	done := make(chan struct{})
	// go cache.cleanupExpiredKeys(60*time.Second, done)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		<-sigChan
		close(done)
		// if err := cache.SaveSnapshot("cache_snapshot.json"); err != nil {
		// 	fmt.Println("Error Saving Snapshot!!!")
		// } else {
		// 	fmt.Println("Snapshot saved successfully!!!")
		// }
		fmt.Println("shutting down.... ")
		os.Exit(0)
	}()
	fmt.Println("starting server ")
	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			key := r.URL.Query().Get("key")
			value := r.URL.Query().Get("value")
			ttlStr := r.URL.Query().Get("ttl")
			var ttl time.Duration
			if ttlStr != "" {
				if sec, err := strconv.Atoi(ttlStr); err == nil {
					ttl = time.Duration(sec) * time.Second
				}
			}
			nodes := hashRing.GetReplicas(key, 5)
			for _, node := range nodes {
				fmt.Printf("[POST] Writing key '%s' to node %s\n", key, node)
				caches[node].Set(key, value, ttl)
			}
			fmt.Println("total nodes are ", len(hashRing.sortedHashes))
			return
		}

		http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/cache/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/cache/"):]
		nodes := hashRing.GetReplicas(key, 100)
		switch r.Method {
		case http.MethodGet:
			for _, node := range nodes {
				if val, found := caches[node].Get(key); found {
					fmt.Fprintf(w, "value (from %s): %s\n", node, val)

				}
			}
			http.Error(w, "key not found or expired", http.StatusNotFound)

		case http.MethodDelete:
			deleted := false
			for _, node := range nodes {
				if caches[node].Del(key) {
					deleted = true
				}
			}
			if deleted {
				fmt.Fprintf(w, "deleted key: %s\n", key)
			} else {
				http.Error(w, "key not found ", http.StatusNotFound)
			}

		default:
			http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		metrics := make(map[string]map[string]int64)

		for _, node := range hashRing.nodes {
			cache, exists := caches[node]
			if !exists {
				continue // or log missing node
			}

			cache.mu.RLock()
			metrics[string(node)] = map[string]int64{
				"hits":      cache.hits,
				"misses":    cache.misses,
				"evictions": cache.evictions,
				"expired":   cache.expired,
			}
			cache.mu.RUnlock()
		}
		json.NewEncoder(w).Encode(metrics)

	})

	fmt.Println("Mini-Redis server running on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
