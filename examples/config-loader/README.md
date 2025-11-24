# Configuration Loader Example

This example demonstrates how to use govalidator to validate application configuration files, ensuring all required settings are present and valid before starting your application.

## Features

- Validates complex, nested configuration structures
- Checks required fields, data types, and value ranges
- Validates enums for constrained options (log levels, database drivers, etc.)
- Provides detailed error messages for misconfigured settings
- Prevents application startup with invalid configuration
- Demonstrates validation of deeply nested objects

## Running the Example

```bash
cd examples/config-loader
go run main.go
```

The program will:
1. Create two configuration files: `config.json` (valid) and `config-invalid.json` (with errors)
2. Load and validate the valid configuration
3. Attempt to load the invalid configuration (will fail with detailed errors)
4. Display results for both test cases

## Configuration Structure

The example validates a complete application configuration with these sections:

### Server Configuration
- Host (hostname or IP pattern validation)
- Port (1-65535)
- Read/Write timeouts (1-300 seconds)
- TLS settings (enabled flag, cert/key files)

### Database Configuration
- Driver (postgres, mysql, or sqlite)
- Connection settings (host, port, credentials)
- Connection pool settings (max connections, idle connections)

### Redis Configuration
- Connection settings (host, port, password, database)
- Database number (0-15, as per Redis limitations)

### Logging Configuration
- Log level (debug, info, warn, error)
- Format (json or text)
- Output destination (stdout, stderr, file)

### Feature Flags
- Boolean flags for various features
- All flags are required (explicit configuration)

## Sample Valid Configuration

```json
{
  "server": {
    "host": "localhost",
    "port": 8080,
    "readTimeout": 30,
    "writeTimeout": 30,
    "tls": {
      "enabled": false,
      "certFile": "",
      "keyFile": ""
    }
  },
  "database": {
    "driver": "postgres",
    "host": "localhost",
    "port": 5432,
    "username": "appuser",
    "password": "secretpassword",
    "database": "appdb",
    "maxConnections": 25,
    "maxIdleConns": 5,
    "connMaxLifetime": 300
  },
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout"
  },
  "features": {
    "enableCache": true,
    "enableRateLimit": true,
    "enableMetrics": true,
    "enableDebugMode": false,
    "maintenanceMode": false
  }
}
```

## Common Validation Errors

The example demonstrates catching these types of errors:

1. **Type Mismatches**: Boolean fields with string values
2. **Out of Range**: Port numbers outside valid range
3. **Invalid Enums**: Database drivers not in allowed list
4. **Empty Required Fields**: Missing usernames or passwords
5. **Missing Sections**: Required configuration sections not present
6. **Pattern Violations**: Invalid hostname patterns

## Key Patterns Demonstrated

1. **Deeply Nested Validation**: Multiple levels of object nesting
2. **Schema Composition**: Building complex schemas from smaller parts
3. **Enum Validation**: Using `OneOfValidator` for constrained values
4. **Range Validation**: Numeric bounds checking
5. **Pattern Matching**: Hostname pattern validation with regex
6. **Required vs Optional**: Mix of required and optional fields
7. **Error Presentation**: User-friendly error messages with `DetailedErrorPresenter`

## Use Cases

This validation pattern is useful for:
- Application configuration files (JSON, parsed YAML, etc.)
- Environment-specific settings validation
- Configuration management systems
- Infrastructure as Code validation
- Deployment configuration verification
- Multi-environment setup validation

## Integration Tips

To integrate this pattern into your application:

1. Define your config schema at package initialization
2. Load and validate config before starting any services
3. Use `DetailedErrorPresenter` for user-friendly error messages
4. Fail fast with clear error messages if config is invalid
5. Consider adding custom validators for business-specific rules

Example integration:

```go
func main() {
    config, err := LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }

    // Proceed with application startup
    startServer(config)
}
```
