package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Header struct {
	Nonce int    `json:"nonce"`
	Hash  string `json:"hash"`
}

type Payload struct {
	Sequence     int    `json:"sequence"`
	Timestamp    int64  `json:"timestamp"`
	Data         string `json:"data"`
	PreviousHash string `json:"previousHash"`
}

type Block struct {
	Header  `json:"header"`
	Payload `json:"payload"`
}

type BlockChain struct {
	Chain      []Block `json:"chain"`
	prefix     string
	difficulty int
}

func NewBlockChain(prefix byte, difficulty int) BlockChain {
	blockChain := BlockChain{
		[]Block{},
		string(prefix),
		difficulty,
	}
	genesis := blockChain.createGenesisBlock()
	blockChain.Chain = append(blockChain.Chain, genesis)
	return blockChain
}

func (c BlockChain) createGenesisBlock() Block {
	payload := Payload{
		Sequence:     0,
		Timestamp:    time.Now().UnixMilli(),
		Data:         "Genesis Block",
		PreviousHash: "",
	}
	header := Header{
		Nonce: 0,
		Hash:  c.HashPayload(payload),
	}
	return Block{
		header,
		payload,
	}
}

func (c BlockChain) GetLastBlock() Block {
	return c.Chain[len(c.Chain)-1]
}

func (c BlockChain) CreateBlockPayload(data string) Payload {
	lastBlock := c.GetLastBlock()
	payload := Payload{
		Sequence:     lastBlock.Payload.Sequence + 1,
		Timestamp:    time.Now().UnixMilli(),
		Data:         data,
		PreviousHash: lastBlock.Header.Hash,
	}
	return payload
}

func (c BlockChain) MineBlock(payload Payload) Block {
	nonce := 0
	startTime := time.Now().UnixMilli()

	for {
		payloadHash := c.HashPayload(payload)
		proofingHash := c.Hash(fmt.Sprintf("%s%d", payloadHash, nonce))
		if c.IsHashProofed(proofingHash) {
			endTime := time.Now().UnixMilli()
			shortHash := payloadHash[0:12]
			mineTime := endTime - startTime

			fmt.Printf("Mined block %d in %d milliseconds. Hash: %s (%d attempts)\n", payload.Sequence, mineTime, shortHash, nonce)

			header := Header{
				Nonce: nonce,
				Hash:  payloadHash,
			}
			return Block{
				header,
				payload,
			}
		}
		nonce++
	}
}

func (c *BlockChain) PushBlock(block Block) {
	if c.verifyBlock(block) {
		c.Chain = append(c.Chain, block)
		fmt.Printf("Pushed block %d sequence %s hash\n", block.Payload.Sequence, block.Header.Hash)
	}
}

func (c BlockChain) verifyBlock(block Block) bool {
	lastBlock := c.GetLastBlock()
	if block.Payload.PreviousHash != lastBlock.Header.Hash {
		return false
	}
	payloadHash := c.HashPayload(block.Payload)
	check := c.Hash(fmt.Sprintf("%s%d", payloadHash, block.Header.Nonce))
	if !c.IsHashProofed(check) {
		return false
	}
	return true
}

func (c BlockChain) HashPayload(payload Payload) string {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return c.Hash(string(payloadJson))
}

func (c BlockChain) Hash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func (c BlockChain) IsHashProofed(hash string) bool {
	check := strings.Repeat(c.prefix, c.difficulty)
	return strings.HasPrefix(hash, check)
}
