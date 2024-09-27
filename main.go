package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var VERBOSE bool = false

func main() {
	inputFile := flag.String("input-file", "", "The file to read its content and add to block chain")
	chainFile := flag.String("chain-file", "", "The chain input/output file to read and write")
	help := flag.Bool("help", false, "Show help menu")
	verbose := flag.Bool("verbose", false, "Show information about mining")

	flag.Parse()

	VERBOSE = *verbose

	if *help {
		showHelp()
		return
	}

	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		panic(fmt.Sprintf("%s file does not exist", *inputFile))
	}
	if _, err := os.Stat(*chainFile); os.IsNotExist(err) {
		createBlockChain(*inputFile, *chainFile)
	} else {
		insertIntoBlockChain(*inputFile, *chainFile)
	}
}

func createBlockChain(inputFile string, chainFile string) {
	blockChain := NewBlockChain('0', 3)
	content, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	payload := blockChain.CreateBlockPayload(string(content))
	block := blockChain.MineBlock(payload)
	blockChain.PushBlock(block)

	chain, _ := json.MarshalIndent(blockChain, "", "  ")
	file, _ := os.Create(chainFile)
	defer file.Close()

	file.Write(chain)
}

func insertIntoBlockChain(inputFile string, chainFile string) {
	blockChain := NewBlockChain('0', 3)
	chainContent, _ := os.ReadFile(chainFile)
	json.Unmarshal(chainContent, &blockChain)

	content, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	payload := blockChain.CreateBlockPayload(string(content))
	block := blockChain.MineBlock(payload)
	blockChain.PushBlock(block)

	chain, _ := json.MarshalIndent(blockChain, "", "  ")
	file, _ := os.Create(chainFile)
	defer file.Close()

	file.Write(chain)
}

func showHelp() {
	fmt.Println("Welcome to Avraham's custom Blockchain")
	fmt.Println("Flags: ")
	fmt.Println("--input-file specify which file to read it's content and add to block chain")
	fmt.Println("--chain-file specify the json file which it's located the block chain, if it doesn't exist, a new will be created")
	fmt.Println("--help show this menu")
	fmt.Println("--verbose show details of block chain and mining")
}
