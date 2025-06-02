package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"

	"github.com/omkar273/police/internal/config"
	ierr "github.com/omkar273/police/internal/errors"
	"github.com/omkar273/police/internal/logger"
)

// EncryptionService defines the interface for encryption and hashing operations
type EncryptionService interface {
	// Encrypt encrypts plaintext using RSA public key
	Encrypt(plaintext string) (string, error)

	// Decrypt decrypts ciphertext using RSA private key
	Decrypt(ciphertext string) (string, error)

	// Hash creates a one-way hash of the input value using SHA-256
	Hash(value string) string
}

type rsaEncryptionService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	logger     *logger.Logger
}

// NewEncryptionService creates a new encryption service using RSA keys from config
func NewEncryptionService(cfg *config.Configuration, logger *logger.Logger) (EncryptionService, error) {
	if cfg.Secrets.PrivateKey == "" || cfg.Secrets.PublicKey == "" {
		return nil, ierr.NewError("RSA keys not configured").
			WithHint("Both private and public RSA keys must be configured").
			Mark(ierr.ErrSystem)
	}

	// Parse private key
	privateKeyPEM := cfg.Secrets.PrivateKey
	privateKeyBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if privateKeyBlock == nil {
		return nil, ierr.NewError("failed to decode private key PEM").
			WithHint("Private key must be in valid PEM format").
			Mark(ierr.ErrSystem)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		// Try PKCS8 format if PKCS1 fails
		privateKeyInterface, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			return nil, ierr.WithError(err).
				WithHint("Failed to parse private key - ensure it's in PKCS1 or PKCS8 format").
				Mark(ierr.ErrSystem)
		}
		var ok bool
		privateKey, ok = privateKeyInterface.(*rsa.PrivateKey)
		if !ok {
			return nil, ierr.NewError("private key is not RSA").
				WithHint("Private key must be an RSA key").
				Mark(ierr.ErrSystem)
		}
	}

	// Parse public key
	publicKeyPEM := cfg.Secrets.PublicKey
	publicKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if publicKeyBlock == nil {
		return nil, ierr.NewError("failed to decode public key PEM").
			WithHint("Public key must be in valid PEM format").
			Mark(ierr.ErrSystem)
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to parse public key").
			Mark(ierr.ErrSystem)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, ierr.NewError("public key is not RSA").
			WithHint("Public key must be an RSA key").
			Mark(ierr.ErrSystem)
	}

	return &rsaEncryptionService{
		privateKey: privateKey,
		publicKey:  publicKey,
		logger:     logger,
	}, nil
}

// Encrypt encrypts plaintext using RSA public key with OAEP padding
func (s *rsaEncryptionService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// RSA encryption with OAEP padding and SHA-256 hash
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		s.publicKey,
		[]byte(plaintext),
		nil,
	)
	if err != nil {
		return "", ierr.WithError(err).
			WithHint("Failed to encrypt data with RSA public key").
			Mark(ierr.ErrSystem)
	}

	// Encode the result as base64 for storage
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

// Decrypt decrypts base64-encoded ciphertext using RSA private key
func (s *rsaEncryptionService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Decode the base64-encoded ciphertext
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ierr.WithError(err).
			WithHint("Failed to decode ciphertext").
			Mark(ierr.ErrSystem)
	}

	// RSA decryption with OAEP padding and SHA-256 hash
	plaintext, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		s.privateKey,
		decoded,
		nil,
	)
	if err != nil {
		return "", ierr.WithError(err).
			WithHint("Failed to decrypt data with RSA private key").
			Mark(ierr.ErrSystem)
	}

	return string(plaintext), nil
}

// Hash creates a one-way hash of the input value using SHA-256
func (s *rsaEncryptionService) Hash(value string) string {
	if value == "" {
		return ""
	}

	// Create a new SHA-256 hasher
	hasher := sha256.New()

	// Write the value to the hasher
	hasher.Write([]byte(value))

	// Get the hash sum and convert to hex string
	return hex.EncodeToString(hasher.Sum(nil))
}
