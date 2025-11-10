# Integration Service Demo - cURL Commands
# Copy and paste these commands to test the service

# 1. Health Check
curl -X GET http://localhost:9090/health

# 2. Valid Integration Request (Success)
curl -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "token": "valid_token123",
    "action": "fetch_user_data",
    "metadata": {
      "source": "web",
      "version": "1.0"
    }
  }'

# 3. Invalid Token Request (Auth Failure)
curl -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "token": "invalid_token",
    "action": "fetch_user_data"
  }'

# 4. Another Valid Request with Different Action
curl -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user456",
    "token": "valid_xyz789",
    "action": "update_profile",
    "metadata": {
      "source": "mobile",
      "version": "2.0"
    }
  }'

# 5. Direct Auth Service Call - Valid Token
curl -X POST http://localhost:9090/auth/validate \
  -H "Content-Type: application/json" \
  -d '{
    "token": "valid_test123",
    "user_id": "testuser"
  }'

# 6. Direct Auth Service Call - Invalid Token
curl -X POST http://localhost:9090/auth/validate \
  -H "Content-Type: application/json" \
  -d '{
    "token": "bad_token",
    "user_id": "testuser"
  }'

# 7. Direct Database Service Call
curl -X GET "http://localhost:9090/database/fetch?user_id=testuser&action=test_action"

# 8. Direct Notification Service Call
curl -X POST http://localhost:9090/notification/send \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "testuser",
    "message": "This is a test notification"
  }'

# 9. Request using example JSON file
curl -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d @examples/valid-request.json

# 10. Pretty print JSON output (requires jq)
curl -s -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user999",
    "token": "valid_token999",
    "action": "get_analytics"
  }' | jq '.'

# 11. Show response headers
curl -i -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user777",
    "token": "valid_token777",
    "action": "test_headers"
  }'

# 12. Verbose output for debugging
curl -v -X POST http://localhost:9090/api/process \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "debug_user",
    "token": "valid_debug_token",
    "action": "debug_action"
  }'
