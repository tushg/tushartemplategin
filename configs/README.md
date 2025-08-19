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

## SSL/TLS Configuration

The application supports SSL/TLS for secure HTTPS communication:

### SSL Settings
- **enabled**: Enable/disable SSL/TLS
- **port**: SSL port (default: :443)
- **certFile**: Path to SSL certificate file
- **keyFile**: Path to SSL private key file
- **redirectHTTP**: Redirect HTTP to HTTPS

### Certificate Setup
1. **Place your SSL certificate and private key files** in the `certs/` directory
2. **Update the `certFile` and `keyFile` paths** in your `config.json`
3. **Set `enabled: true`** to activate SSL/TLS
4. **Restart the service** to load the new configuration

### Example SSL Configuration
```json
{
  "server": {
    "ssl": {
      "enabled": true,
      "port": ":443",
      "certFile": "./certs/server.crt",
      "keyFile": "./certs/server.key",
      "redirectHTTP": true
    }
  }
}
```

### Production Certificate Requirements
- **Valid SSL certificate** from a trusted CA (Let's Encrypt, DigiCert, etc.)
- **Private key file** in PEM format
- **Certificate file** in PEM format (including full chain if needed)
- **Proper file permissions** (600 for key, 644 for cert)

## Production Considerations

- Use environment variables for sensitive data (passwords, API keys)
- Keep configuration files in version control (without secrets)
- Use different configuration files for different environments
- Validate configuration before deployment
- **SSL certificates should be stored securely** and not committed to version control
