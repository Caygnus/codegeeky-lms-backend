# Keys Directory

This directory is used to store cryptographic keys for the application. The keys are automatically loaded by the configuration system if they are not provided via environment variables or configuration files.

## Quick Start

### For Development (Unencrypted Keys)

```bash
make generate-dev-keys
```

### For Production (Encrypted Keys)

```bash
make generate-keys
```

## Supported File Names

The application will automatically look for keys in the following order:

### Private Key

1. `private_key.pem` (preferred)
2. `private.pem`
3. `private_key`
4. `id_rsa`

### Public Key

1. `public_key.pem` (preferred)
2. `public.pem`
3. `public_key`
4. `id_rsa.pub`

## Manual Key Generation

If you prefer to generate keys manually:

### RSA Key Pair (Unencrypted - Development Only)

```bash
# Generate private key
openssl genpkey -algorithm RSA -out private_key.pem -pkcs8

# Extract public key
openssl pkey -in private_key.pem -pubout -out public_key.pem

# Set proper permissions
chmod 600 private_key.pem
chmod 644 public_key.pem
```

### RSA Key Pair (Encrypted - Production)

```bash
# Generate encrypted private key
openssl genpkey -algorithm RSA -out private_key.pem -pkcs8 -aes256

# Extract public key
openssl pkey -in private_key.pem -pubout -out public_key.pem

# Set proper permissions
chmod 600 private_key.pem
chmod 644 public_key.pem
```

## Security Notes

⚠️ **IMPORTANT**: This directory is ignored by Git for security reasons.

- **Never commit keys to version control**
- **Use encrypted keys in production**
- **Set proper file permissions** (600 for private keys, 644 for public keys)
- **Rotate keys regularly** in production environments
- **For production deployments**, consider using environment variables or a secure key management system instead of files

## Alternative Configuration Methods

If you don't want to use files, you can also provide keys via:

1. **Environment Variables**:

   ```bash
   export CAYGNUS_SECRETS_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
   export CAYGNUS_SECRETS_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
   ```

2. **Configuration File** (`config.yaml`):

   ```yaml
   secrets:
     private_key: |
       -----BEGIN PRIVATE KEY-----
       ...
       -----END PRIVATE KEY-----
     public_key: |
       -----BEGIN PUBLIC KEY-----
       ...
       -----END PUBLIC KEY-----
   ```

3. **.env File**:
   ```env
   CAYGNUS_SECRETS_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
   CAYGNUS_SECRETS_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
   ```
