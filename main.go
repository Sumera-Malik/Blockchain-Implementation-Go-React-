package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	BlockChainDisplayName = "Sumera Malik Blockchain"
	OwnerName             = "Sumera Malik"
	FirstRollNumber       = "21i-1579"
)

var BC Blockchain

func main() {
	// Initialize blockchain
	BC = Blockchain{
		Name:       OwnerName,
		Difficulty: 4, // adjust if mining is too slow/fast
	}
	BC.Blocks = append(BC.Blocks, BC.createGenesis())

	// Requirement: first transaction must be roll number
	BC.AddTransaction(FirstRollNumber)
	// Mine initial block so your roll number is recorded on-chain
	if _, err := BC.MinePending(); err != nil {
		log.Fatalf("failed to mine initial block: %v", err)
	}

	http.HandleFunc("/name", withCORS(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"name": BC.Name})
	}))

	http.HandleFunc("/view", withCORS(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"blocks":      BC.Blocks,
			"difficulty":  BC.Difficulty,
			"pendingTx":   BC.PendingTx,
			"displayName": BlockChainDisplayName,
		})
	}))

	http.HandleFunc("/pending", withCORS(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"pending": BC.PendingTx})
	}))

	http.HandleFunc("/tx", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var body struct{ Data string `json:"data"` }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		BC.AddTransaction(body.Data)
		json.NewEncoder(w).Encode(map[string]any{
			"ok":          true,
			"added":       body.Data,
			"pendingSize": len(BC.PendingTx),
		})
	}))

	http.HandleFunc("/mine", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		blk, err := BC.MinePending()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "block": blk})
	}))

	http.HandleFunc("/search", withCORS(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		type hit struct {
			BlockIndex int      `json:"blockIndex"`
			Matches    []string `json:"matches"`
			Hash       string   `json:"hash"`
		}
		var results []hit
		if q != "" {
			for _, b := range BC.Blocks {
				var matches []string
				for _, d := range b.Data {
					if containsCaseInsensitive(d, q) {
						matches = append(matches, d)
					}
				}
				if len(matches) > 0 {
					results = append(results, hit{BlockIndex: b.Index, Matches: matches, Hash: b.Hash})
				}
			}
		}
		json.NewEncoder(w).Encode(map[string]any{"query": q, "results": results})
	}))

	fmt.Println("🚀 Backend running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- helpers (CORS + case-insensitive contains) ---

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

func containsCaseInsensitive(s, sub string) bool {
	// simple case-folding search
	sLow, subLow := toLower(s), toLower(sub)
	return contains(sLow, subLow)
}

func toLower(s string) string {
	out := []rune(s)
	for i, r := range out {
		if 'A' <= r && r <= 'Z' {
			out[i] = r + ('a' - 'A')
		}
	}
	return string(out)
}

func contains(s, sub string) bool {
	// naive search to avoid extra imports
	if len(sub) == 0 { return true }
	if len(sub) > len(s) { return false }
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
