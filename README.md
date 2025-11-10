# Integration Service Demo

A demonstration service that simulates an integration flow with multiple service dependencies: authentication, database, and notification services.

## Overview

This service demonstrates a typical integration pattern where:
1. **Receive Request** - Client sends a request to the integration endpoint
2. **Validate Auth** - Request is authenticated via auth service
3. **Fetch Data** - Data is retrieved from database service
4. **Send Notification** - Notification is sent to user
5. **Respond** - Final response is sent back to client

## Architecture

The service includes:
- **Main Integration Endpoint** (`/api/process`) - Orchestrates the entire flow
- **Auth Service** (`/auth/validate`) - Simulated authentication service
- **Database Service** (`/database/fetch`) - Simulated database lookup
- **Notification Service** (`/notification/send`) - Simulated notification sender
- **Health Check** (`/health`) - Service health status

## Running the Service

### Using Go

```bash
go run main.go
```

### Using Docker

```bash
docker build -t integration-demo .
docker run -p 9090:9090 integration-demo
```

The service will start on port 9090.

## API Usage

### Main Integration Endpoint

**POST** `/api/process`

Request body:
```json
{
  "user_id": "user123",
  "token": "valid_token123",
  "action": "fetch_user_data",
  "metadata": {
    "source": "web",
    "version": "1.0"
  }
}
```

**Note**: Tokens must start with `valid_` to pass authentication.

Success response:
```json
{
  "success": true,
  "message": "Request processed successfully for action: fetch_user_data",
  "data": {
    "user_id": "user123",
    "action": "fetch_user_data",
    "records": 42,
    "last_access": "2025-11-09T12:00:00Z",
    "permissions": ["read", "write", "execute"],
    "quota_remaining": 856
  },
  "processed_at": "2025-11-10T10:30:00Z",
  "request_id": "REQ-1699612345-1234"
}
```

Error response:
```json
{
  "success": false,
  "message": "Authentication failed",
  "processed_at": "2025-11-10T10:30:00Z",
  "request_id": "REQ-1699612345-1234"
}
```

### Example cURL Commands

**Valid request:**
```bash
curl -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "token": "valid_token123",
    "action": "fetch_user_data"
  }'
```

**Invalid token (will fail auth):**
```bash
curl -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "token": "invalid_token",
    "action": "fetch_user_data"
  }'
```

**Health check:**
```bash
curl http://localhost:9090/health
```

## Individual Service Endpoints

### Auth Service
**POST** `/auth/validate`
```json
{
  "token": "valid_token123",
  "user_id": "user123"
}
```

### Database Service
**GET** `/database/fetch?user_id=user123&action=fetch_user_data`

### Notification Service
**POST** `/notification/send`
```json
{
  "user_id": "user123",
  "message": "Your action completed successfully"
}
```

## Logging

The service provides detailed logging for each request:
- Request ID for tracking
- Step-by-step execution flow
- Service call outcomes
- Processing duration
- Success/failure indicators

Example log output:
```
[REQ-1699612345-1234] Received integration request
[REQ-1699612345-1234] Request parsed - UserID: user123, Action: fetch_user_data
[REQ-1699612345-1234] Step 1: Calling Auth Service...
[AUTH] Validated token for user user123: true
[REQ-1699612345-1234] ✓ Auth validated successfully
[REQ-1699612345-1234] Step 2: Calling Database Service...
[DATABASE] Fetched data for user user123, action fetch_user_data
[REQ-1699612345-1234] ✓ Data fetched successfully
[REQ-1699612345-1234] Step 3: Calling Notification Service...
[NOTIFICATION] Sent notification to user user123
[REQ-1699612345-1234] ✓ Notification sent successfully
[REQ-1699612345-1234] Processing complete in 287ms
```

## Features

- **Request Tracking** - Each request gets a unique ID for tracing
- **Simulated Latency** - Services include realistic processing delays
- **Error Handling** - Proper error responses for various failure scenarios
- **Graceful Degradation** - Continues processing even if notification fails
- **Health Monitoring** - Health check endpoint for service status
- **Detailed Logging** - Comprehensive logs for debugging and monitoring

## Authentication

For the demo, authentication logic is simplified:
- Tokens starting with `valid_` are considered valid
- All other tokens will fail authentication

Examples:
- ✅ `valid_token123` - Valid
- ✅ `valid_abc456` - Valid
- ❌ `token123` - Invalid
- ❌ `invalid_token` - Invalid

## Service Timeouts

All internal service calls have a 5-second timeout to prevent hanging requests.

## Use Cases

This demo is useful for:
- Understanding integration patterns
- Testing distributed tracing tools
- Demonstrating service orchestration
- Learning error handling in microservices
- Observability and monitoring demos
- API gateway patterns

