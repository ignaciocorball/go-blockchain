// Package main is the entry point for the UFChain blockchain application.
// It initializes the blockchain, database, and API server components,
// and orchestrates their interaction to run the blockchain node.
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ignaciocorball/go-blockchain/api"
	"github.com/ignaciocorball/go-blockchain/blockchain"
	"github.com/ignaciocorball/go-blockchain/storage"
)

// main initializes and starts the UFChain blockchain node.
// The function performs the following steps in order:
// 1. Creates a genesis block to initialize the blockchain
// 2. Initializes the blockchain with the genesis block
// 3. Sets up the Badger database for persistent storage
// 4. Persists the genesis block to the database
// 5. Starts the API server to handle external requests
//
// The genesis block is special as it:
//   - Has no transactions
//   - Has no previous block hash
//   - Is created by a special genesis validator
//
// The database is configured to store blocks in "./storage/badger"
// and is properly closed when the application exits.
//
// The API server runs on the default port (1323) and provides
// endpoints for blockchain operations.
func main() {
	// Create the genesis block with:
	// - Empty transaction list
	// - Empty previous hash
	// - Special genesis validator
	genesisBlock := blockchain.NewBlock([]*blockchain.Transaction{}, []byte{}, []byte("genesis-validator"))

	// Initialize the blockchain with the genesis block
	bc := blockchain.NewBlockchain(genesisBlock)

	// Initialize the Badger database for persistent storage
	// The database will be stored in the ./storage/badger directory
	db := storage.OpenDB("./storage/badger")

	// Configurar el manejo de señales para un cierre limpio
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Goroutine para manejar el cierre limpio
	go func() {
		<-sigChan
		fmt.Println("\nCerrando la aplicación...")
		db.CloseDB()
		os.Exit(0)
	}()

	// Persist the genesis block to the database
	// This ensures the blockchain can be recovered if the application restarts
	err := db.SaveBlock(genesisBlock)
	if err != nil {
		log.Printf("Error saving genesis block: %v", err)
		db.CloseDB()
		os.Exit(1)
	}

	// Start the API server with the blockchain and database instances
	// This will begin listening for incoming requests
	fmt.Println("Iniciando servidor en http://localhost:1323")
	api.StartServer(bc, db)
}
