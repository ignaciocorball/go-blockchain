package blockchain

import (
	"bytes"
	"fmt"
)

// UTXO represents an unspent transaction output
type UTXO struct {
	TransactionID []byte // ID of the transaction that created this UTXO
	OutputIndex   int    // Index of the output in the transaction
	Value         int    // Amount of tokens
	PublicKey     []byte // Public key of the owner
}

// UTXOSet manages the set of unspent UTXOs
type UTXOSet struct {
	UTXOs map[string]*UTXO // Map of unspent UTXOs, key = "txID_outputIndex"
}

// NewUTXOSet creates a new UTXO set
func NewUTXOSet() *UTXOSet {
	return &UTXOSet{
		UTXOs: make(map[string]*UTXO),
	}
}

// AddUTXO adds a new UTXO to the set
func (us *UTXOSet) AddUTXO(txID []byte, outputIndex int, value int, publicKey []byte) {
	key := fmt.Sprintf("%x_%d", txID, outputIndex)
	us.UTXOs[key] = &UTXO{
		TransactionID: txID,
		OutputIndex:   outputIndex,
		Value:         value,
		PublicKey:     publicKey,
	}
}

// RemoveUTXO removes a UTXO from the set
func (us *UTXOSet) RemoveUTXO(txID []byte, outputIndex int) {
	key := fmt.Sprintf("%x_%d", txID, outputIndex)
	delete(us.UTXOs, key)
}

// GetUTXOsForAddress returns all UTXOs for a specific address
func (us *UTXOSet) GetUTXOsForAddress(address []byte) []*UTXO {
	var utxos []*UTXO
	for _, utxo := range us.UTXOs {
		if bytes.Equal(utxo.PublicKey, address) {
			utxos = append(utxos, utxo)
		}
	}
	return utxos
}

// GetBalance calculates the total balance for an address
func (us *UTXOSet) GetBalance(address []byte) int {
	var balance int
	for _, utxo := range us.GetUTXOsForAddress(address) {
		balance += utxo.Value
	}
	return balance
}
