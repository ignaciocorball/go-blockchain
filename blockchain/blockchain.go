// Package blockchain implements the core blockchain data structure and operations.
// It provides functionality for creating and managing the blockchain, including
// block creation, validation, and chain management.
package blockchain

import (
	"bytes"
	"fmt"
)

// Blockchain represents the main blockchain structure.
// It maintains an ordered list of blocks, where each block is linked to its
// previous block through cryptographic hashes, forming an immutable chain.
type Blockchain struct {
	Blocks []*Block // Ordered list of blocks in the chain
	UTXOs  *UTXOSet
}

// GetBlock retrieves a block from the blockchain by its hash.
// Parameters:
//   - hash: The cryptographic hash of the block to find
//
// Returns:
//   - A slice containing the found block if successful
//   - An error if the block is not found
//
// Note: Returns a slice of blocks to maintain consistency with potential
// future implementations that might support multiple blocks with the same hash
// (though this is not currently supported).
func (bc *Blockchain) GetBlock(hash []byte) ([]*Block, error) {
	for _, block := range bc.Blocks {
		if bytes.Equal(block.Hash, hash) {
			return []*Block{block}, nil
		}
	}
	return nil, fmt.Errorf("block not found")
}

// NewBlockchain creates a new blockchain instance with a genesis block.
// Parameters:
//   - genesisBlock: The first block in the chain that initializes the blockchain
//
// Returns a new blockchain instance containing only the genesis block.
// The genesis block is special as it has no previous block and typically
// contains initial system state or configuration.
func NewBlockchain(genesisBlock *Block) *Blockchain {
	bc := &Blockchain{
		Blocks: []*Block{genesisBlock},
		UTXOs:  NewUTXOSet(),
	}

	// Process the genesis block
	bc.UpdateUTXOs(genesisBlock)

	return bc
}

// UpdateUTXOs updates the UTXO set based on a new block
func (bc *Blockchain) UpdateUTXOs(block *Block) {
	for _, tx := range block.Transactions {
		// Remove spent UTXOs
		for _, input := range tx.Input {
			bc.UTXOs.RemoveUTXO(input.TransactionID, input.OutputIndex)
		}

		// Add new UTXOs
		for i, output := range tx.Output {
			bc.UTXOs.AddUTXO(tx.ID, i, output.Value, output.PublicKey)
		}
	}
}

// AddBlock creates and adds a new block to the blockchain.
// Parameters:
//   - transactions: List of transactions to be included in the new block
//   - validator: Public key of the validator who created this block
//
// The function:
// 1. Gets the previous block (last block in the chain)
// 2. Creates a new block with the provided transactions
// 3. Links it to the previous block using the previous block's hash
// 4. Adds the new block to the chain
//
// Returns the newly created block.
func (bc *Blockchain) AddBlock(transactions []*Transaction, validator []byte) *Block {
	// Verify all transactions
	for _, tx := range transactions {
		if !tx.Verify() {
			panic("Invalid transaction signature")
		}
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(transactions, prevBlock.Hash, validator)

	// Update UTXOs before adding the block
	bc.UpdateUTXOs(newBlock)

	bc.Blocks = append(bc.Blocks, newBlock)
	return newBlock
}

// GetBalance returns the balance of an address
func (bc *Blockchain) GetBalance(address []byte) int {
	return bc.UTXOs.GetBalance(address)
}
