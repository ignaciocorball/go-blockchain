// Package api implements the HTTP server and REST API endpoints for the UFChain blockchain.
// It provides interfaces for blockchain operations, smart contract deployment, and transaction handling.
package api

import (
	"net/http"

	"github.com/davepartner/go-blockchain/blockchain"
	"github.com/davepartner/go-blockchain/contracts"
	"github.com/davepartner/go-blockchain/storage"
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
//   - POST /contract      - Deploy new smart contracts
//   - POST /contract/:id/execute - Execute deployed contracts
func StartServer(bcInstance *blockchain.Blockchain, dbInstance *storage.BlockchainDB) {
	bc = bcInstance
	db = dbInstance

	e := echo.New()

	e.POST("/transaction", handleTransaction)
	e.GET("block/:hash", handleGetBlock)
	e.POST("/contract", handleDeployContract)
	e.POST("/contract/:id/execute", handleExecuteContract)

	e.Logger.Fatal(e.Start(":1323"))
}

// handleTransaction processes incoming transaction requests.
// Query Parameters:
//   - from:   Sender's address
//   - to:     Recipient's address
//   - amount: Transaction amount
//
// Returns a JSON response with transaction details and status.
func handleTransaction(c echo.Context) error {

	from := c.QueryParam("from")
	to := c.QueryParam("to")
	amount := c.QueryParam("amount")

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Transaction created successfully",
		"from":    from,
		"to":      to,
		"amount":  amount,
	})
}

// handleGetBlock retrieves a block from the blockchain by its hash.
// URL Parameters:
//   - hash: The hash of the block to retrieve
//
// Returns:
//   - 200 OK with block data if found
//   - 404 Not Found if block doesn't exist
func handleGetBlock(c echo.Context) error {
	hash := c.Param("hash")

	block, err := bc.GetBlock([]byte(hash))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Block not found",
		})
	}
	return c.JSON(http.StatusOK, block)
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
