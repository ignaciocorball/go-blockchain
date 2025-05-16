// Package blockchain implements the consensus mechanism for the UFChain blockchain.
// This file specifically handles the Proof of Stake (PoS) consensus algorithm,
// which is used to select validators for block creation based on their stake in the system.
package blockchain

import (
	"crypto/rand"
	"log"
	"math/big"
)

// PosValidator represents a validator in the Proof of Stake system.
// Each validator has:
//   - PublicKey: The cryptographic public key used to identify the validator
//   - Stake: The amount of tokens the validator has staked in the system
//
// The stake amount determines the validator's probability of being selected
// to create the next block. Higher stake means higher probability of selection.
type PosValidator struct {
	PublicKey []byte // Validator's public key for identification
	Stake     int    // Amount of tokens staked by the validator
}

// ProofOfStake implements the Proof of Stake consensus algorithm.
// It selects a validator to create the next block based on their stake in the system.
//
// Parameters:
//   - validators: A map of validator public keys to their validator information
//
// The selection process:
//  1. Calculates the total stake across all validators
//  2. Generates a cryptographically secure random number between 0 and total stake
//  3. Selects a validator based on their proportional stake
//     (validators with higher stakes have higher probability of selection)
//
// Returns:
//   - The public key of the selected validator as a string
//
// Panics:
//   - If there's an error generating the random number
//   - If no validator can be selected (should never happen in normal operation)
//
// Example:
//
//	If there are two validators with stakes 70 and 30:
//	- First validator has 70% chance of being selected
//	- Second validator has 30% chance of being selected
func ProofOfStake(validators map[string]*PosValidator) string {
	// Calculate total stake across all validators
	totalStake := 0
	for _, validator := range validators {
		totalStake += validator.Stake
	}

	// Generate a cryptographically secure random number between 0 and total stake
	randomBig, err := rand.Int(rand.Reader, big.NewInt(int64(totalStake)))
	if err != nil {
		log.Panic(err)
	}

	random := randomBig.Int64()

	// Select a validator based on their stake proportion
	// Validators with higher stakes have a higher probability of being selected
	for _, validator := range validators {
		random -= int64(validator.Stake)
		if random <= 0 {
			return string(validator.PublicKey)
		}
	}

	// This should never happen in normal operation
	// as long as there are validators with positive stakes
	log.Panic("Unable to find a validator")
	return ""
}
