#!/bin/bash

# Development Authentication Example Script
# This script demonstrates how to use the development authentication solutions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Development Authentication Examples${NC}"
echo "=================================="
echo

# Check if server is running
echo -e "${YELLOW}1. Checking if server is running...${NC}"
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Server is running on http://localhost:8080${NC}"
else
    echo -e "${RED}❌ Server is not running. Please start your server first.${NC}"
    echo "   Run: go run cmd/server/main.go"
    exit 1
fi

echo

# Example 1: Development API Key
echo -e "${YELLOW}2. Testing with Development API Key${NC}"
echo "----------------------------------------"

# Check if dev API key is configured
if [ -z "$CAYGNUS_SUPABASE_DEV_API_KEY" ]; then
    echo -e "${RED}❌ CAYGNUS_SUPABASE_DEV_API_KEY not set${NC}"
    echo "   Add to your .env file:"
    echo "   CAYGNUS_SUPABASE_DEV_API_KEY=your_dev_key_here"
    echo
else
    echo -e "${GREEN}✅ Development API key found${NC}"
    echo "Testing authenticated endpoint..."
    
    response=$(curl -s -w "\n%{http_code}" \
        -H "Authorization: Bearer $CAYGNUS_SUPABASE_DEV_API_KEY" \
        http://localhost:8080/api/v1/user/me)
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✅ Success! (HTTP $http_code)${NC}"
        echo "Response: $body"
    else
        echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
        echo "Response: $body"
    fi
fi

echo

# Example 2: Long-lived JWT Token
echo -e "${YELLOW}3. Testing with Long-lived JWT Token${NC}"
echo "----------------------------------------"

if [ -z "$DEV_JWT_TOKEN" ]; then
    echo -e "${RED}❌ DEV_JWT_TOKEN not set${NC}"
    echo "   To generate a token, run:"
    echo "   make generate-dev-token JWT_SECRET=your_jwt_secret"
    echo "   Then add the token to your .env file:"
    echo "   DEV_JWT_TOKEN=your_generated_token_here"
    echo
else
    echo -e "${GREEN}✅ Development JWT token found${NC}"
    echo "Testing authenticated endpoint..."
    
    response=$(curl -s -w "\n%{http_code}" \
        -H "Authorization: Bearer $DEV_JWT_TOKEN" \
        http://localhost:8080/api/v1/user/me)
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✅ Success! (HTTP $http_code)${NC}"
        echo "Response: $body"
    else
        echo -e "${RED}❌ Failed (HTTP $http_code)${NC}"
        echo "Response: $body"
    fi
fi

echo

# Example 3: Without authentication (should fail)
echo -e "${YELLOW}4. Testing without authentication (should fail)${NC}"
echo "------------------------------------------------"

response=$(curl -s -w "\n%{http_code}" \
    http://localhost:8080/api/v1/user/me)

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n -1)

if [ "$http_code" = "401" ]; then
    echo -e "${GREEN}✅ Correctly rejected (HTTP $http_code)${NC}"
    echo "Response: $body"
else
    echo -e "${RED}❌ Unexpected response (HTTP $http_code)${NC}"
    echo "Response: $body"
fi

echo

# Summary
echo -e "${BLUE}📋 Summary${NC}"
echo "=========="
echo -e "${GREEN}✅ Development API Key:${NC} Fast, simple, no expiration"
echo -e "${GREEN}✅ Long-lived JWT Token:${NC} Real JWT validation, 1 year validity"
echo -e "${GREEN}✅ Both solutions:${NC} Only work in development mode"
echo
echo -e "${YELLOW}💡 Tips:${NC}"
echo "- Use Development API Key for quick testing"
echo "- Use Long-lived JWT Token for realistic testing"
echo "- Never use development tokens in production"
echo "- Keep your JWT secrets secure"
echo
echo -e "${BLUE}📚 Documentation:${NC}"
echo "- See docs/DEVELOPMENT_AUTHENTICATION.md for detailed guide"
echo "- Run 'make generate-dev-token JWT_SECRET=your_secret' to generate tokens" 