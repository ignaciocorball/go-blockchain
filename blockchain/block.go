// Package blockchain implements the core blockchain data structures and operations
// for the UFChain blockchain system. It provides functionality for creating and
// managing blocks, transactions, and the overall blockchain structure.
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// Block represents a single block in the blockchain. Each block contains:
// - Timestamp: When the block was created
// - Transactions: List of transactions included in this block
// - Hash: The cryptographic hash of this block
// - PrevHash: The hash of the previous block in the chain
// - Validator: The public key of the validator who created this block
// - Nonce: A number used in the proof-of-work/proof-of-stake mechanism
type Block struct {
	Timestamp    string
	Transactions []*Transaction
	Hash         []byte
	PrevHash     []byte
	Validator    []byte
	Nonce        int
}

// NewBlock creates and returns a new block in the blockchain.
// Parameters:
//   - transactions: List of transactions to be included in the block
//   - prevHash: Hash of the previous block in the chain
//   - validator: Public key of the validator creating this block
//
// The function initializes a new block with the current timestamp,
// calculates its hash, and returns the complete block structure.
func NewBlock(transactions []*Transaction, prevHash []byte, validator []byte) *Block {
	block := &Block{
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Validator:    validator,
		Nonce:        0,
	}

	block.Hash = block.calculateHash()
	return block
}

// calculateHash generates the cryptographic hash of the block.
// The hash is calculated by combining:
// - The previous block's hash
// - All transaction IDs in the block
// - The block's timestamp
//
// Returns a SHA-256 hash of the combined data as a byte slice.
func (b *Block) calculateHash() []byte {
	var txHashes []byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID...)
	}

	hash := sha256.Sum256(bytes.Join([][]byte{
		b.PrevHash,
		txHashes,
		[]byte(b.Timestamp),
	}, []byte{}))

	return hash[:]
}

// Serialize converts the block into a byte array for storage or transmission.
// Uses gob encoding to serialize the entire block structure.
// Returns the serialized block as a byte slice.
// Panics if serialization fails.
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock reconstructs a Block from its serialized byte array.
// Takes a byte slice containing the serialized block data.
// Returns a pointer to the reconstructed Block.
// Panics if deserialization fails.
func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}

	return &block
}
