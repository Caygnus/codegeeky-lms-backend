# Configuration Guide

This application requires proper configuration to run. The configuration system supports multiple sources with the following priority order:

1. **Environment Variables** (highest priority)
2. **Config YAML file**
3. **.env file**
4. **Key files in ./keys/ folder** (lowest priority)

## Required Configuration

The following fields are **required** and the application will fail to start if any are missing:

### Server Configuration

- `server.address` - Server bind address (e.g., `:8080`)

### Logging Configuration

- `logging.level` - Log level (`debug`, `info`, `warn`, `error`)

### PostgreSQL Configuration

- `postgres.host` - Database host
- `postgres.port` - Database port
- `postgres.user` - Database username
- `postgres.password` - Database password
- `postgres.dbname` - Database name
- `postgres.sslmode` - SSL mode (`disable`, `require`, `verify-ca`, `verify-full`)

### Supabase Configuration

- `supabase.url` - Supabase project URL
- `supabase.key` - Supabase API key
- `supabase.jwt_secret` - Supabase JWT secret for local token validation

### Secrets Configuration

- `secrets.private_key` - Private key for cryptographic operations
- `secrets.public_key` - Public key for cryptographic operations

## Configuration Methods

### Method 1: YAML Configuration File

Create a `config.yaml` file in either `./internal/config/` or `./config/` directory:

```yaml
server:
  address: ":8080"

logging:
  level: "debug"

postgres:
  host: "localhost"
  port: 5432
  user: "your_username"
  password: "your_password"
  dbname: "your_database"
  sslmode: "disable"
  max_open_conns: 10
  max_idle_conns: 5
  conn_max_lifetime_minutes: 60
  auto_migrate: true

supabase:
  url: "https://your-project.supabase.co"
  key: "your-supabase-key"
  jwt_secret: "your-supabase-jwt-secret"

secrets:
  private_key: |
    -----BEGIN PRIVATE KEY-----
    your_private_key_content_here
    -----END PRIVATE KEY-----
  public_key: |
    -----BEGIN PUBLIC KEY-----
    your_public_key_content_here
    -----END PUBLIC KEY-----
```

### Method 2: Environment Variables

Set environment variables with the `CAYGNUS_` prefix:

```bash
export CAYGNUS_SERVER_ADDRESS=":8080"
export CAYGNUS_LOGGING_LEVEL="debug"
export CAYGNUS_POSTGRES_HOST="localhost"
export CAYGNUS_POSTGRES_PORT=5432
export CAYGNUS_POSTGRES_USER="your_username"
export CAYGNUS_POSTGRES_PASSWORD="your_password"
export CAYGNUS_POSTGRES_DBNAME="your_database"
export CAYGNUS_POSTGRES_SSLMODE="disable"
export CAYGNUS_POSTGRES_AUTO_MIGRATE=true
export CAYGNUS_SUPABASE_URL="https://your-project.supabase.co"
export CAYGNUS_SUPABASE_KEY="your-supabase-key"
export CAYGNUS_SUPABASE_JWT_SECRET="your-supabase-jwt-secret"
export CAYGNUS_SECRETS_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nyour_private_key_content_here\n-----END PRIVATE KEY-----"
export CAYGNUS_SECRETS_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\nyour_public_key_content_here\n-----END PUBLIC KEY-----"
```

### Method 3: .env File

Create a `.env` file in the project root:

```env
CAYGNUS_SERVER_ADDRESS=:8080
CAYGNUS_LOGGING_LEVEL=debug
CAYGNUS_POSTGRES_HOST=localhost
CAYGNUS_POSTGRES_PORT=5432
CAYGNUS_POSTGRES_USER=your_username
CAYGNUS_POSTGRES_PASSWORD=your_password
CAYGNUS_POSTGRES_DBNAME=your_database
CAYGNUS_POSTGRES_SSLMODE=disable
CAYGNUS_POSTGRES_AUTO_MIGRATE=true
CAYGNUS_SUPABASE_URL=https://your-project.supabase.co
CAYGNUS_SUPABASE_KEY=your-supabase-key
CAYGNUS_SUPABASE_JWT_SECRET=your-supabase-jwt-secret
CAYGNUS_SECRETS_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nyour_private_key_content_here\n-----END PRIVATE KEY-----"
CAYGNUS_SECRETS_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----\nyour_public_key_content_here\n-----END PUBLIC KEY-----"
```

### Method 4: Key Files (Recommended for Development)

Create key files in the `./keys/` folder. The application will automatically load them if secrets are not provided via other methods:

**Primary file names (checked first):**

- `./keys/private_key.pem` - Private key file
- `./keys/public_key.pem` - Public key file

**Alternative file names (fallback):**

- `./keys/private.pem` or `./keys/private_key`
- `./keys/public.pem` or `./keys/public_key`
- `./keys/id_rsa` (private key)
- `./keys/id_rsa.pub` (public key)

Example key files:

**./keys/private_key.pem:**

```
-----BEGIN PRIVATE KEY-----
your_private_key_content_here
-----END PRIVATE KEY-----
```

**./keys/public_key.pem:**

```
-----BEGIN PUBLIC KEY-----
your_public_key_content_here
-----END PUBLIC KEY-----
```

## Optional Configuration

The following fields have default values and are optional:

- `postgres.max_open_conns` (default: 10)
- `postgres.max_idle_conns` (default: 5)
- `postgres.conn_max_lifetime_minutes` (default: 60)
- `postgres.auto_migrate` (default: false)

## Validation

The application performs strict validation on startup and will fail with detailed error messages if:

1. Required fields are missing
2. Field values don't meet validation criteria
3. No configuration source is available

## Error Handling

If configuration validation fails, you'll see an error message like:

```
configuration validation failed: Key: 'Configuration.Server.Address' Error:Field validation for 'Address' failed on the 'required' tag

Please ensure you have either:
1. A valid config.yaml file in ./internal/config/ or ./config/
2. A .env file with required variables
3. Environment variables with CAYGNUS_ prefix
4. Key files in ./keys/ folder (private_key.pem, public_key.pem)

Required fields: server.address, logging.level, postgres.host, postgres.port, postgres.user, postgres.password, postgres.dbname, postgres.sslmode, secrets.private_key, secrets.public_key
```

This ensures the application fails fast with clear guidance on what needs to be configured.

## Key Management Best Practices

1. **For Development**: Use key files in the `./keys/` folder (add `keys/` to `.gitignore`)
2. **For Production**: Use environment variables or secure secret management systems
3. **Never commit keys to version control**
4. **Use proper file permissions** (600) for key files
5. **Rotate keys regularly** in production environments
