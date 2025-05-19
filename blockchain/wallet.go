package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/hex"
	"log"
)

// Wallet represents a wallet in the blockchain
type Wallet struct {
	PrivateKeyBytes []byte // Serialized private key
	PublicKey       []byte // Public key
	Address         string // Address derived from the public key
}

// NewWallet creates a new wallet with an ECDSA key pair
func NewWallet() *Wallet {
	// Use P-256 curve for keys
	curve := elliptic.P256()

	// Generate private/public key pair
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	// Serialize private key to PEM format
	privateKeyBytes, err := x509.MarshalECPrivateKey(private)
	if err != nil {
		log.Panic(err)
	}

	// Get public key in bytes format
	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	// Generate address from public key
	address := generateAddress(publicKey)

	return &Wallet{
		PrivateKeyBytes: privateKeyBytes,
		PublicKey:       publicKey,
		Address:         address,
	}
}

// GetPrivateKey retrieves the ECDSA private key
func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := x509.ParseECPrivateKey(w.PrivateKeyBytes)
	if err != nil {
		log.Panic(err)
	}
	return privateKey
}

// generateAddress creates a unique and readable address from the public key
func generateAddress(publicKey []byte) string {
	// Hash of the public key
	hash := sha256.Sum256(publicKey)

	// Take the last 20 bytes
	addressBytes := hash[len(hash)-20:]

	// Convert to hexadecimal and add "0x" prefix
	return "0x" + hex.EncodeToString(addressBytes)
}

// GetBalance calculates the wallet's balance
func (w *Wallet) GetBalance(bc *Blockchain) int {
	return bc.GetBalance(w.PublicKey)
}

// Serialize converts the wallet to bytes for storage
func (w *Wallet) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(w)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeWallet reconstructs a wallet from bytes
func DeserializeWallet(data []byte) *Wallet {
	var wallet Wallet

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&wallet)
	if err != nil {
		log.Panic(err)
	}

	return &wallet
}
