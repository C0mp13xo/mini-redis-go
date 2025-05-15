package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	cache := NewCache(1000)
	done := make(chan struct{})
	go cache.cleanupExpiredKeys(60*time.Second, done)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		close(done)
		fmt.Println("shutting down.... ")
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
			cache.Set(key, value, ttl)
			fmt.Fprintf(w, "set key=%s value =%s ttl=%v \n", key, value, ttl)
			return
		}
		http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/cache/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/cache/"):]
		switch r.Method {
		case http.MethodGet:
			if val, found := cache.Get(key); !found {
				http.Error(w, "key not found or expired ", http.StatusNotFound)
			} else {
				fmt.Fprintf(w, "value: %s \n", val)
			}
		case http.MethodDelete:
			if cache.Del(key) {
				fmt.Fprintf(w, "deleted key: %s \n", key)
			} else {
				http.Error(w, "key not found ", http.StatusNotFound)
			}
		default:
			http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hits: %d\n, misses: %d\n, evictions: %d\n, expired: %d\n", cache.hits, cache.misses, cache.evictions, cache.expired)
	})

	fmt.Println("Mini-Redis server running on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
