# Supabase JWT Token Validation

This document describes the implementation of local JWT token validation for Supabase authentication, replacing the previous HTTP-based token validation.

## Overview

The `ValidateToken` method in the Supabase authentication provider has been modified to validate Supabase access tokens locally by decoding and verifying JWT tokens instead of making HTTP calls to Supabase's API.

## Benefits

- **Performance**: Eliminates HTTP calls for token validation, reducing latency
- **Reliability**: Removes dependency on Supabase API availability
- **Scalability**: Local validation scales better under high load
- **Cost**: Reduces API usage and potential rate limiting

## Implementation Details

### JWT Validation Process

1. **Parse JWT Token**: Uses the `github.com/golang-jwt/jwt/v5` library to parse the token
2. **Verify Signature**: Validates the token signature using the Supabase JWT secret
3. **Check Expiration**: Automatically validates token expiration
4. **Validate Claims**: Ensures required claims are present and valid
5. **Extract User Info**: Extracts user ID, email, and phone from token claims

### Required Configuration

Add the JWT secret to your configuration:

```yaml
supabase:
  url: "https://your-project.supabase.co"
  key: "your-supabase-key"
  jwt_secret: "your-supabase-jwt-secret"
```

Or via environment variable:

```bash
export CAYGNUS_SUPABASE_JWT_SECRET="your-supabase-jwt-secret"
```

### Finding Your JWT Secret

1. Go to your Supabase Dashboard
2. Navigate to Settings > API
3. Copy the "JWT Secret" value
4. Add it to your configuration

## Token Validation Rules

The implementation validates the following:

### Required Claims

- `sub`: User ID (must be present and non-empty)
- `aud`: Audience (must be "authenticated" for user tokens)
- `role`: Role (must be "authenticated" for user tokens)
- `exp`: Expiration (automatically validated by JWT library)

### Optional Claims

- `email`: User email address
- `phone`: User phone number
- `iat`: Issued at timestamp

### Supported Algorithms

- **HS256**: HMAC SHA-256 (Supabase default)

## Usage Example

```go
// The ValidateToken method is called automatically by your auth middleware
claims, err := authProvider.ValidateToken(ctx, accessToken)
if err != nil {
    // Token is invalid
    return err
}

// Use the validated claims
userID := claims.UserID
email := claims.Email
phone := claims.Phone
```

## Error Handling

The implementation returns specific errors for different validation failures:

- `invalid token: signature is invalid` - Wrong JWT secret or tampered token
- `token is invalid` - General token validation failure
- `token is expired` - Token has expired
- `invalid token audience` - Wrong audience claim
- `invalid token role` - Wrong role claim
- `token missing user ID` - Missing or empty sub claim

## Testing

Comprehensive tests are included in `internal/auth/supabase_test.go`:

```bash
go test ./internal/auth -v
```

Test scenarios include:

- Valid JWT token validation
- Invalid JWT secret handling
- Missing user ID detection
- Invalid audience rejection
- Expired token rejection

## Migration from HTTP Validation

### Before (HTTP-based)

```go
func (p *supabaseProvider) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
    claims, err := p.supabase.Admin.GetUser(ctx, token)
    if err != nil {
        return nil, err
    }
    return &auth.Claims{
        UserID: claims.ID,
    }, nil
}
```

### After (JWT-based)

```go
func (p *supabaseProvider) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
    // Parse and validate JWT locally
    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method and return secret
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(p.cfg.Supabase.JWTSecret), nil
    })

    // Validate and extract claims...
}
```

## Security Considerations

1. **JWT Secret Security**: Keep your JWT secret secure and never expose it in client-side code
2. **Token Expiration**: JWT tokens have built-in expiration validation
3. **Signature Verification**: All tokens are cryptographically verified
4. **Claim Validation**: Required claims are validated to ensure token integrity

## Troubleshooting

### Common Issues

1. **"signature is invalid"**

   - Check that your JWT secret is correct
   - Ensure the secret matches your Supabase project

2. **"invalid token audience"**

   - Ensure you're using user access tokens, not anon keys
   - Check that the token was issued for authenticated users

3. **"token is expired"**
   - Token has expired and needs to be refreshed
   - Check your token refresh logic

### Debug Logging

Enable debug logging to see detailed validation information:

```yaml
logging:
  level: "debug"
```

## Dependencies

- `github.com/golang-jwt/jwt/v5`: JWT parsing and validation library

## References

- [Supabase JWT Documentation](https://supabase.com/docs/guides/auth/jwts)
- [JWT.io Debugger](https://jwt.io/) - For testing and debugging JWT tokens
- [golang-jwt/jwt Documentation](https://github.com/golang-jwt/jwt)
