package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadPublicKey(filename string) (*ecdsa.PublicKey, error) {
	const op = "lib.jwt.LoadPublicKey"
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// decode pem block
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("%s: error decode pem block", op)
	}

	// Parse the public key
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%s: not an ECDSA public key", op)
	}

	return publicKey, nil
}

func LoadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	// Read the private key file
	const op = "lib.jwt.LoadPrivateKey"
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("%s: failed to decode PEM block containing EC private key", op)
	}

	// Parse the EC private key
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return privateKey, nil
}
