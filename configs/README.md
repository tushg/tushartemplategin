# Configuration Guide

This directory contains configuration files for the Tushar Template Gin application.

## Configuration Formats

The application supports two configuration formats:

### 1. JSON Configuration (Default)
- **File**: `config.json`
- **Example**: `config.example.json`
- **Advantages**: 
  - Machine-readable and parseable
  - Widely supported by tools and editors
  - No indentation issues
  - Better for automation and CI/CD

### 2. YAML Configuration (Alternative)
- **File**: `config.yaml`
- **Example**: `config.example.yaml`
- **Advantages**:
  - Human-readable with comments
  - Less verbose than JSON
  - Support for complex data structures

## Switching Between Formats

To switch from YAML to JSON:

1. **Update the configuration loading code** in `pkg/config/config.go`:
   ```go
   // For JSON
   viper.SetConfigType("json")
   
   // For YAML
   viper.SetConfigType("yaml")
   ```

2. **Rename your configuration file**:
   - For JSON: `config.json`
   - For YAML: `config.yaml`

3. **Update the example file** accordingly

## Configuration Structure

Both formats support the same configuration structure:

```json
{
  "server": {
    "port": ":8080",
    "mode": "debug"
  },
  "log": {
    "level": "info",
    "format": "json"
  },
  "database": {
    "type": "postgres",
    "postgres": {
      "host": "localhost",
      "port": 5432,
      "name": "tushar_db"
    }
  }
}
```

## Environment Variable Overrides

Environment variables can override configuration values:

```bash
# Database type
export DB_TYPE=sqlite

# PostgreSQL settings
export DB_POSTGRES_HOST=prod-db.example.com
export DB_POSTGRES_PASSWORD=prod_password

# SQLite settings
export DB_SQLITE_FILE_PATH=/tmp/test.db
```

## File Locations

The application searches for configuration files in this order:
1. `./configs/` directory
2. Current working directory (`.`)

## Validation

Configuration is automatically validated on startup:
- Required fields are checked
- Database type validation
- Type-specific configuration validation
- Environment variable overrides are applied

## Production Considerations

- Use environment variables for sensitive data (passwords, API keys)
- Keep configuration files in version control (without secrets)
- Use different configuration files for different environments
- Validate configuration before deployment
