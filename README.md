# CA Service - SSH Certificate Authority

A lightweight SSH Certificate Authority service that signs ED25519 public keys and issues short-lived SSH certificates.

## What It Does

This service acts as a Certificate Authority (CA) that:
- Accepts SSH public keys (ED25519 only)
- Signs them with a CA private key
- Issues certificates valid for 30 minutes
- Requires authentication via bearer tokens

Basically it will issue short-lived SSH certificates to your infrastructure instead of relying on traditional key distribution.

## How It Works

1. **Client** sends a public key + authentication token
2. **Service** validates the token and public key
3. **Service** signs the key with its CA private key
4. **Client** receives a signed certificate to use for SSH authentication

## Prerequisites

- Go 1.25+ (if building from source)
- A CA private key pair (ED25519 format)
- HTTPS certificates for the server

## Setup

### 1. Generate CA Key Pair

If you don't have CA keys, generate them:

```bash
ssh-keygen -t ed25519 -f ./certs/ca-key-pair/ca_key -N ""
```

This creates:
- `./certs/ca-key-pair/ca_key`
- `./certs/ca-key-pair/ca_key.pub`

### 2. Generate HTTPS Certificates

The service runs over HTTPS. Generate self-signed certificates:

```bash
# Using mkcert (for local development)
mkcert -key-file ./certs/https/localhost+2-key.pem -cert-file ./certs/https/localhost+2.pem localhost 127.0.0.1 ::1
```

This creates:
- `./certs/https/localhost+2-key.pem` (private key)
- `./certs/https/localhost+2.pem` (certificate)

### 3. Create `.env` File

Create a `.env` file in the project root with your authentication tokens.
Tokens are used in the Authorization header.
Principals are identities the certificate will be valid for

```bash
# Format: token:principal1,principal2;token2:principal3
```

Example with real tokens:
```bash
CA_ACCESS_TOKEN=prod_abc123:admin,root;test_xyz789:test-user
```

## Usage

### Start the Server

```bash
go run ./cmd/CA-service/main.go
```

The service will start on `https://localhost:8443`

### Request Format

Sign a public key by sending a POST request to `/sign`:

```bash
curl -k \
  -X POST https://localhost:8443/sign \
  -H "Authorization: Bearer your-secret-token" \
  -H "Content-Type: application/json" \
  -d '{
    "public_key": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFmj... user@host"
  }'
```

**Parameters:**
- `Authorization` header: Bearer token (must match a token in `CA_ACCESS_TOKEN`)
- `public_key`: SSH public key in OpenSSH format (ED25519 only)

### Response Format

**Success (200 OK):**
```json
{
  "signed_cert": "ssh-ed25519-cert-v01@openssh.com AAAAHHNzaC1lZDI1NTE5LWNlcnQtdjAxQG9wZW5zc2guY29t..."
}
```

**Error (400+):**
```json
{
  "error": "err msg example"
}
```

## Building

```bash
make build
# or
go build -o bin/CA-service ./cmd/CA-service
```

## Tests

Run tests:
```bash
go test ./...
```
