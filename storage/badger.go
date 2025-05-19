// Package storage implements the persistent storage layer for the UFChain blockchain.
// This file uses BadgerDB, a key-value store optimized for SSDs, to store blockchain data.
// BadgerDB provides ACID transactions and high performance for blockchain operations.
package storage

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/ignaciocorball/go-blockchain/blockchain"
)

// BlockchainDB wraps the Badger database instance and provides
// blockchain-specific storage operations. It handles:
//   - Block storage and retrieval
//   - Transaction management
//   - Database lifecycle
type BlockchainDB struct {
	DB *badger.DB // Badger database instance
}

// OpenDB initializes and opens a new Badger database instance.
// Parameters:
//   - path: File system path where the database will be stored
//
// Configuration:
//   - Truncate: Enabled to allow value log truncation
//   - Logger: Disabled to reduce noise
//   - SyncWrites: Enabled to ensure data durability
//   - NumVersionsToKeep: Set to 1 to avoid version conflicts
//
// Returns a new BlockchainDB instance.
// Panics if database initialization fails.
func OpenDB(path string) *BlockchainDB {
	opts := badger.DefaultOptions(path)
	opts.Truncate = true
	opts.Logger = nil
	opts.SyncWrites = true
	opts.NumVersionsToKeep = 1
	opts.ValueLogFileSize = 1024 * 1024 * 10 // 10MB
	opts.MaxTableSize = 64 << 20             // 64MB
	opts.ValueLogMaxEntries = 1000000        // 1M entries

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}
	return &BlockchainDB{DB: db}
}

// SaveBlock stores a block in the database.
// Parameters:
//   - block: The block to be stored
//
// The function:
// 1. Starts a new transaction
// 2. Serializes the block
// 3. Stores it using the block's hash as the key
// 4. Commits the transaction
//
// Returns:
//   - nil if storage is successful
//   - error if storage fails
func (bdb *BlockchainDB) SaveBlock(block *blockchain.Block) error {
	txn := bdb.DB.NewTransaction(true)
	defer txn.Discard()

	// Serialize the block
	blockData := block.Serialize()

	// Save the block using its hash as the key
	err := txn.Set(block.Hash, blockData)
	if err != nil {
		return fmt.Errorf("error saving block: %v", err)
	}

	// Commit the transaction
	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("error committing block: %v", err)
	}

	return nil
}

// GetBlock retrieves a block from the database by its hash.
// Parameters:
//   - hash: The hash of the block to retrieve
//
// The function:
// 1. Starts a read-only transaction
// 2. Retrieves the block data
// 3. Deserializes the block
//
// Returns:
//   - The retrieved block if found
//   - nil and error if block doesn't exist or retrieval fails
//
// Note: Uses Badger's View transaction for read-only operations
func (bdb *BlockchainDB) GetBlock(hash []byte) (*blockchain.Block, error) {
	var block *blockchain.Block

	err := bdb.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			block = blockchain.DeserializeBlock(val)
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	return block, nil
}

// SaveWallet stores a wallet in the database
func (bdb *BlockchainDB) SaveWallet(address string, wallet *blockchain.Wallet) error {
	txn := bdb.DB.NewTransaction(true)
	defer txn.Discard()

	// Serialize the wallet
	walletData := wallet.Serialize()

	// Save the wallet using its address as the key
	key := []byte("wallet_" + address)
	err := txn.Set(key, walletData)
	if err != nil {
		return fmt.Errorf("error saving wallet: %v", err)
	}

	// Commit the transaction
	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("error committing wallet: %v", err)
	}

	return nil
}

// GetWallet retrieves a wallet from the database
func (bdb *BlockchainDB) GetWallet(address string) (*blockchain.Wallet, error) {
	txn := bdb.DB.NewTransaction(false)
	defer txn.Discard()

	key := []byte("wallet_" + address)
	item, err := txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("wallet not found: %s", address)
		}
		return nil, fmt.Errorf("error getting wallet: %v", err)
	}

	var walletData []byte
	err = item.Value(func(val []byte) error {
		walletData = append([]byte{}, val...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error reading wallet data: %v", err)
	}

	return blockchain.DeserializeWallet(walletData), nil
}

// CloseDB safely closes the database connection.
// This function should be called when the application is shutting down
// to ensure proper cleanup of resources.
func (bdb *BlockchainDB) CloseDB() {
	if bdb.DB != nil {
		err := bdb.DB.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
}
