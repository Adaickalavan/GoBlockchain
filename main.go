package main

import (
	"fmt"
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
