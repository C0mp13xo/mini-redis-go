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
	cache := NewCache()
	done := make(chan struct{})
	go cache.cleanupExpiredKeys(done)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		close(done)
	}()

	fmt.Println("shutting down.... ")
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")

		if val, found := cache.Get(key); !found {
			http.Error(w, "key not found or expired ", http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "value: %s \n", val)
		}

	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if cache.Del(key) {
			fmt.Fprintf(w, "deleted key: %s \n", key)
		} else {
			http.Error(w, "key not found ", http.StatusNotFound)
		}
	})
	fmt.Println("Mini-Redis server running on :8080")
	http.ListenAndServe(":8080", nil)
}
