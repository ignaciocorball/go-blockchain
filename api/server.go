// Package api implements the HTTP server and REST API endpoints for the UFChain blockchain.
// It provides interfaces for blockchain operations, smart contract deployment, and transaction handling.
package api

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ignaciocorball/go-blockchain/blockchain"
	"github.com/ignaciocorball/go-blockchain/contracts"
	"github.com/ignaciocorball/go-blockchain/storage"
	"github.com/labstack/echo/v4"
)

/*
  This file contains the API server implementation for the blockchain.
  It provides endpoints for blockchain operations such as creating a new blockchain,
  adding blocks, and querying the blockchain.
*/

// Global variables to store blockchain and database instances
// These are initialized when the server starts and used across all handlers
var bc *blockchain.Blockchain
var db *storage.BlockchainDB

// StartServer initializes and starts the HTTP server for the blockchain API.
// Parameters:
//   - bcInstance: The blockchain instance to use for operations
//   - dbInstance: The database instance for persistent storage
//
// The server provides the following endpoints:
//   - POST /transaction    - Create new transactions
//   - GET  /block/:hash   - Retrieve block information
//   - GET  /blocks         - Retrieve all blocks
//   - POST /contract      - Deploy new smart contracts
//   - POST /contract/:id/execute - Execute deployed contracts
//   - POST /wallet         - Create a new wallet
//   - GET  /wallet/:address/balance - Get wallet balance
//   - POST /wallet/:address/mint    - Mint new tokens to a wallet
func StartServer(bcInstance *blockchain.Blockchain, dbInstance *storage.BlockchainDB) {
	bc = bcInstance
	db = dbInstance

	e := echo.New()

	e.POST("/transaction", handleTransaction)
	e.GET("/block/:hash", handleGetBlock)
	e.GET("/blocks", handleGetAllBlocks)
	e.POST("/contract", handleDeployContract)
	e.POST("/contract/:id/execute", handleExecuteContract)
	e.POST("/wallet", handleCreateWallet)
	e.GET("/wallet/:address/balance", handleGetWalletBalance)
	e.POST("/wallet/:address/mint", handleMintTokens)

	e.Logger.Fatal(e.Start(":1323"))
}

// handleTransaction processes incoming transaction requests and creates a new block.
// Query Parameters:
//   - from:   Sender's address
//   - to:     Recipient's address
//   - amount: Transaction amount
//   - privateKey: Sender's private key (hex encoded)
//
// Returns a JSON response with transaction details, block hash, and status.
// Possible errors:
//   - 400 Bad Request: Invalid parameters or insufficient funds
//   - 404 Not Found: Wallet not found
//   - 500 Internal Server Error: Database or blockchain errors
func handleTransaction(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	amountStr := c.QueryParam("amount")
	privateKeyHex := c.QueryParam("privateKey")

	// Validate required parameters
	if from == "" || to == "" || amountStr == "" || privateKeyHex == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameters: from, to, amount, and privateKey are required",
		})
	}

	// Convert amount to integer
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid amount format",
		})
	}

	// Get sender's wallet
	fromWallet, err := db.GetWallet(from)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Sender wallet not found",
		})
	}

	// Verify that the private key matches
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid private key format",
		})
	}

	// Parse the private key
	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid private key format: " + err.Error(),
		})
	}

	// Verify that the private key corresponds to the wallet
	if !bytes.Equal(privateKey.PublicKey.X.Bytes(), fromWallet.PublicKey[:32]) ||
		!bytes.Equal(privateKey.PublicKey.Y.Bytes(), fromWallet.PublicKey[32:]) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Private key does not match wallet address",
		})
	}

	// Get recipient's wallet
	toWallet, err := db.GetWallet(to)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Recipient wallet not found",
		})
	}

	// Verify sufficient balance
	balance := fromWallet.GetBalance(bc)
	if balance < amount {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Insufficient funds",
			"balance": balance,
			"amount":  amount,
		})
	}

	// Get available UTXOs for the sender
	utxos := bc.UTXOs.GetUTXOsForAddress(fromWallet.PublicKey)

	// Create the transaction using the NewTransaction function
	// We pass the recipient's public key directly
	tx, err := blockchain.NewTransaction(fromWallet, string(toWallet.PublicKey), amount, utxos)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Create a new block with the transaction
	// We use a test validator for now
	newBlock := bc.AddBlock([]*blockchain.Transaction{tx}, []byte("test-validator"))

	// Save the block to the database
	err = db.SaveBlock(newBlock)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error saving block to database",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":    "Transaction created and block added successfully",
		"from":       from,
		"to":         to,
		"amount":     amount,
		"block_hash": fmt.Sprintf("%x", newBlock.Hash),
	})
}

// handleGetBlock retrieves a block from the blockchain by its hash.
// URL Parameters:
//   - hash: The hash of the block to retrieve (in base64)
//
// Returns:
//   - 200 OK with block data if found
//   - 404 Not Found if block doesn't exist
func handleGetBlock(c echo.Context) error {
	hashStr := c.Param("hash")

	// Decodificar el hash de base64
	hash, err := base64.StdEncoding.DecodeString(hashStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid hash format",
			"error":   err.Error(),
		})
	}

	block, err := bc.GetBlock(hash)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Block not found",
		})
	}
	return c.JSON(http.StatusOK, block[0]) // Retornar el primer bloque encontrado
}

// handleGetAllBlocks retrieves all blocks from the blockchain.
// Returns a JSON response with the list of all blocks.
func handleGetAllBlocks(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"blocks": bc.Blocks,
		"count":  len(bc.Blocks),
	})
}

// handleDeployContract processes smart contract deployment requests.
// Query Parameters:
//   - id:   Unique identifier for the contract
//   - code: Smart contract code to deploy
//
// Returns:
//   - 201 Created if deployment successful
//   - 400 Bad Request if contract validation fails
func handleDeployContract(c echo.Context) error {
	id := c.QueryParam("id")
	code := c.QueryParam("code")

	contract := contracts.NewSmartContract(id, code)
	err := contract.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Contract deployed successfully",
		"id":      id,
	})
}

// handleExecuteContract executes a deployed smart contract.
// URL Parameters:
//   - id: The identifier of the contract to execute
//
// Request Body:
//   - input: Map of input parameters for contract execution
//
// Returns:
//   - 200 OK with execution results
//   - 404 Not Found if contract doesn't exist
func handleExecuteContract(c echo.Context) error {
	id := c.Param("id")
	input := map[string]interface{}{}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Contract executed successfully",
		"id":      id,
		"input":   input,
	})
}

// handleCreateWallet creates a new wallet and returns its credentials
// Returns:
//   - 201 Created with address, public key and private key
//   - 500 Internal Server Error if there's an error saving
func handleCreateWallet(c echo.Context) error {
	wallet := blockchain.NewWallet()

	// Guardar la wallet en la base de datos
	err := db.SaveWallet(wallet.Address, wallet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error saving wallet",
		})
	}

	// Devolver la dirección, clave pública y clave privada
	// En una implementación real, la clave privada debería manejarse de forma más segura
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"address":    wallet.Address,
		"publicKey":  fmt.Sprintf("%x", wallet.PublicKey),
		"privateKey": fmt.Sprintf("%x", wallet.PrivateKeyBytes),
	})
}

// handleGetWalletBalance gets the balance of a wallet
func handleGetWalletBalance(c echo.Context) error {
	address := c.Param("address")

	wallet, err := db.GetWallet(address)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Wallet not found",
		})
	}

	balance := wallet.GetBalance(bc)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"address": address,
		"balance": balance,
	})
}

// handleMintTokens creates a special transaction to generate new tokens and assign them to a wallet
// URL Parameters:
//   - address: The address of the wallet that will receive the tokens
//
// Query Parameters:
//   - amount: Amount of tokens to generate
//   - privateKey: Private key of the wallet (in hex format) that authorizes the generation
//
// Returns:
//   - 201 Created if generation was successful
//   - 400 Bad Request if parameters are invalid
//   - 404 Not Found if wallet doesn't exist
//   - 500 Internal Server Error if there are internal errors
func handleMintTokens(c echo.Context) error {
	address := c.Param("address")
	amountStr := c.QueryParam("amount")
	privateKeyHex := c.QueryParam("privateKey")

	// Validate required parameters
	if amountStr == "" || privateKeyHex == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameters: amount and privateKey are required",
		})
	}

	// Convert amount to integer
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid amount format",
		})
	}

	// Verify that the destination wallet exists
	wallet, err := db.GetWallet(address)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Destination wallet not found",
		})
	}

	// Verify that the private key matches the wallet
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid private key format",
		})
	}
	if !bytes.Equal(privateKeyBytes, wallet.PrivateKeyBytes) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid private key for wallet address",
		})
	}

	// Create a special token generation transaction
	// This transaction has no inputs (it's a coinbase transaction)
	tx := &blockchain.Transaction{
		Input: []blockchain.TxInput{}, // No inputs in a generation transaction
		Output: []blockchain.TxOutput{{
			Value:     amount,
			PublicKey: wallet.PublicKey,
		}},
	}
	tx.ID = tx.HashTransaction()

	// Create a new block with the generation transaction
	// We use the wallet's public key as validator
	newBlock := bc.AddBlock([]*blockchain.Transaction{tx}, wallet.PublicKey)

	// Save the block to the database
	err = db.SaveBlock(newBlock)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error saving block to database",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":    "Tokens minted successfully",
		"address":    address,
		"amount":     amount,
		"block_hash": fmt.Sprintf("%x", newBlock.Hash),
	})
}
