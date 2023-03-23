package main

import (
	"fmt"
	"log"

	"go-block/chain"
	"go-block/wallet"
)

func init() {
	log.SetPrefix("Block: ")
}

func main() {
	minerWallet := wallet.NewWallet()
	AWallet := wallet.NewWallet()
	BWallet := wallet.NewWallet()

	// Wallet
	transaction := wallet.NewTransaction(AWallet.PrivateKey(), AWallet.PublicKey(), AWallet.BlockchainAddress(), BWallet.BlockchainAddress(), 1.0)
	fmt.Printf("Signature %s \n", transaction.GenerateSignature())

	// Blockchain
	blockchain := chain.NewBlockchain(minerWallet.BlockchainAddress(), 5000)
	isAdded := blockchain.AddTransaction(AWallet.BlockchainAddress(), BWallet.BlockchainAddress(), 1.0,
		AWallet.PublicKey(), transaction.GenerateSignature())
	fmt.Println("Added? ", isAdded)

	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("A %.1f \n", blockchain.CalculateTotalAmount(AWallet.BlockchainAddress()))
	fmt.Printf("B %.1f \n", blockchain.CalculateTotalAmount(BWallet.BlockchainAddress()))
	fmt.Printf("Miner %.1f \n", blockchain.CalculateTotalAmount(minerWallet.BlockchainAddress()))
}
