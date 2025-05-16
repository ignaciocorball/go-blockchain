// Package storage implements the persistent storage layer for the UFChain blockchain.
// This file uses BadgerDB, a key-value store optimized for SSDs, to store blockchain data.
// BadgerDB provides ACID transactions and high performance for blockchain operations.
package storage

import (
	"log"

	"github.com/davepartner/go-blockchain/blockchain"
	"github.com/dgraph-io/badger"
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
//
// Returns a new BlockchainDB instance.
// Panics if database initialization fails.
func OpenDB(path string) *BlockchainDB {
	opts := badger.DefaultOptions(path)
	opts.Truncate = true // Allow value log truncation
	opts.Logger = nil    // Disable logging to reduce noise
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
//
// Returns:
//   - nil if storage is successful
//   - error if storage fails
//
// Note: Uses Badger's Update transaction for atomic writes
func (bdb *BlockchainDB) SaveBlock(block *blockchain.Block) error {
	return bdb.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.Hash, block.Serialize())
		if err != nil {
			return err
		}
		return nil
	})
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

// CloseDB safely closes the database connection.
// This function should be called when the application is shutting down
// to ensure proper cleanup of resources.
//
// Panics if database closure fails, as this indicates a serious
// system-level issue that should be addressed immediately.
func (bdb *BlockchainDB) CloseDB() {
	err := bdb.DB.Close()
	if err != nil {
		log.Panic(err)
	}
}
