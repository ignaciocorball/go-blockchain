// Package blockchain implements the transaction system for the UFChain blockchain.
// This file handles the creation, signing, and serialization of transactions
// using ECDSA (Elliptic Curve Digital Signature Algorithm) for security.
package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
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
//   - TransactionID: The ID of the transaction containing the UTXO being spent
//   - OutputIndex: The index of the UTXO in the original transaction
//   - Signature: The cryptographic signature proving ownership
//   - PublicKey: The public key of the sender
type TxInput struct {
	TransactionID []byte // ID of the transaction containing the UTXO being spent
	OutputIndex   int    // Index of the UTXO in the original transaction
	Signature     []byte // ECDSA signature of the transaction
	PublicKey     []byte // Sender's public key
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
//   - fromWallet: The sender's wallet
//   - toPublicKey: The recipient's public key (as a string)
//   - amount: The amount to transfer
//   - utxos: The list of UTXOs available for this transaction
//
// The function:
// 1. Verifies that there are enough UTXOs to cover the amount
// 2. Creates a new transaction with the appropriate inputs and outputs
// 3. Signs the transaction using the sender's private key
// 4. Returns the complete transaction
//
// Returns nil and an error if the transaction cannot be created.
func NewTransaction(fromWallet *Wallet, toPublicKey string, amount int, utxos []*UTXO) (*Transaction, error) {
	var inputs []TxInput
	var outputs []TxOutput
	var totalInput int

	// Verify that there are enough UTXOs to cover the amount
	for _, utxo := range utxos {
		if bytes.Equal(utxo.PublicKey, fromWallet.PublicKey) {
			totalInput += utxo.Value
			input := TxInput{
				TransactionID: utxo.TransactionID,
				OutputIndex:   utxo.OutputIndex,
				PublicKey:     fromWallet.PublicKey,
			}
			inputs = append(inputs, input)

			if totalInput >= amount {
				break
			}
		}
	}

	if totalInput < amount {
		return nil, fmt.Errorf("insufficient funds: have %d, need %d", totalInput, amount)
	}

	// Create the output for the recipient using their public key directly
	outputs = append(outputs, TxOutput{
		Value:     amount,
		PublicKey: []byte(toPublicKey), // The public key is already in the correct format
	})

	// Create change output if necessary
	if totalInput > amount {
		outputs = append(outputs, TxOutput{
			Value:     totalInput - amount,
			PublicKey: fromWallet.PublicKey,
		})
	}

	tx := &Transaction{
		Input:  inputs,
		Output: outputs,
	}

	// Sign the transaction
	tx.ID = tx.HashTransaction()
	for i := range tx.Input {
		tx.Input[i].Signature = tx.Sign(fromWallet.GetPrivateKey())
	}

	return tx, nil
}

// Sign signs the transaction with the private key
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) []byte {
	// Create a copy of the transaction without signatures
	txCopy := tx.TrimmedCopy()

	// Sign the transaction hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, txCopy.ID)
	if err != nil {
		log.Panic(err)
	}

	// Concatenate r and s in DER format
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

// Verify verifies the transaction signature
func (tx *Transaction) Verify() bool {
	// Create a copy of the transaction without signatures
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for _, input := range tx.Input {
		// The public key is already in the correct format (X || Y)
		if len(input.PublicKey) != 64 { // 32 bytes for X + 32 bytes for Y
			return false
		}

		// Extract X and Y from the public key
		x := new(big.Int).SetBytes(input.PublicKey[:32])
		y := new(big.Int).SetBytes(input.PublicKey[32:])

		// Create the public key
		publicKey := &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		}

		// Split the signature into r and s
		if len(input.Signature) != 64 { // 32 bytes for r + 32 bytes for s
			return false
		}
		r := new(big.Int).SetBytes(input.Signature[:32])
		s := new(big.Int).SetBytes(input.Signature[32:])

		// Verify the signature
		if !ecdsa.Verify(publicKey, txCopy.ID, r, s) {
			return false
		}
	}
	return true
}

// TrimmedCopy creates a trimmed copy of the transaction without signatures
func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, input := range tx.Input {
		inputs = append(inputs, TxInput{
			TransactionID: input.TransactionID,
			OutputIndex:   input.OutputIndex,
			PublicKey:     input.PublicKey,
		})
	}

	for _, output := range tx.Output {
		outputs = append(outputs, TxOutput{
			Value:     output.Value,
			PublicKey: output.PublicKey,
		})
	}

	txCopy := &Transaction{
		Input:  inputs,
		Output: outputs,
	}
	txCopy.ID = txCopy.HashTransaction()
	return txCopy
}

// HashTransaction generates a unique hash for the transaction.
// The hash is created by combining:
//   - Transaction ID
//   - Output index
//   - Sender's public key
//   - Recipient's public key
//   - Transaction amount
//
// Returns a SHA-256 hash of the combined data.
func (tx *Transaction) HashTransaction() []byte {
	var hash [32]byte

	// Use transaction data directly without creating a copy
	var data [][]byte
	for _, input := range tx.Input {
		data = append(data, input.TransactionID)
		data = append(data, []byte(fmt.Sprintf("%d", input.OutputIndex)))
		data = append(data, input.PublicKey)
	}
	for _, output := range tx.Output {
		data = append(data, output.PublicKey)
		data = append(data, []byte(fmt.Sprintf("%d", output.Value)))
	}

	hash = sha256.Sum256(bytes.Join(data, []byte{}))
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
