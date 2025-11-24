# GoValidator Examples

This directory contains real-world, runnable examples demonstrating how to use govalidator in different scenarios.

## Available Examples

### 1. HTTP API Validation (`http-api/`)

Demonstrates validation of REST API requests with:
- User registration endpoint
- Product creation endpoint
- Comprehensive field validation
- Clean error responses for clients

**Use Case**: Building REST APIs with robust input validation

```bash
cd http-api && go run main.go
```

### 2. CSV File Validator (`csv-validator/`)

Shows how to validate CSV file imports with:
- Row-by-row validation
- Type conversion (string to number)
- Detailed error reporting with row numbers
- Export of invalid rows to JSON

**Use Case**: Data import validation, bulk uploads, ETL pipelines

```bash
cd csv-validator && go run main.go
```

### 3. Configuration Loader (`config-loader/`)

Validates application configuration files with:
- Deeply nested object validation
- Database, Redis, and server settings
- Enum validation for configuration options
- Fail-fast configuration loading

**Use Case**: Application startup configuration validation

```bash
cd config-loader && go run main.go
```

### 4. Webhook Handler (`webhook-handler/`)

Handles and validates webhook payloads from external services:
- GitHub push events
- Stripe payment events
- Signature verification
- Complex nested payload validation

**Use Case**: Webhook integrations, event-driven architectures

```bash
cd webhook-handler && go run main.go
```

## Quick Start

Each example is self-contained and can be run independently:

```bash
# Clone the repository
git clone https://github.com/gstachniukrsk/govalidator.git
cd govalidator/examples

# Run any example
cd http-api
go run main.go
```

## Common Patterns Demonstrated

All examples showcase these key patterns:

1. **Modern Schema API**: Using the fluent builder pattern for schema definition
2. **Field-Level Validation**: Combining multiple validators per field
3. **Nested Structures**: Validating objects within objects and arrays
4. **Error Presentation**: Different error formatting strategies
5. **Type Safety**: Converting unstructured JSON to validated data
6. **Extra Fields Control**: Using `ExtraForbid` vs `ExtraIgnore`
7. **Required vs Optional**: Explicit nullability control

## Testing the Examples

You can test the examples using curl or by examining the source code:

```bash
# Terminal 1: Start the HTTP API example
cd http-api && go run main.go

# Terminal 2: Test with curl
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret123","role":"user"}'
```

## Example Selection Guide

Choose the example that best matches your use case:

| Use Case | Example | Key Features |
|----------|---------|--------------|
| REST API validation | `http-api` | Request/response handling |
| Bulk data import | `csv-validator` | File processing, batch validation |
| App configuration | `config-loader` | Nested config, fail-fast loading |
| Third-party webhooks | `webhook-handler` | Signature verification, event handling |

## Learning Path

Recommended order for learning:

1. **http-api**: Start here to learn basic validation patterns
2. **csv-validator**: Learn file processing and batch validation
3. **config-loader**: Understand deeply nested validation
4. **webhook-handler**: Master complex real-world scenarios

## Contributing Examples

Have an idea for a new example? Contributions are welcome!

Good example candidates:
- Form validation (web forms)
- GraphQL input validation
- CLI argument validation
- Database query validation
- Environment variable validation
- Multi-file validation

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## Questions or Issues?

- Check the main [README.md](../README.md) for API documentation
- Review [CLAUDE.md](../CLAUDE.md) for architecture details
- Open an issue on [GitHub](https://github.com/gstachniukrsk/govalidator/issues)

## License

All examples are provided under the same [MIT License](../LICENSE) as the main library.
