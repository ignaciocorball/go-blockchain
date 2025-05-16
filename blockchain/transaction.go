// Package blockchain implements the transaction system for the UFChain blockchain.
// This file handles the creation, signing, and serialization of transactions
// using ECDSA (Elliptic Curve Digital Signature Algorithm) for security.
package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// Transaction represents a transfer of value in the blockchain.
// Each transaction consists of:
//   - ID: A unique identifier (hash) of the transaction
//   - Input: The source of the transaction (previous unspent output)
//   - Output: The destination and amount of the transfer
type Transaction struct {
	ID     []byte     // Transaction hash
	Input  []TxInput  // Transaction inputs (sources)
	Output []TxOutput // Transaction outputs (destinations)
}

// TxInput represents the source of a transaction.
// It contains:
//   - Signature: The cryptographic signature proving ownership
//   - PublicKey: The public key of the sender
type TxInput struct {
	Signature []byte // ECDSA signature of the transaction
	PublicKey []byte // Sender's public key
}

// TxOutput represents the destination of a transaction.
// It contains:
//   - Value: The amount being transferred
//   - PublicKey: The public key of the recipient
type TxOutput struct {
	Value     int    // Amount to transfer
	PublicKey []byte // Recipient's public key
}

// NewTransaction creates a new transaction in the blockchain.
// Parameters:
//   - privateKey: The sender's private key for signing
//   - recipient: The recipient's public key
//   - amount: The amount to transfer
//
// The function:
// 1. Creates a new transaction with input and output
// 2. Generates a hash of the transaction
// 3. Signs the transaction using the sender's private key
// 4. Returns the complete transaction
//
// Panics if signing fails.
func NewTransaction(privateKey ecdsa.PrivateKey, recipient []byte, amount int) *Transaction {
	txIn := TxInput{}
	txOut := TxOutput{Value: amount, PublicKey: recipient}

	tx := Transaction{
		Input:  []TxInput{txIn},
		Output: []TxOutput{txOut},
	}

	tx.ID = tx.hashTransaction()
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, tx.ID)

	if err != nil {
		log.Panic(err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	txIn.Signature = signature

	return &tx
}

// hashTransaction generates a unique hash for the transaction.
// The hash is created by combining:
//   - Sender's public key
//   - Recipient's public key
//   - Transaction amount
//
// Returns a SHA-256 hash of the combined data.
func (tx *Transaction) hashTransaction() []byte {
	var hash [32]byte

	hash = sha256.Sum256(bytes.Join([][]byte{
		tx.Input[0].PublicKey,
		tx.Output[0].PublicKey,
		[]byte(string(tx.Output[0].Value)),
	}, []byte{}))

	return hash[:]
}

// Serialize converts the transaction into a byte array for storage or transmission.
// Uses gob encoding to serialize the entire transaction structure.
// Returns the serialized transaction as a byte slice.
// Panics if serialization fails.
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)

	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

// DeserializeTransaction reconstructs a Transaction from its serialized byte array.
// Takes a byte slice containing the serialized transaction data.
// Returns a pointer to the reconstructed Transaction.
// Panics if deserialization fails.
func DeserializeTransaction(data []byte) *Transaction {
	var transaction Transaction

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return &transaction
}
