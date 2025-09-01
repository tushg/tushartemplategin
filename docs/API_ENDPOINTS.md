# API Endpoints Documentation

This document describes all available REST API endpoints in the Tushar Template Gin application.

## Base URL
```
http://localhost:8080/api/v1
```

## Health Endpoints

### GET /health
Get overall health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "tushar-service",
  "version": "1.0.0"
}
```

### GET /health/ready
Kubernetes readiness probe endpoint.

**Response:**
```json
{
  "status": "ready",
  "timestamp": "2024-01-01T12:00:00Z",
  "database": "not_required",
  "service": "tushar-service"
}
```

### GET /health/live
Kubernetes liveness probe endpoint.

**Response:**
```json
{
  "status": "alive",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "tushar-service"
}
```

## Product Registration Endpoints

### POST /products
Create a new product.

**Request Body:**
```json
{
  "name": "Sample Product",
  "description": "A sample product description",
  "category": "Electronics",
  "price": 99.99,
  "sku": "SKU-001",
  "stock": 100,
  "is_active": true
}
```

**Response (201 Created):**
```json
{
  "product": {
    "id": 1,
    "name": "Sample Product",
    "description": "A sample product description",
    "category": "Electronics",
    "price": 99.99,
    "sku": "SKU-001",
    "stock": 100,
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### GET /products
List all products with pagination and filtering.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)
- `category` (optional): Filter by category
- `is_active` (optional): Filter by active status (true/false)
- `search` (optional): Search in name, description, or SKU

**Example:**
```
GET /products?page=1&limit=10&category=Electronics&is_active=true&search=laptop
```

**Response (200 OK):**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Sample Product",
      "description": "A sample product description",
      "category": "Electronics",
      "price": 99.99,
      "sku": "SKU-001",
      "stock": 100,
      "is_active": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

### GET /products/:id
Get a specific product by ID.

**Response (200 OK):**
```json
{
  "product": {
    "id": 1,
    "name": "Sample Product",
    "description": "A sample product description",
    "category": "Electronics",
    "price": 99.99,
    "sku": "SKU-001",
    "stock": 100,
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**Response (404 Not Found):**
```json
{
  "error": "Not Found",
  "details": "product with id 999 not found"
}
```

### PUT /products/:id
Update a specific product.

**Request Body (all fields optional):**
```json
{
  "name": "Updated Product Name",
  "description": "Updated description",
  "category": "Updated Category",
  "price": 149.99,
  "sku": "SKU-002",
  "stock": 50,
  "is_active": false
}
```

**Response (200 OK):**
```json
{
  "product": {
    "id": 1,
    "name": "Updated Product Name",
    "description": "Updated description",
    "category": "Updated Category",
    "price": 149.99,
    "sku": "SKU-002",
    "stock": 50,
    "is_active": false,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:30:00Z"
  }
}
```

### DELETE /products/:id
Delete a specific product.

**Response (204 No Content):**
No response body.

**Response (404 Not Found):**
```json
{
  "error": "Not Found",
  "details": "product with id 999 not found"
}
```

### GET /products/sku/:sku
Get a product by SKU.

**Response (200 OK):**
```json
{
  "product": {
    "id": 1,
    "name": "Sample Product",
    "description": "A sample product description",
    "category": "Electronics",
    "price": 99.99,
    "sku": "SKU-001",
    "stock": 100,
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### PATCH /products/:id/stock
Update product stock quantity.

**Request Body:**
```json
{
  "stock": 150
}
```

**Response (200 OK):**
```json
{
  "message": "Stock updated successfully",
  "id": 1,
  "stock": 150
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Bad Request",
  "details": "Invalid request data"
}
```

### 404 Not Found
```json
{
  "error": "Not Found",
  "details": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal Server Error",
  "details": "An unexpected error occurred"
}
```

## Validation Rules

### Product Creation/Update
- `name`: Required, 1-255 characters
- `description`: Optional, max 1000 characters
- `category`: Required, 1-100 characters
- `price`: Required, must be >= 0
- `sku`: Required, 1-50 characters, must be unique
- `stock`: Required, must be >= 0
- `is_active`: Optional boolean

### Pagination
- `page`: Optional, minimum 1
- `limit`: Optional, minimum 1, maximum 100

## Example Usage

### Create a Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming Laptop",
    "description": "High-performance gaming laptop",
    "category": "Electronics",
    "price": 1299.99,
    "sku": "LAPTOP-001",
    "stock": 25,
    "is_active": true
  }'
```

### List Products
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&limit=5&category=Electronics"
```

### Update Product
```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 1199.99,
    "stock": 30
  }'
```

### Delete Product
```bash
curl -X DELETE http://localhost:8080/api/v1/products/1
```

### Update Stock
```bash
curl -X PATCH http://localhost:8080/api/v1/products/1/stock \
  -H "Content-Type: application/json" \
  -d '{
    "stock": 50
  }'
```
