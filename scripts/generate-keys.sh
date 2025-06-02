#!/bin/bash

# Script to generate RSA key pairs for development
# This script creates private and public keys in the keys/ folder

set -e

KEYS_DIR="keys"
PRIVATE_KEY_FILE="$KEYS_DIR/private_key.pem"
PUBLIC_KEY_FILE="$KEYS_DIR/public_key.pem"

echo "🔐 Generating RSA key pair for development..."

# Create keys directory if it doesn't exist
mkdir -p "$KEYS_DIR"

# Generate private key (2048-bit RSA)
echo "📝 Generating private key..."
openssl genpkey -algorithm RSA -out "$PRIVATE_KEY_FILE" -pkcs8 -aes256

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
echo "⚠️  IMPORTANT: The private key is password-protected."
echo "   You'll need to enter the password when the application loads the key."
echo ""
echo "🔒 For production, consider using environment variables or a secure key management system." 