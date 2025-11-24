# HTTP API Validation Example

This example demonstrates how to use govalidator in a REST API to validate incoming HTTP requests.

## Features

- User registration endpoint with comprehensive validation
- Product creation endpoint with nested array validation
- Clean error responses with detailed validation messages
- Demonstrates ExtraForbid mode to reject unexpected fields
- Pattern matching for email validation
- Range validation for numeric values

## Running the Example

```bash
cd examples/http-api
go run main.go
```

The server will start on `http://localhost:8080`.

## API Endpoints

### POST /api/users/register

Register a new user with validation.

**Valid Request:**
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secret123",
    "role": "user",
    "age": 25
  }'
```

**Invalid Request (validation errors):**
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jo",
    "email": "invalid-email",
    "password": "short"
  }'
```

### POST /api/products

Create a new product with validation.

**Valid Request:**
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "categories": ["electronics", "computers"],
    "tags": ["new", "featured"]
  }'
```

**Invalid Request (negative price):**
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": -100,
    "categories": []
  }'
```

## Validation Rules

### User Registration
- **username**: Required, 3-20 characters
- **email**: Required, valid email format
- **password**: Required, minimum 8 characters
- **age**: Optional, between 13-120
- **role**: Required, must be "admin", "user", or "guest"
- Extra fields are forbidden

### Product Creation
- **name**: Required, 3-100 characters
- **description**: Optional, max 1000 characters
- **price**: Required, minimum 0.01
- **categories**: Required array of non-empty strings
- **tags**: Optional array of strings
- Extra fields are forbidden

## Key Patterns Demonstrated

1. **Schema Definition**: Using the modern Schema API with builder pattern
2. **Flat Error Lists**: Using `ValidateFlat()` for API-friendly error responses
3. **Nested Validation**: Validating arrays and objects
4. **Custom Validators**: Using RegexpValidator for email validation
5. **Range Validation**: Using Min/Max validators for numeric fields
6. **Extra Fields Control**: Using `ExtraForbid` to reject unexpected fields
