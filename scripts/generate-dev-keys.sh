#!/bin/bash

# Script to generate unencrypted RSA key pairs for development
# This script creates private and public keys in the keys/ folder
# WARNING: These keys are NOT password-protected - use only for development!

set -e

KEYS_DIR="keys"
PRIVATE_KEY_FILE="$KEYS_DIR/private_key.pem"
PUBLIC_KEY_FILE="$KEYS_DIR/public_key.pem"

echo "🔐 Generating RSA key pair for development (unencrypted)..."

# Create keys directory if it doesn't exist
mkdir -p "$KEYS_DIR"

# Generate private key (2048-bit RSA, unencrypted)
echo "📝 Generating private key..."
openssl genpkey -algorithm RSA -out "$PRIVATE_KEY_FILE" -pkcs8

# Extract public key from private key
echo "📝 Extracting public key..."
openssl pkey -in "$PRIVATE_KEY_FILE" -pubout -out "$PUBLIC_KEY_FILE"

# Set proper permissions
chmod 600 "$PRIVATE_KEY_FILE"
chmod 644 "$PUBLIC_KEY_FILE"

echo "✅ Key pair generated successfully!"
echo "   Private key: $PRIVATE_KEY_FILE"
echo "   Public key:  $PUBLIC_KEY_FILE"
echo ""
echo "⚠️  WARNING: These keys are NOT password-protected!"
echo "   Use only for development. Never use in production."
echo ""
echo "🔒 For production, use the generate-keys.sh script or a secure key management system." 