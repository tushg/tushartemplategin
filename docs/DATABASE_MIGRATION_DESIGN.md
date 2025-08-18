# Database Migration Design Document
## Using go-migrate (golang-migrate) for PostgreSQL

### Table of Contents
1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Migration Strategy](#migration-strategy)
4. [File Structure](#file-structure)
5. [Implementation Details](#implementation-details)
6. [Deployment Strategy](#deployment-strategy)
7. [Versioning & Rollback](#versioning--rollback)
8. [Security Considerations](#security-considerations)
9. [Testing Strategy](#testing-strategy)
10. [Monitoring & Observability](#monitoring--observability)

---

## Overview

This document outlines the implementation of database migrations using `go-migrate` (golang-migrate) library for the Tushar Template Gin project. The migration system will handle both fresh installations and incremental upgrades while maintaining data integrity and providing rollback capabilities.

### Goals
- **Fresh Install**: Create complete database schema from scratch
- **Incremental Updates**: Apply schema changes without data loss
- **Version Control**: Track migration versions and history
- **Rollback Support**: Ability to revert to previous schema versions
- **Zero Downtime**: Support for blue-green deployments
- **Audit Trail**: Complete history of schema changes

---

## Architecture

### High-Level Architecture
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Application  │    │   Migration      │    │   PostgreSQL    │
│   Startup      │───▶│   Engine         │───▶│   Database      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Config        │    │   Version        │    │   Schema        │
│   Management    │    │   Tracking       │    │   State         │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Components
1. **Migration Engine**: Core migration logic using go-migrate
2. **Version Tracker**: Database table to track applied migrations
3. **Migration Files**: SQL files for up/down migrations
4. **Configuration**: Environment-specific migration settings
5. **Health Checks**: Validation of migration status

---

## Migration Strategy

### 1. Fresh Install Strategy
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Empty DB      │───▶│   Apply All      │───▶│   Full Schema   │
│                 │    │   Migrations     │    │   Ready         │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Process:**
- Check if `schema_migrations` table exists
- If not, create it
- Apply all migrations in sequence (001, 002, 003...)
- Validate final schema state
- Mark installation as complete

### 2. Incremental Update Strategy
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Current       │───▶│   Apply New      │───▶│   Updated       │
│   Schema v1.2   │    │   Migrations     │    │   Schema v1.3   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Process:**
- Query `schema_migrations` table for current version
- Identify pending migrations
- Apply migrations sequentially
- Update version tracking
- Validate schema integrity

### 3. Rollback Strategy
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Current       │───▶│   Rollback       │───▶│   Previous      │
│   Schema v1.3   │    │   Migrations     │    │   Schema v1.2   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Process:**
- Identify target rollback version
- Execute down migrations in reverse order
- Validate rollback success
- Update version tracking

---

## File Structure

```
project-root/
├── migrations/
│   ├── 000001_create_initial_schema.up.sql
│   ├── 000001_create_initial_schema.down.sql
│   ├── 000002_add_health_status_table.up.sql
│   ├── 000002_add_health_status_table.down.sql
│   ├── 000003_add_indexes.up.sql
│   ├── 000003_add_indexes.down.sql
│   └── README.md
├── pkg/
│   └── migrations/
│       ├── engine.go          # Migration engine
│       ├── config.go          # Migration configuration
│       ├── validator.go       # Schema validation
│       └── health.go          # Migration health checks
├── scripts/
│   ├── migrate.sh             # Migration script
│   └── rollback.sh            # Rollback script
└── configs/
    └── migrations.yaml        # Migration configuration
```

### Migration File Naming Convention
- **Format**: `{version}_{description}.{direction}.sql`
- **Version**: 6-digit zero-padded number (000001, 000002, etc.)
- **Description**: Descriptive name with underscores
- **Direction**: `up.sql` (apply) or `down.sql` (rollback)

**Example:**
```
000001_create_initial_schema.up.sql
000001_create_initial_schema.down.sql
000002_add_user_authentication.up.sql
000002_add_user_authentication.down.sql
```

---

## Implementation Details

### 1. Migration Engine Interface
```go
type MigrationEngine interface {
    // Initialize migration system
    Initialize(ctx context.Context) error
    
    // Apply pending migrations
    Migrate(ctx context.Context) error
    
    // Rollback to specific version
    Rollback(ctx context.Context, version int) error
    
    // Get current migration status
    Status(ctx context.Context) (*MigrationStatus, error)
    
    // Validate schema integrity
    Validate(ctx context.Context) error
    
    // Force specific version (for recovery)
    Force(ctx context.Context, version int) error
}
```

### 2. Version Tracking Table
```sql
-- This table is automatically created by go-migrate
CREATE TABLE IF NOT EXISTS schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL DEFAULT false,
    applied_at timestamp with time zone DEFAULT now()
);

-- Index for performance
CREATE INDEX IF NOT EXISTS idx_schema_migrations_version ON schema_migrations(version);
CREATE INDEX IF NOT EXISTS idx_schema_migrations_applied_at ON schema_migrations(applied_at);
```

### 3. Configuration Structure
```yaml
# configs/migrations.yaml
migrations:
  database:
    driver: "postgres"
    host: "${DB_HOST}"
    port: "${DB_PORT}"
    name: "${DB_NAME}"
    username: "${DB_USERNAME}"
    password: "${DB_PASSWORD}"
    sslmode: "${DB_SSL_MODE}"
  
  options:
    migrations_path: "./migrations"
    timeout: 300s
    lock_timeout: 30s
    max_retries: 3
    retry_delay: 5s
  
  validation:
    check_constraints: true
    check_indexes: true
    check_foreign_keys: true
    schema_consistency: true
```

---

## Deployment Strategy

### 1. Application Startup Flow
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   App Start     │───▶│   Check DB       │───▶│   Run           │
│                 │    │   Connection     │    │   Migrations    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Validate      │───▶│   Start HTTP     │───▶│   Ready for     │
│   Schema        │    │   Server         │    │   Traffic       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 2. Blue-Green Deployment
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Blue          │───▶│   Apply          │───▶│   Blue          │
│   Environment   │    │   Migrations     │    │   Updated       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Switch        │───▶│   Green          │───▶│   Production    │
│   Traffic       │    │   Environment    │    │   Traffic       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 3. Rollback Deployment
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Issue         │───▶│   Rollback       │───▶│   Previous      │
│   Detected      │    │   Migrations     │    │   Version       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Validate      │───▶│   Restore        │───▶│   Service       │
│   Rollback      │    │   Traffic        │    │   Restored      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

---

## Versioning & Rollback

### 1. Version Numbering Strategy
```
Version Format: MAJOR.MINOR.PATCH.MIGRATION

Examples:
- 1.0.0.001: Initial schema
- 1.0.0.002: Add health status table
- 1.1.0.003: Add user authentication
- 1.1.1.004: Fix index performance
- 2.0.0.005: Breaking schema changes
```

### 2. Migration Dependencies
```
Migration Dependencies Graph:
001 (Initial Schema) ← 002 (Health Table) ← 003 (Indexes)
     ↓
004 (User Auth) ← 005 (User Permissions)
     ↓
006 (Audit Logs)
```

### 3. Rollback Scenarios
- **Development Rollback**: Rollback to any previous version
- **Staging Rollback**: Rollback to last stable version
- **Production Rollback**: Rollback to last known good version
- **Emergency Rollback**: Immediate rollback to safe version

---

## Security Considerations

### 1. Migration Security
- **Authentication**: Database credentials management
- **Authorization**: Migration execution permissions
- **Audit Logging**: Track all migration executions
- **Encryption**: Secure storage of database credentials

### 2. Data Protection
- **Backup Strategy**: Pre-migration backups
- **Data Validation**: Post-migration data integrity checks
- **Rollback Safety**: Ensure rollback doesn't lose data
- **Testing**: Test migrations in staging environment

### 3. Access Control
- **Principle of Least Privilege**: Minimal database permissions
- **Network Security**: Database connection security
- **Monitoring**: Alert on unauthorized migration attempts

---

## Testing Strategy

### 1. Migration Testing
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Unit Tests    │───▶│   Integration    │───▶│   End-to-End    │
│   (Individual   │    │   Tests          │    │   Tests         │
│   Migrations)   │    │   (Migration     │    │   (Full         │
└─────────────────┘    │   Sequences)     │    │   Workflow)     │
                       └──────────────────┘    └─────────────────┘
```

### 2. Test Environments
- **Local Development**: SQLite for fast testing
- **CI/CD Pipeline**: PostgreSQL test database
- **Staging Environment**: Production-like database
- **Production**: Final validation

### 3. Test Data Management
- **Seed Data**: Realistic test data sets
- **Data Cleanup**: Automatic test data cleanup
- **Data Isolation**: Separate test databases

---

## Monitoring & Observability

### 1. Migration Metrics
```
Metrics to Track:
- Migration execution time
- Success/failure rates
- Rollback frequency
- Schema change impact
- Database performance impact
```

### 2. Health Checks
```
Health Check Endpoints:
- /health/migrations/status
- /health/migrations/version
- /health/migrations/history
- /health/migrations/validation
```

### 3. Alerting
```
Alert Conditions:
- Migration failures
- Long migration execution times
- Schema validation failures
- Rollback events
- Version mismatches
```

---

## Implementation Phases

### Phase 1: Foundation (Week 1)
- [ ] Set up go-migrate dependency
- [ ] Create migration engine interface
- [ ] Implement basic migration functionality
- [ ] Create initial migration files

### Phase 2: Core Features (Week 2)
- [ ] Implement version tracking
- [ ] Add rollback functionality
- [ ] Create configuration management
- [ ] Add basic validation

### Phase 3: Advanced Features (Week 3)
- [ ] Implement health checks
- [ ] Add monitoring and metrics
- [ ] Create deployment scripts
- [ ] Add comprehensive testing

### Phase 4: Production Ready (Week 4)
- [ ] Security hardening
- [ ] Performance optimization
- [ ] Documentation completion
- [ ] Production deployment

---

## Conclusion

This migration system design provides a robust, scalable, and industry-standard approach to database schema management. By implementing go-migrate with proper versioning, rollback capabilities, and comprehensive monitoring, we ensure data integrity and operational reliability.

The system is designed to handle both fresh installations and incremental updates while maintaining backward compatibility and providing emergency rollback capabilities. The modular architecture allows for easy extension and maintenance as the project evolves.

---

## References

- [go-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Migration Best Practices](https://www.postgresql.org/docs/current/ddl.html)
- [Database Migration Patterns](https://martinfowler.com/articles/evodb.html)
- [Zero-Downtime Deployment Strategies](https://martinfowler.com/bliki/BlueGreenDeployment.html)
