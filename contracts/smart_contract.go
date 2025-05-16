// Package contracts implements the smart contract system for the UFChain blockchain.
// This file handles the creation, validation, and execution of smart contracts,
// providing a basic framework for programmable transactions on the blockchain.
package contracts

import (
	"errors"
	"fmt"
	"time"
)

// SmartContract represents a programmable contract on the blockchain.
// Each contract consists of:
//   - ID: A unique identifier for the contract
//   - Code: The contract's executable code
//   - State: A key-value store for the contract's persistent state
//   - CreatedAt: Timestamp of contract creation
//
// The contract's state is mutable and persists between executions,
// allowing for stateful contract behavior.
type SmartContract struct {
	ID        string                 // Unique identifier for the contract
	Code      string                 // Contract's executable code
	State     map[string]interface{} // Contract's persistent state storage
	CreatedAt time.Time              // Contract creation timestamp
}

// NewSmartContract creates a new smart contract instance.
// Parameters:
//   - id: Unique identifier for the contract
//   - code: The contract's executable code
//
// Returns a new smart contract with:
//   - Initialized state map
//   - Creation timestamp
//   - Provided ID and code
func NewSmartContract(id string, code string) *SmartContract {
	return &SmartContract{
		ID:        id,
		Code:      code,
		State:     make(map[string]interface{}),
		CreatedAt: time.Now(),
	}
}

// Execute runs the smart contract with the provided input.
// Parameters:
//   - input: Map of input parameters for contract execution
//
// The function:
// 1. Logs the execution attempt
// 2. Stores the input in the contract's state
// 3. Returns the current state of the contract
//
// Returns:
//   - The contract's state after execution
//   - Any error that occurred during execution
//
// Note: This is a basic implementation that only stores the input.
// A full implementation would parse and execute the contract code.
func (sc *SmartContract) Execute(input map[string]interface{}) (map[string]interface{}, error) {
	fmt.Println("Executing smart contract", input)
	sc.State["lastExecution"] = input
	return sc.State, nil
}

// Validate performs basic validation of the smart contract.
// Checks:
//   - Contract code is not empty
//
// Returns:
//   - nil if validation passes
//   - error if validation fails
//
// Note: This is a basic validation. A full implementation would:
//   - Validate contract code syntax
//   - Check for security vulnerabilities
//   - Verify resource limits
//   - Validate state initialization
func (sc *SmartContract) Validate() error {
	if sc.Code == "" {
		return errors.New("smart contract code is required")
	}
	return nil
}
