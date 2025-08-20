# 🐳 **Docker Containerization Design Specification**

> **📋 Simplified Approach: Pure config.json Only with Volume Mounting for SSL Certificates**

**Version**: 1.0  
**Date**: January 19, 2025  
**Status**: Design Complete  
**Base Image**: SUSE Linux with Go Runtime  
**Architecture**: Pure config.json Configuration + Volume Mounting

---

## 📋 **Table of Contents**

1. [Design Overview](#design-overview)
2. [Architecture Principles](#architecture-principles)
3. [Base Image Strategy](#base-image-strategy)
4. [Configuration Management](#configuration-management)
5. [SSL Certificate Handling](#ssl-certificate-handling)
6. [Container Structure](#container-structure)
7. [Deployment Strategies](#deployment-strategies)
8. [Security Considerations](#security-considerations)
9. [Implementation Details](#implementation-details)
10. [Testing Strategy](#testing-strategy)
11. [Production Deployment](#production-deployment)

---

## 🎯 **Design Overview**

### **Design Goals**
- **Universal Docker image** that runs anywhere
- **Pure config.json configuration** - no environment variables
- **Volume mounting for SSL certificates** and external configuration
- **SUSE Linux base image** with pre-installed Go runtime
- **Production-ready** with security best practices

### **Key Principles**
1. **Simplicity**: Single configuration source (config.json)
2. **Flexibility**: Volume mounting for all external resources
3. **Security**: Non-root user, minimal attack surface
4. **Portability**: Same image works across all environments
5. **Maintainability**: Clear separation of concerns

---

## 🏗️ **Architecture Principles**

### **1. Configuration-First Design**
```
┌─────────────────────────────────────┐
│         Configuration Priority       │
├─────────────────────────────────────┤
│ Priority 1: Mounted config.json    │
│         (External volume)           │
├─────────────────────────────────────┤
│ Priority 2: Default config.json    │
│         (Built into image)          │
└─────────────────────────────────────┘
```

### **2. Volume-Centric Architecture**
```
┌─────────────────────────────────────┐
│         Container Image             │
│  ┌─────────────────────────────┐   │
│  │      Go Binary              │   │
│  │   + Default config.json     │   │
│  │   + Application Logic       │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────┐
│      Volume Mounts                  │
│  ┌─────────────┐ ┌─────────────┐   │
│  │ config.json │ │ SSL Certs   │   │
│  │ (External)  │ │ (External)  │   │
│  └─────────────┘ └─────────────┘   │
└─────────────────────────────────────┘
```

### **3. Stateless Application Design**
- **No persistent data** stored in container
- **All configuration** from external sources
- **All certificates** from volume mounts
- **All data** stored in external databases

---

## 🐧 **Base Image Strategy**

### **SUSE Linux Base Image Selection**

#### **Option 1: SUSE Linux Enterprise Server (SLES)**
```dockerfile
FROM registry.suse.com/suse/sle15:latest
```

**Benefits:**
- ✅ **Enterprise-grade** security and support
- ✅ **Long-term support** (LTS)
- ✅ **Go runtime** pre-installed
- ✅ **Security patches** and updates

**Considerations:**
- ❌ **Larger image size** (~200-300MB)
- ❌ **License requirements** for production use

#### **Option 2: OpenSUSE Tumbleweed**
```dockerfile
FROM opensuse/tumbleweed:latest
```

**Benefits:**
- ✅ **Rolling release** with latest packages
- ✅ **Smaller image size** (~150-200MB)
- ✅ **Go runtime** available
- ✅ **No license restrictions**

**Considerations:**
- ❌ **Frequent updates** may introduce instability
- ❌ **Less predictable** for production

#### **Option 3: SUSE Linux Enterprise Micro**
```dockerfile
FROM registry.suse.com/suse/sle-micro:latest
```

**Benefits:**
- ✅ **Minimal attack surface** (~30-50MB)
- ✅ **Security-focused** design
- ✅ **Go runtime** available
- ✅ **Enterprise support**

**Considerations:**
- ❌ **Limited package availability**
- ❌ **More complex** setup

### **Recommended Base Image**
```dockerfile
# Use SUSE Linux Enterprise Server for production stability
FROM registry.suse.com/suse/sle15:latest
```

**Rationale:**
- **Production stability** with long-term support
- **Pre-installed Go runtime** reduces build complexity
- **Security patches** and enterprise support
- **Predictable behavior** across deployments

---

## ⚙️ **Configuration Management**

### **Configuration File Structure**

#### **Default config.json (Built into Image)**
```json
{
  "server": {
    "port": ":8080",
    "mode": "debug",
    "ssl": {
      "enabled": false,
      "port": ":443",
      "certFile": "/app/certs/server.crt",
      "keyFile": "/app/certs/server.key",
      "redirectHTTP": false
    }
  },
  "database": {
    "type": "sqlite",
    "sqlite": {
      "filePath": "./data/app.db"
    }
  },
  "log": {
    "level": "info",
    "format": "json",
    "output": "stdout",
    "filePath": "",
    "maxSize": 100,
    "maxBackup": 3,
    "maxAge": 28,
    "compress": true,
    "addCaller": true,
    "addStack": false
  }
}
```

#### **External Configuration Override**
```json
{
  "server": {
    "mode": "production",
    "ssl": {
      "enabled": true,
      "redirectHTTP": true
    }
  },
  "database": {
    "type": "postgres",
    "postgres": {
      "host": "postgres-service",
      "port": 5432,
      "name": "app_db",
      "username": "app_user",
      "password": "app_password",
      "sslMode": "require"
    }
  },
  "log": {
    "level": "info",
    "format": "json"
  }
}
```

### **Configuration Loading Strategy**

```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// Configuration loading with priority system
func LoadConfig() (*Config, error) {
    var config *Config
    
    // Priority 1: Try to load external config.json (mounted volume)
    if externalConfig, err := loadExternalConfig(); err == nil {
        config = externalConfig
        log.Info("Loaded external configuration from mounted volume")
    } else {
        // Priority 2: Fall back to default config.json (built into image)
        if defaultConfig, err := loadDefaultConfig(); err == nil {
            config = defaultConfig
            log.Info("Loaded default configuration from image")
        } else {
            return nil, fmt.Errorf("failed to load any configuration: %v", err)
        }
    }
    
    // Validate configuration
    if err := validateConfig(config); err != nil {
        return nil, fmt.Errorf("configuration validation failed: %v", err)
    }
    
    return config, nil
}

func loadExternalConfig() (*Config, error) {
    // Try multiple possible mount points
    configPaths := []string{
        "/app/config/config.json",           // Primary mount point
        "/config/config.json",               // Alternative mount point
        "./config.json",                     // Current directory
    }
    
    for _, path := range configPaths {
        if config, err := loadConfigFile(path); err == nil {
            return config, nil
        }
    }
    
    return nil, fmt.Errorf("no external configuration found")
}

func loadDefaultConfig() (*Config, error) {
    return loadConfigFile("/app/config/default.json")
}

func loadConfigFile(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

---

## 🔒 **SSL Certificate Handling**

### **Certificate Volume Mounting Strategy**

#### **Directory Structure**
```
/app/certs/
├── server.crt    # SSL certificate (644 permissions)
└── server.key    # Private key (600 permissions)
```

#### **Volume Mount Points**
```bash
# Primary mount point
-v /host/path/certs:/app/certs

# Alternative mount points
-v /host/path/certs:/etc/ssl/certs
-v /host/path/certs:/config/certs
```

### **SSL Configuration in config.json**
```json
{
  "server": {
    "ssl": {
      "enabled": true,
      "port": ":443",
      "certFile": "/app/certs/server.crt",
      "keyFile": "/app/certs/server.key",
      "redirectHTTP": true
    }
  }
}
```

### **Certificate File Permissions**
```bash
# Secure certificate files
chmod 600 /app/certs/server.key    # Private key - owner read/write only
chmod 644 /app/certs/server.crt    # Certificate - owner read/write, others read
chown appuser:appgroup /app/certs/*
```

---

## 🏗️ **Container Structure**

### **Dockerfile Design**

```dockerfile
# Use SUSE Linux Enterprise Server with Go runtime
FROM registry.suse.com/suse/sle15:latest

# Set working directory
WORKDIR /app

# Install required packages
RUN zypper --non-interactive install \
    ca-certificates \
    && zypper clean --all

# Create application user and group
RUN groupadd -g 1001 appgroup && \
    useradd -u 1001 -g appgroup -m -s /bin/bash appuser

# Create necessary directories
RUN mkdir -p /app/config /app/certs /app/data /app/logs

# Copy application binary
COPY --chown=appuser:appgroup ./main /app/

# Copy default configuration
COPY --chown=appuser:appgroup ./configs/config.json /app/config/default.json

# Set proper permissions
RUN chmod +x /app/main && \
    chmod 644 /app/config/default.json

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 443

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Default command
CMD ["/app/main"]
```

### **Container File System Layout**
```
/app/
├── main                    # Application binary
├── config/
│   └── default.json       # Default configuration
├── certs/                  # SSL certificates (mounted)
├── data/                   # Application data (mounted)
└── logs/                   # Application logs (mounted)
```

### **Volume Mount Points**
```bash
# Configuration
-v /host/config.json:/app/config/config.json

# SSL Certificates
-v /host/certs:/app/certs

# Data persistence
-v /host/data:/app/data

# Log files
-v /host/logs:/app/logs
```

---

## 🚀 **Deployment Strategies**

### **1. Local Development**

#### **Docker Run**
```bash
# Simple run with default config
docker run -p 8080:8080 your-app:latest

# With external configuration
docker run -p 8080:8080 \
    -v $(pwd)/config.json:/app/config/config.json \
    your-app:latest

# With SSL certificates
docker run -p 8080:8080 -p 443:443 \
    -v $(pwd)/config.json:/app/config/config.json \
    -v $(pwd)/certs:/app/certs \
    your-app:latest
```

#### **Docker Compose**
```yaml
version: '3.8'

services:
  app:
    image: your-app:latest
    ports:
      - "8080:8080"
      - "443:443"
    volumes:
      # External configuration
      - ./config.json:/app/config/config.json
      # SSL certificates
      - ./certs:/app/certs
      # Data persistence
      - ./data:/app/data
      # Log files
      - ./logs:/app/logs
    restart: unless-stopped

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=app_db
      - POSTGRES_USER=app_user
      - POSTGRES_PASSWORD=app_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:
```

### **2. Staging Environment**

#### **Docker Run with Staging Config**
```bash
docker run -d -p 8080:8080 -p 443:443 \
    -v /etc/app/config-staging.json:/app/config/config.json \
    -v /etc/ssl/certs-staging:/app/certs \
    -v /var/app/data:/app/data \
    -v /var/app/logs:/app/logs \
    --name app-staging \
    your-app:latest
```

### **3. Production Environment**

#### **Docker Run with Production Config**
```bash
docker run -d -p 8080:8080 -p 443:443 \
    -v /etc/app/config-prod.json:/app/config/config.json \
    -v /etc/ssl/certs:/app/certs \
    -v /var/app/data:/app/data \
    -v /var/app/logs:/app/logs \
    --name app-production \
    --restart unless-stopped \
    your-app:latest
```

---

## 🔒 **Security Considerations**

### **1. Container Security**
- **Non-root user**: Application runs as `appuser` (UID 1001)
- **Minimal packages**: Only essential packages installed
- **Read-only mounts**: SSL certificates mounted as read-only
- **No secrets**: No sensitive data in image layers

### **2. File Permissions**
```bash
# Application binary
chmod 755 /app/main

# Configuration files
chmod 644 /app/config/default.json

# SSL certificates
chmod 600 /app/certs/server.key    # Private key
chmod 644 /app/certs/server.crt    # Certificate

# Data directories
chmod 755 /app/data /app/logs
```

### **3. Network Security**
- **Port exposure**: Only necessary ports (8080, 443)
- **Health checks**: Built-in health monitoring
- **SSL/TLS**: Proper certificate handling

---

## 🔧 **Implementation Details**

### **1. Build Process**

#### **Local Build**
```bash
# Build image locally
docker build -t your-app:latest .

# Test locally
docker run -p 8080:8080 your-app:latest
```

#### **CI/CD Build**
```bash
# Build in CI/CD pipeline
docker build -t your-app:$BUILD_NUMBER .
docker tag your-app:$BUILD_NUMBER your-app:latest

# Push to registry
docker push your-app:$BUILD_NUMBER
docker push your-app:latest
```

### **2. Configuration Management**

#### **Environment-Specific Configs**
```bash
# Development
cp configs/config.json configs/config-dev.json

# Staging
cp configs/config.json configs/config-staging.json

# Production
cp configs/config.json configs/config-prod.json
```

#### **Configuration Validation**
```go
func validateConfig(config *Config) error {
    // Validate server configuration
    if config.Server.Port == "" {
        return fmt.Errorf("server port is required")
    }
    
    // Validate SSL configuration
    if config.Server.SSL.Enabled {
        if config.Server.SSL.CertFile == "" {
            return fmt.Errorf("SSL certificate file is required when SSL is enabled")
        }
        if config.Server.SSL.KeyFile == "" {
            return fmt.Errorf("SSL key file is required when SSL is enabled")
        }
        
        // Check if certificate files exist
        if _, err := os.Stat(config.Server.SSL.CertFile); err != nil {
            return fmt.Errorf("SSL certificate file not found: %s", config.Server.SSL.CertFile)
        }
        if _, err := os.Stat(config.Server.SSL.KeyFile); err != nil {
            return fmt.Errorf("SSL key file not found: %s", config.Server.SSL.KeyFile)
        }
    }
    
    // Validate database configuration
    if err := validateDatabaseConfig(config.Database); err != nil {
        return err
    }
    
    return nil
}
```

### **3. SSL Certificate Management**

#### **Certificate Generation (Development)**
```bash
# Generate self-signed certificate for development
openssl req -x509 -newkey rsa:4096 \
    -keyout server.key \
    -out server.crt \
    -days 365 \
    -nodes \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
```

#### **Certificate Deployment**
```bash
# Copy certificates to host
sudo mkdir -p /etc/ssl/certs
sudo cp server.crt /etc/ssl/certs/
sudo cp server.key /etc/ssl/certs/

# Set proper permissions
sudo chmod 600 /etc/ssl/certs/server.key
sudo chmod 644 /etc/ssl/certs/server.crt
sudo chown root:root /etc/ssl/certs/*
```

---

## 🧪 **Testing Strategy**

### **1. Unit Testing**
```go
func TestConfigLoading(t *testing.T) {
    // Test default configuration loading
    config, err := LoadConfig()
    assert.NoError(t, err)
    assert.NotNil(t, config)
    
    // Test configuration validation
    err = validateConfig(config)
    assert.NoError(t, err)
}
```

### **2. Integration Testing**
```bash
# Test with default configuration
docker run --rm your-app:latest

# Test with external configuration
docker run --rm \
    -v $(pwd)/test-config.json:/app/config/config.json \
    your-app:latest

# Test with SSL certificates
docker run --rm \
    -v $(pwd)/test-config.json:/app/config/config.json \
    -v $(pwd)/test-certs:/app/certs \
    your-app:latest
```

### **3. End-to-End Testing**
```bash
# Start application
docker run -d --name test-app \
    -v $(pwd)/test-config.json:/app/config/config.json \
    -v $(pwd)/test-certs:/app/certs \
    -p 8080:8080 -p 443:443 \
    your-app:latest

# Test HTTP endpoint
curl http://localhost:8080/api/v1/health

# Test HTTPS endpoint
curl -k https://localhost:443/api/v1/health

# Cleanup
docker stop test-app
docker rm test-app
```

---

## 🚀 **Production Deployment**

### **1. Production Configuration**

#### **Production config.json**
```json
{
  "server": {
    "port": ":8080",
    "mode": "release",
    "ssl": {
      "enabled": true,
      "port": ":443",
      "certFile": "/app/certs/server.crt",
      "keyFile": "/app/certs/server.key",
      "redirectHTTP": true
    }
  },
  "database": {
    "type": "postgres",
    "postgres": {
      "host": "postgres-service",
      "port": 5432,
      "name": "app_db",
      "username": "app_user",
      "password": "app_password",
      "sslMode": "require",
      "maxOpenConns": 25,
      "maxIdleConns": 5,
      "connMaxLifetime": "5m",
      "connMaxIdleTime": "1m",
      "timeout": "30s",
      "maxRetries": 3,
      "retryDelay": "1s",
      "healthCheckInterval": "30s"
    }
  },
  "log": {
    "level": "info",
    "format": "json",
    "output": "stdout",
    "maxSize": 100,
    "maxBackup": 3,
    "maxAge": 28,
    "compress": true,
    "addCaller": true,
    "addStack": false
  }
}
```

### **2. Production Deployment Script**

#### **deploy.sh**
```bash
#!/bin/bash

# Production deployment script
set -e

# Configuration
APP_NAME="your-app"
APP_VERSION="latest"
CONFIG_PATH="/etc/app/config-prod.json"
CERTS_PATH="/etc/ssl/certs"
DATA_PATH="/var/app/data"
LOGS_PATH="/var/app/logs"

# Stop existing container
echo "Stopping existing container..."
docker stop $APP_NAME || true
docker rm $APP_NAME || true

# Pull latest image
echo "Pulling latest image..."
docker pull $APP_NAME:$APP_VERSION

# Start new container
echo "Starting new container..."
docker run -d \
    --name $APP_NAME \
    --restart unless-stopped \
    -p 8080:8080 \
    -p 443:443 \
    -v $CONFIG_PATH:/app/config/config.json:ro \
    -v $CERTS_PATH:/app/certs:ro \
    -v $DATA_PATH:/app/data \
    -v $LOGS_PATH:/app/logs \
    $APP_NAME:$APP_VERSION

# Wait for health check
echo "Waiting for health check..."
sleep 10

# Verify deployment
if curl -f http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "Deployment successful!"
else
    echo "Deployment failed!"
    docker logs $APP_NAME
    exit 1
fi
```

### **3. Monitoring and Logging**

#### **Health Check Endpoint**
```go
// Health check for container orchestration
func healthCheckHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "timestamp": time.Now(),
        "service": "your-app",
        "version": "1.0.0",
    })
}
```

#### **Logging Configuration**
```go
// Structured logging for production
func setupLogging(config *LogConfig) {
    // Configure logging based on config.json
    // Output to stdout for container logs
    // Structured JSON format for log aggregation
}
```

---

## 📋 **Summary**

### **Design Benefits**
1. ✅ **Universal compatibility**: Same image works everywhere
2. ✅ **Simple configuration**: Pure config.json approach
3. ✅ **Secure**: Volume mounting for SSL certificates
4. ✅ **Production ready**: SUSE Linux base with Go runtime
5. ✅ **Easy deployment**: Volume mounting for all external resources

### **Key Features**
- **SUSE Linux base image** with pre-installed Go runtime
- **Pure config.json configuration** - no environment variables
- **Volume mounting** for SSL certificates and external configuration
- **Non-root user** execution for security
- **Health checks** for container orchestration
- **Multi-environment support** with configuration files

### **Deployment Scenarios**
1. **Local Development**: Volume mounting with local config
2. **Staging**: Volume mounting with staging config
3. **Production**: Volume mounting with production config
4. **CI/CD**: Automated builds and deployments

This design provides a **clean, simple, and secure** approach to containerizing your microservice while maintaining **maximum flexibility** across different deployment environments.

---

## ❓ **Next Steps**

1. **Review and approve** this design specification
2. **Implement Dockerfile** based on SUSE Linux base
3. **Create configuration templates** for different environments
4. **Set up CI/CD pipeline** for automated builds
5. **Test deployment** across different scenarios
6. **Deploy to production** with monitoring

Would you like me to proceed with implementing this design?
