package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

// Message represents the new data being added to the new block in a blockchain
type Message struct {
	BPM int
}

func genesis() {
	t := time.Now()
	genesisBlock := Block{Index: 0, Timestamp: t.String(), BPM: 0, Hash: "", PrevHash: ""}
	Blockchain = append(Blockchain, genesisBlock)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	genesis()
	log.Fatal(run())
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) Block {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
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

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err == nil {
		fmt.Println("DANGER INSIDE DECODE REWEROROR")
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	oldBlock := Blockchain[len(Blockchain)-1]
	newBlock := generateBlock(oldBlock, m.BPM)

	if isBlockValid(newBlock, oldBlock) {
		Blockchain = append(Blockchain, newBlock)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		http.Error(w, "HTTP 500: Internal Server Error", http.StatusInternalServerError)
		// w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
