package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

// NodeID is the unique identifier for a real node.
type NodeID string

// HashRing represents the consistent hash ring.
type HashRing struct {
	mu           sync.RWMutex
	nodes        map[string]NodeID
	sortedHashes []string
	virtualNodes int
}

func NewHashRing(vnodes int) *HashRing {
	return &HashRing{
		nodes:        make(map[string]NodeID),
		sortedHashes: []string{},
		virtualNodes: vnodes,
	}
}

func (r *HashRing) AddNode(node NodeID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < r.virtualNodes; i++ {
		hash := hashKey(fmt.Sprintf("%s-%d", node, i))
		r.nodes[hash] = node
		r.sortedHashes = append(r.sortedHashes, hash)
	}
	sort.Strings(r.sortedHashes)
}

func (r *HashRing) GetReplicas(key string, n int) []NodeID {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h := hashKey(key)
	start := sort.SearchStrings(r.sortedHashes, h)
	var replicas []NodeID
	for i := 0; len(replicas) < n && i < len(r.sortedHashes); i++ {
		idx := (start + i) % len(r.sortedHashes)
		node := r.nodes[r.sortedHashes[idx]]
		if !contains(replicas, node) {
			replicas = append(replicas, node)
		}
	}
	return replicas
}

func contains(nodes []NodeID, node NodeID) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}

func hashKey(key string) string {
	h := sha1.Sum([]byte(key))
	return hex.EncodeToString(h[:])
}
