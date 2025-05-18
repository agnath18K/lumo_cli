# Lumo Authentication System

This document provides detailed information about Lumo's authentication system for the REST API.

## Overview

Lumo's REST API includes a robust authentication system that protects all endpoints from unauthorized access. The authentication system uses JWT (JSON Web Tokens) for secure, stateless authentication and includes features like token refresh, password management, and secure credential storage.

## Default Credentials

When the authentication system is first initialized, a default user is created with the following credentials:

- **Username**: `admin`
- **Password**: `lumo`

**Important**: For security reasons, it is strongly recommended to change the default password immediately after the first login.

## Authentication Configuration

You can configure the authentication system using the following commands:

```bash
# Enable authentication for the REST server
lumo config:server auth enable

# Disable authentication for the REST server
lumo config:server auth disable

# Change the default admin password
lumo config:server auth password
```

The authentication settings are stored in the Lumo configuration file (`~/.config/lumo/config.json`) with the following options:

```json
{
  "enable_auth": true,
  "jwt_secret": "your-secret-key",
  "token_expiration_hours": 24,
  "refresh_expiration_days": 7
}
```

## Authentication Endpoints

The authentication system provides the following endpoints:

### Login

```bash
# Login to get a JWT token
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"lumo"}' \
  http://localhost:7531/api/v1/auth/login
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "admin",
  "expires_in": 86400
}
```

### Refresh Token

```bash
# Refresh an expired token
curl -X POST -H "Content-Type: application/json" \
  -d '{"refresh_token":"your-refresh-token"}' \
  http://localhost:7531/api/v1/auth/refresh
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "admin",
  "expires_in": 86400
}
```

### Change Password

```bash
# Change password (requires authentication)
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{"current_password":"lumo","new_password":"new-secure-password"}' \
  http://localhost:7531/api/v1/auth/change-password
```

Response:

```json
{
  "success": true,
  "message": "Password updated successfully"
}
```

## Using Authentication with API Endpoints

All API endpoints (except for `/ping` and `/api/v1/status`) require authentication when the authentication system is enabled. To authenticate, include the JWT token in the `Authorization` header:

```bash
# Execute a command with authentication
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{"command":"What is the capital of France?"}' \
  http://localhost:7531/api/v1/execute
```

## Web Interface Authentication

The web interface includes a login page that authenticates the user using the same credentials as the API. After successful authentication, the web interface stores the JWT token in the browser's localStorage and includes it in all API requests.

## Security Considerations

1. **Change Default Password**: Always change the default password immediately after the first login.
2. **Secure JWT Secret**: The JWT secret is automatically generated on first run, but you can set a custom secret in the configuration file.
3. **Token Expiration**: Tokens expire after 24 hours by default, but you can configure the expiration time in the configuration file.
4. **HTTPS**: For production use, it's recommended to use HTTPS to encrypt the communication between the client and the server.
5. **Firewall**: Configure your firewall to restrict access to the Lumo server port (7531 by default).

## Credential Storage

User credentials are stored locally in the `~/.config/lumo/credentials.json` file. Passwords are securely hashed using bcrypt with a cost factor of 12.

## Implementation Details

The authentication system is implemented using the following components:

1. **JWT Tokens**: JSON Web Tokens are used for stateless authentication.
2. **Bcrypt**: Passwords are hashed using bcrypt with a cost factor of 12.
3. **Middleware**: All API endpoints are protected by an authentication middleware.
4. **Local Storage**: Credentials are stored locally in the user's config directory.
5. **Token Refresh**: Refresh tokens are used to obtain new access tokens without requiring the user to log in again.

## Troubleshooting

If you encounter authentication issues, try the following:

1. **Check Credentials**: Verify that you're using the correct username and password.
2. **Check Token Expiration**: Tokens expire after 24 hours by default. Use the refresh token to obtain a new token.
3. **Check Server Status**: Make sure the Lumo server is running.
4. **Check Authentication Status**: Verify that authentication is enabled using `lumo config:server show`.
5. **Reset Credentials**: If you've forgotten your password, you can reset the credentials by deleting the `~/.config/lumo/credentials.json` file. This will recreate the default user on the next server start.

## Examples

### Python Example

```python
import requests
import json

# Base URL for the Lumo REST API
base_url = "http://localhost:7531"

# Login to get authentication tokens
login_payload = {
    "username": "admin",
    "password": "lumo"  # Replace with your actual password
}
login_response = requests.post(
    f"{base_url}/api/v1/auth/login",
    headers={"Content-Type": "application/json"},
    data=json.dumps(login_payload)
)
auth_data = login_response.json()
token = auth_data["token"]
refresh_token = auth_data["refresh_token"]
print(f"Logged in as: {auth_data['username']}")

# Execute a command with authentication
payload = {
    "command": "What is the capital of France?"
}
response = requests.post(
    f"{base_url}/api/v1/execute",
    headers={
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}"
    },
    data=json.dumps(payload)
)
print("Response:", response.json()["output"])
```

### JavaScript Example

```javascript
// Login to get authentication tokens
async function login(username, password) {
  const response = await fetch('http://localhost:7531/api/v1/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  });
  
  if (!response.ok) {
    throw new Error('Login failed');
  }
  
  return await response.json();
}

// Execute a command with authentication
async function executeCommand(token, command) {
  const response = await fetch('http://localhost:7531/api/v1/execute', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ command })
  });
  
  if (!response.ok) {
    throw new Error('Command execution failed');
  }
  
  return await response.json();
}

// Usage example
async function main() {
  try {
    const auth = await login('admin', 'lumo');
    console.log(`Logged in as: ${auth.username}`);
    
    const result = await executeCommand(auth.token, 'What is the capital of France?');
    console.log('Response:', result.output);
  } catch (error) {
    console.error('Error:', error.message);
  }
}

main();
```
