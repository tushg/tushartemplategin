# Tushar Template Gin

A Go web server template with health check endpoints and basic API structure.

## Project Structure

```
tushartemplategin/
├── cmd/
│   └── server/          # Main application entry point
├── configs/             # Configuration files
├── internal/
│   └── health/          # Health check endpoints
├── pkg/
│   ├── config/          # Configuration management
│   ├── logger/          # Logging utilities
│   └── server/          # HTTP server implementation
└── go.mod               # Go module file
```

## Available Endpoints

- `GET /` - Root endpoint with welcome message
- `GET /health` - Basic health check
- `GET /health/detailed` - Detailed health check with system status
- `GET /api/status` - API status information

## Running the Service

### Prerequisites
- Go 1.24.5 or later

### Steps to Run

1. **Navigate to the project directory:**
   ```bash
   cd tushartemplategin
   ```

2. **Run the server:**
   ```bash
   go run cmd/server/main.go
   ```

   Or build and run:
   ```bash
   go build -o server cmd/server/main.go
   ./server
   ```

3. **The server will start on port 8080 by default**

   You can change the port using environment variables:
   ```bash
   PORT=3000 go run cmd/server/main.go
   ```

## Testing the Endpoints

### Using curl

1. **Test root endpoint:**
   ```bash
   curl http://localhost:8080/
   ```

2. **Test health check:**
   ```bash
   curl http://localhost:8080/health
   ```

3. **Test detailed health check:**
   ```bash
   curl http://localhost:8080/health/detailed
   ```

4. **Test API status:**
   ```bash
   curl http://localhost:8080/api/status
   ```

### Using a web browser

Open these URLs in your browser:
- http://localhost:8080/
- http://localhost:8080/health
- http://localhost:8080/health/detailed
- http://localhost:8080/api/status

### Using Postman or similar tools

Import these requests:
- GET http://localhost:8080/
- GET http://localhost:8080/health
- GET http://localhost:8080/health/detailed
- GET http://localhost:8080/api/status

## Environment Variables

- `PORT` - Server port (default: 8080)
- `ENV` - Environment (default: development)
- `LOG_LEVEL` - Log level (default: info)

## Stopping the Service

Press `Ctrl+C` in the terminal where the server is running to gracefully shut it down.

## Expected Responses

### Health Check (`/health`)
```json
{
  "status": "healthy",
  "timestamp": "2025-01-XX...",
  "service": "tushartemplategin",
  "version": "1.0.0"
}
```

### Root Endpoint (`/`)
```json
{
  "message": "Welcome to Tushar Template Gin",
  "service": "tushartemplategin",
  "version": "1.0.0"
}
```

### API Status (`/api/status`)
```json
{
  "status": "running",
  "environment": "development",
  "port": 8080
}
```
