package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	blockChain := NewBlockChain('0', 2)
	for i := 1; i < 5; i++ {
		payload := blockChain.CreateBlockPayload(fmt.Sprintf("Block %d", i))
		block := blockChain.MineBlock(payload)
		blockChain.PushBlock(block)
	}

	chain, _ := json.MarshalIndent(blockChain, "", "  ")
	file, _ := os.Create("chain.json")
	defer file.Close()

	file.Write(chain)
}
