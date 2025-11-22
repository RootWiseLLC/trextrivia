# Coolify Environment Variables Setup

When setting environment variables in Coolify, you need to paste the **raw PEM key content** (not base64 encoded).

## JWT_PRIVATE_KEY

Paste the entire contents of your `jwtRS512.key` file, including the header and footer:

```
-----BEGIN RSA PRIVATE KEY-----
MIIJKgIBAAKCAgEAuymgeZj4WohZ7JLPIIJUEzrWjwlUoaajnA6C3SjnSjlH4r8j
7c/ANEM7VP+Onn6h1E4CzIruCz+OkaHhfrkJrB6kiiV0byfxBGBOnBebMOt6EUgF
... (rest of your private key)
-----END RSA PRIVATE KEY-----
```

## JWT_PUBLIC_KEY

Paste the entire contents of your `jwtRS512.key.pub` file:

```
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAuymgeZj4WohZ7JLPIIJu
... (rest of your public key)
-----END PUBLIC KEY-----
```

## How to get the keys:

```bash
# View private key
cat jwtRS512.key

# View public key  
cat jwtRS512.key.pub
```

Copy the ENTIRE output (including `-----BEGIN` and `-----END` lines) and paste into Coolify.

## Important Notes:

- **DO NOT base64 encode** - Paste the raw PEM format
- **Include newlines** - The multi-line format is required
- **Include headers** - The `-----BEGIN` and `-----END` lines are part of the key
- Coolify will handle multi-line values correctly

## Other Environment Variables:

```bash
POSTGRES_PASSWORD=your_strong_random_password
DOMAIN=trextrivia.com
```

These are simple single-line values.
