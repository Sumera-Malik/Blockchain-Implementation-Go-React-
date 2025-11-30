package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
	"strconv"
)

type Block struct {
	Index      int      `json:"index"`
	Timestamp  string   `json:"timestamp"`
	Data       []string `json:"data"`       // transactions in this block
	PrevHash   string   `json:"prevHash"`
	Hash       string   `json:"hash"`
	Nonce      int      `json:"nonce"`
	MerkleRoot string   `json:"merkleRoot"`
}

type Blockchain struct {
	Name        string
	Blocks      []Block
	PendingTx   []string
	Difficulty  int // number of leading zeros required
}

func calculateHash(b Block) string {
	record := strconv.Itoa(b.Index) + b.Timestamp + strings.Join(b.Data, "|") +
		b.PrevHash + strconv.Itoa(b.Nonce) + b.MerkleRoot
	h := sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}

func hasLeadingZeros(hexHash string, zeros int) bool {
	prefix := strings.Repeat("0", zeros)
	return strings.HasPrefix(hexHash, prefix)
}

func (bc *Blockchain) createGenesis() Block {
	genesis := Block{
		Index:     0,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      []string{"Genesis: " + bc.Name + " Blockchain"},
		PrevHash:  "",
		Nonce:     0,
	}
	genesis.MerkleRoot = merkleRoot(genesis.Data)
	genesis.Hash = calculateHash(genesis)
	return genesis
}

func (bc *Blockchain) lastBlock() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) AddTransaction(tx string) {
	if strings.TrimSpace(tx) != "" {
		bc.PendingTx = append(bc.PendingTx, tx)
	}
}

func (bc *Blockchain) MinePending() (*Block, error) {
	if len(bc.PendingTx) == 0 {
		return nil, ErrNoPendingTx
	}
	prev := bc.lastBlock()
	newBlock := Block{
		Index:     prev.Index + 1,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      append([]string{}, bc.PendingTx...), // copy
		PrevHash:  prev.Hash,
		Nonce:     0,
	}
	newBlock.MerkleRoot = merkleRoot(newBlock.Data)

	// Proof of Work
	for {
		newBlock.Hash = calculateHash(newBlock)
		if hasLeadingZeros(newBlock.Hash, bc.Difficulty) {
			break
		}
		newBlock.Nonce++
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	bc.PendingTx = []string{} // clear pool
	return &newBlock, nil
}

var ErrNoPendingTx = &bcError{"no pending transactions to mine"}

type bcError struct{ msg string }
func (e *bcError) Error() string { return e.msg }
