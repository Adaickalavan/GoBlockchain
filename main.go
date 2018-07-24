package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	Index     int    // Position of the block in blockchain
	Timestamp string // Time data is written
	BPM       int    // Data to be stored, Beats per minute
	Hash      string // SHA256 identifier of the block
	PrevHash  string // SHA256 identifier of the previous block
}

// Blockchain represents the sequence of blocks forming a blockchain
var Blockchain []Block

func main() {
	fmt.Println("hi")
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	switch {
	case oldBlock.Index+1 != newBlock.Index:
		return false
	case oldBlock.Hash != newBlock.PrevHash:
		return false
	case calculateHash(newBlock) != newBlock.Hash:
		return false
	}
	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}