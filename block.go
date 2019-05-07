package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/payload"
)

// Block represents the base structure for recording credit data
type Block struct {
	Timestamp    string
	Amount       float64
	To           string
	From         string
	Hash         string
	PreviousHash string
}

func calculateHash(block Block) string {

	amount := strconv.FormatFloat(block.Amount, 'g', 1, 64)

	record := amount + block.From +
		block.PreviousHash +
		block.PreviousHash +
		block.Timestamp +
		block.To

	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// GenerateBlock creates a new Block
func GenerateBlock(oldBlock *Block, amount float64, to string, from string) (Block, error) {

	var newBlock Block

	newBlock.Timestamp = time.Now().String()
	newBlock.Amount = amount
	newBlock.PreviousHash = oldBlock.PreviousHash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

// IsBlockValid compares two block to determine if they are the same
func IsBlockValid(newBlock *Block, oldBlock *Block) bool {

	if oldBlock.Hash != newBlock.Hash {
		return false
	}

	if calculateHash(*newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// ReplaceChain replaces the current chain if the blocks are longer
func ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > 2 {
	}
}

func (block Block) Read(reader payload.Reader) (noise.Message, error) {
	bytes, err := reader.ReadBytes()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &block)
	if err != nil {
		panic(err)
	}
	fmt.Println(block.Hash)

	return block, err
}

func (block Block) Write() []byte {
	bytes, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	return payload.NewWriter(nil).WriteBytes(bytes).Bytes()
}
