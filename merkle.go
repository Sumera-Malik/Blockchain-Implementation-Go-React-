package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func merkleRoot(tx []string) string {
	if len(tx) == 0 {
		return ""
	}
	// Start with leaf hashes
	var layer []string
	for _, t := range tx {
		h := sha256.Sum256([]byte(t))
		layer = append(layer, hex.EncodeToString(h[:]))
	}
	// Build up until single root
	for len(layer) > 1 {
		var next []string
		for i := 0; i < len(layer); i += 2 {
			if i+1 < len(layer) {
				concat := layer[i] + layer[i+1]
				h := sha256.Sum256([]byte(concat))
				next = append(next, hex.EncodeToString(h[:]))
			} else {
				// If odd, duplicate the last hash (common Merkle convention)
				concat := layer[i] + layer[i]
				h := sha256.Sum256([]byte(concat))
				next = append(next, hex.EncodeToString(h[:]))
			}
		}
		layer = next
	}
	return layer[0]
}
