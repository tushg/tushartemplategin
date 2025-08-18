# Database Package

This package provides a **production-ready, industry-standard database abstraction layer** with support for multiple database types through a factory pattern.

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Application  │    │   Database       │    │   Database      │
│   Layer        │───▶│   Factory        │───▶│   Implementation│
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Repository   │    │   Configuration  │    │   PostgreSQL    │
│   Layer        │    │   Validation     │    │   SQLite        │
└─────────────────┘    └──────────────────┘    │   MySQL         │
                                              └─────────────────┘
```

## ✨ Features

- **🔌 Multi-Database Support**: PostgreSQL, SQLite, MySQL
- **🏭 Factory Pattern**: Clean database instantiation
- **⚙️ Configuration Driven**: Environment-specific settings
- **🛡️ Validation**: Comprehensive configuration validation
- **📊 Monitoring**: Built-in health checks and metrics
- **🧪 Testable**: Fully testable architecture
- **🚀 Production Ready**: Industry-standard patterns

## 🚀 Quick Start

### 1. Basic Usage

```go
import (
    "tushartemplategin/pkg/database"
    "tushartemplategin/pkg/config"
    "tushartemplategin/pkg/logger"
)

// Create logger
logger := logger.NewLogger(&logger.Config{Level: "info"})

// Create database factory
factory := database.NewDatabaseFactory(logger)

// Create database instance
db, err := factory.CreateDatabase(&cfg.Database)
if err != nil {
    log.Fatalf("Failed to create database: %v", err)
}

// Connect to database
if err := db.Connect(ctx); err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
```

### 2. Configuration

```yaml
# config.yaml
database:
  type: "postgres"  # postgres, sqlite, mysql
  
  postgres:
    host: "localhost"
    port: 5432
    name: "myapp"
    username: "user"
    password: "pass"
    
  sqlite:
    filePath: "./data/app.db"
    
  mysql:
    host: "localhost"
    port: 3306
    name: "myapp"
    username: "user"
    password: "pass"
```

### 3. Environment Variables

```bash
# Override database type
export DB_TYPE=sqlite

# Override specific settings
export DB_POSTGRES_HOST=prod-db.example.com
export DB_POSTGRES_PASSWORD=prod_password
export DB_SQLITE_FILE_PATH=/tmp/test.db
```

## 🗄️ Supported Databases

### PostgreSQL
- **Driver**: `github.com/lib/pq`
- **Features**: Full ACID compliance, advanced indexing, JSON support
- **Use Case**: Production applications, complex queries, high concurrency

### SQLite
- **Driver**: Built-in Go driver
- **Features**: File-based, zero-configuration, embedded
- **Use Case**: Development, testing, embedded applications, single-user

### MySQL
- **Driver**: `github.com/go-sql-driver/mysql`
- **Features**: High performance, replication, clustering
- **Use Case**: Web applications, high-traffic sites, distributed systems

## 🔧 Database Factory

### Core Interface

```go
type DatabaseFactory struct {
    logger logger.Logger
}

func NewDatabaseFactory(log logger.Logger) *DatabaseFactory
func (df *DatabaseFactory) CreateDatabase(cfg *config.DatabaseConfig) (Database, error)
```

### Database Type Constants

```go
const (
    DatabaseTypePostgreSQL DatabaseType = "postgres"
    DatabaseTypeSQLite     DatabaseType = "sqlite"
    DatabaseTypeMySQL      DatabaseType = "mysql"
)
```

### Validation Methods

```go
// Validate database type
func (dt DatabaseType) IsValid() bool

// Check if specific type is supported
func (dt DatabaseType) IsDatabaseTypeSupported(dbType string) bool

// Get database type from config
func (df *DatabaseFactory) GetDatabaseType(cfg *config.DatabaseConfig) DatabaseType
```

## 📋 Configuration Validation

### PostgreSQL Requirements
- ✅ Host (required)
- ✅ Port (1-65535)
- ✅ Database name (required)
- ✅ Username (required)
- ⚠️ Password (optional, generates warning)

### SQLite Requirements
- ✅ File path (required)
- ✅ All other settings have sensible defaults

### MySQL Requirements
- ✅ Host (required)
- ✅ Port (1-65535)
- ✅ Database name (required)
- ✅ Username (required)
- ✅ Password (required)

## 🧪 Testing

### Run Tests

```bash
# Run all database tests
go test ./pkg/database/...

# Run specific test file
go test ./pkg/database/factory_test.go

# Run with verbose output
go test -v ./pkg/database/...

# Run with coverage
go test -cover ./pkg/database/...
```

### Test Coverage

The test suite covers:
- ✅ Database type validation
- ✅ Factory instantiation
- ✅ Configuration validation
- ✅ Error handling
- ✅ Edge cases
- ✅ All public methods

## 🔄 Migration Support

The database factory integrates seamlessly with the migration system:

```go
// Database factory creates the database instance
factory := database.NewDatabaseFactory(logger)
db, err := factory.CreateDatabase(&cfg.Database)

// Migration engine uses the same database instance
migrator := migrations.NewMigrationEngine(db, logger)
err = migrator.Migrate(ctx)
```

## 🚀 Production Deployment

### 1. Environment-Specific Configs

```bash
# Development
DB_TYPE=sqlite
DB_SQLITE_FILE_PATH=./dev.db

# Staging
DB_TYPE=postgres
DB_POSTGRES_HOST=staging-db.example.com

# Production
DB_TYPE=postgres
DB_POSTGRES_HOST=prod-db.example.com
DB_POSTGRES_SSL_MODE=require
```

### 2. Health Checks

```go
// Database health check
if err := db.Health(ctx); err != nil {
    // Handle unhealthy database
    log.Error("Database health check failed", err)
}
```

### 3. Monitoring

```go
// Get database statistics
stats := db.GetConnectionStats()
log.Info("Database stats", stats)
```

## 🔒 Security Considerations

### 1. Credential Management
- Use environment variables for sensitive data
- Never commit passwords to version control
- Use secrets management in production

### 2. Network Security
- Enable SSL/TLS for production databases
- Use firewall rules to restrict access
- Implement proper authentication

### 3. Principle of Least Privilege
- Create database users with minimal permissions
- Use read-only connections where possible
- Audit database access

## 📚 Best Practices

### 1. Configuration Management
```go
// Use configuration validation
if err := validateConfig(cfg); err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}

// Use environment-specific configs
configPath := fmt.Sprintf("configs/config.%s.yaml", os.Getenv("ENV"))
```

### 2. Error Handling
```go
// Always check for errors
db, err := factory.CreateDatabase(cfg)
if err != nil {
    log.Errorf("Database creation failed: %v", err)
    return err
}

// Use proper error wrapping
return fmt.Errorf("failed to connect to %s: %w", dbType, err)
```

### 3. Resource Management
```go
// Always close database connections
defer func() {
    if err := db.Close(); err != nil {
        log.Errorf("Failed to close database: %v", err)
    }
}()

// Use connection pooling appropriately
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
```

## 🔮 Future Enhancements

### Planned Features
- [ ] Connection pooling optimization
- [ ] Database clustering support
- [ ] Advanced monitoring and metrics
- [ ] Automatic failover
- [ ] Database-specific optimizations

### Extension Points
- [ ] Custom database drivers
- [ ] Plugin architecture
- [ ] Custom validation rules
- [ ] Performance profiling

## 📖 Examples

### Complete Application Setup

```go
package main

import (
    "context"
    "log"
    
    "tushartemplategin/pkg/config"
    "tushartemplategin/pkg/database"
    "tushartemplategin/pkg/logger"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Create logger
    appLogger := logger.NewLogger(&cfg.Log)
    
    // Create database factory
    dbFactory := database.NewDatabaseFactory(appLogger)
    
    // Create database instance
    db, err := dbFactory.CreateDatabase(&cfg.Database)
    if err != nil {
        appLogger.Error(context.Background(), "Failed to create database", logger.Fields{
            "error": err.Error(),
            "type":  cfg.Database.Type,
        })
        log.Fatalf("Failed to create database: %v", err)
    }
    
    // Connect to database
    ctx := context.Background()
    if err := db.Connect(ctx); err != nil {
        appLogger.Error(ctx, "Failed to connect to database", logger.Fields{
            "error": err.Error(),
        })
        log.Fatalf("Failed to connect: %v", err)
    }
    
    defer db.Close()
    
    appLogger.Info(ctx, "Successfully connected to database", logger.Fields{
        "type": cfg.Database.Type,
    })
    
    // Your application logic here...
}
```

## 🤝 Contributing

### Adding New Database Support

1. **Create Implementation**
   ```go
   // pkg/database/mysql/connection.go
   type MySQLDB struct {
       config *config.MySQLConfig
       db     *sql.DB
       logger logger.Logger
   }
   ```

2. **Implement Interface**
   ```go
   func (m *MySQLDB) Connect(ctx context.Context) error
   func (m *MySQLDB) Disconnect(ctx context.Context) error
   func (m *MySQLDB) Health(ctx context.Context) error
   // ... other methods
   ```

3. **Update Factory**
   ```go
   case DatabaseTypeMySQL:
       db, err = df.createMySQL(cfg)
   ```

4. **Add Tests**
   ```go
   func TestDatabaseFactory_CreateMySQL(t *testing.T)
   ```

## 📄 License

This package is part of the Tushar Template Gin project and follows the same license terms.

## 🆘 Support

For issues and questions:
1. Check the test files for usage examples
2. Review the configuration examples
3. Check the main application code
4. Create an issue with detailed information

---

**🎯 This database factory pattern follows industry standards and is production-ready for enterprise applications.**
