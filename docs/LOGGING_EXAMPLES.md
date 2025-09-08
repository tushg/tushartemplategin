# 📋 **Logging System Examples and Use Cases**

> **Comprehensive examples for timestamp-based logging system**

---

## 🎯 **Quick Start Examples**

### **Basic Logger Setup**

```go
package main

import (
    "context"
    "log"
    "tushartemplategin/pkg/logger"
)

func main() {
    // Simple configuration
    config := &logger.Config{
        Level:      "info",
        Format:     "json",
        Output:     "file",
        FilePath:   "./logs",
        MaxSize:    20,
        MaxBackups: 3,
        MaxAge:     7,
        Compress:   true,
        AddCaller:  true,
        AddStack:   false,
    }
    
    // Initialize logger
    appLogger, err := logger.NewLogger(config)
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    
    // Use logger
    ctx := context.Background()
    appLogger.Info(ctx, "Application started", map[string]interface{}{
        "version": "1.0.0",
        "build":   "timestamp-logging",
    })
}
```

### **Configuration File Examples**

#### **YAML Configuration**
```yaml
# config.yaml
log:
  level: "info"
  format: "json"
  output: "file"
  filePath: "./logs"
  maxSize: 100
  maxBackup: 3
  maxAge: 0
  compress: true
  addCaller: true
  addStack: false
```

#### **JSON Configuration**
```json
{
  "log": {
    "level": "info",
    "format": "json",
    "output": "file",
    "filePath": "./logs",
    "maxSize": 100,
    "maxBackup": 3,
    "maxAge": 0,
    "compress": true,
    "addCaller": true,
    "addStack": false
  }
}
```

---

## 📁 **File Structure Examples**

### **Generated Log Files**

```
./logs/
├── app_2025-01-19_14-30-25.log          # Current active file
├── app_2025-01-19_14-30-25.log.1.gz     # 1st rotation (compressed)
├── app_2025-01-19_14-30-25.log.2.gz     # 2nd rotation (compressed)
└── app_2025-01-19_14-30-25.log.3.gz     # 3rd rotation (compressed)
```

### **File Naming Convention**

| Pattern | Example | Description |
|---------|---------|-------------|
| `{prefix}_{timestamp}.log` | `app_2025-01-19_14-30-25.log` | Current active file |
| `{prefix}_{timestamp}.log.N.gz` | `app_2025-01-19_14-30-25.log.1.gz` | Rotated file (compressed) |
| `{prefix}_{timestamp}.log.N` | `app_2025-01-19_14-30-25.log.1` | Rotated file (uncompressed) |

---

## 📄 **Log Entry Examples**

### **JSON Format Logs**

#### **Info Level Log**
```json
{
  "level": "info",
  "ts": "2025-01-19T14:30:25.123Z",
  "caller": "main.go:45",
  "msg": "Application started",
  "version": "1.0.0",
  "build": "timestamp-logging",
  "timestamp": "2025-01-19T14:30:25.123Z"
}
```

#### **Error Level Log**
```json
{
  "level": "error",
  "ts": "2025-01-19T14:30:25.123Z",
  "caller": "handler.go:78",
  "msg": "Database connection failed",
  "error": "connection timeout",
  "database": "postgres",
  "host": "localhost:5432",
  "retry_count": 3,
  "timestamp": "2025-01-19T14:30:25.123Z"
}
```

#### **Debug Level Log**
```json
{
  "level": "debug",
  "ts": "2025-01-19T14:30:25.123Z",
  "caller": "service.go:123",
  "msg": "Processing user request",
  "user_id": "12345",
  "request_id": "req-abc-123",
  "method": "GET",
  "path": "/api/users",
  "timestamp": "2025-01-19T14:30:25.123Z"
}
```

### **Console Format Logs**

#### **Info Level Log**
```
2025-01-19T14:30:25.123Z	INFO	main.go:45	Application started	{"version": "1.0.0", "build": "timestamp-logging"}
```

#### **Error Level Log**
```
2025-01-19T14:30:25.123Z	ERROR	handler.go:78	Database connection failed	{"error": "connection timeout", "database": "postgres", "host": "localhost:5432"}
```

---

## 🔄 **Rotation Examples**

### **Size-Based Rotation Timeline**

#### **Configuration**
```yaml
maxSize: 20      # 20 MB per file
maxBackup: 3     # Keep 3 backups
maxAge: 0        # No age limit
```

#### **Timeline**
```
Time 0:00  - app_2025-01-19_14-30-25.log (0 MB)     # Logger starts
Time 0:10  - app_2025-01-19_14-30-25.log (20 MB)    # 1st rotation
Time 0:10  - app_2025-01-19_14-30-25.log.1.gz       # Previous file compressed
Time 0:10  - app_2025-01-19_14-30-25.log (0 MB)     # New file starts
Time 0:20  - app_2025-01-19_14-30-25.log (20 MB)    # 2nd rotation
Time 0:20  - app_2025-01-19_14-30-25.log.1.gz       # Previous backup
Time 0:20  - app_2025-01-19_14-30-25.log.2.gz       # New backup
Time 0:20  - app_2025-01-19_14-30-25.log (0 MB)     # New file starts
Time 0:30  - app_2025-01-19_14-30-25.log (20 MB)    # 3rd rotation
Time 0:30  - app_2025-01-19_14-30-25.log.1.gz       # Previous backup
Time 0:30  - app_2025-01-19_14-30-25.log.2.gz       # Previous backup
Time 0:30  - app_2025-01-19_14-30-25.log.3.gz       # New backup
Time 0:30  - app_2025-01-19_14-30-25.log (0 MB)     # New file starts
Time 0:40  - app_2025-01-19_14-30-25.log (20 MB)    # 4th rotation (maxBackup limit)
Time 0:40  - app_2025-01-19_14-30-25.log.1.gz       # Previous backup
Time 0:40  - app_2025-01-19_14-30-25.log.2.gz       # Previous backup
Time 0:40  - app_2025-01-19_14-30-25.log.3.gz       # Previous backup
Time 0:40  - app_2025-01-19_14-30-25.log (0 MB)     # New file starts
# .4.gz gets DELETED (exceeds maxBackup: 3)
```

### **File Evolution Example**

#### **After 1st Rotation**
```
./logs/
├── app_2025-01-19_14-30-25.log          # Current (0 MB)
└── app_2025-01-19_14-30-25.log.1.gz     # 1st backup (20 MB)
```

#### **After 2nd Rotation**
```
./logs/
├── app_2025-01-19_14-30-25.log          # Current (0 MB)
├── app_2025-01-19_14-30-25.log.1.gz     # 1st backup (20 MB)
└── app_2025-01-19_14-30-25.log.2.gz     # 2nd backup (20 MB)
```

#### **After 3rd Rotation**
```
./logs/
├── app_2025-01-19_14-30-25.log          # Current (0 MB)
├── app_2025-01-19_14-30-25.log.1.gz     # 1st backup (20 MB)
├── app_2025-01-19_14-30-25.log.2.gz     # 2nd backup (20 MB)
└── app_2025-01-19_14-30-25.log.3.gz     # 3rd backup (20 MB)
```

#### **After 4th Rotation (maxBackup limit reached)**
```
./logs/
├── app_2025-01-19_14-30-25.log          # Current (0 MB)
├── app_2025-01-19_14-30-25.log.1.gz     # 1st backup (20 MB)
├── app_2025-01-19_14-30-25.log.2.gz     # 2nd backup (20 MB)
└── app_2025-01-19_14-30-25.log.3.gz     # 3rd backup (20 MB)
# 4th backup gets DELETED (exceeds maxBackup: 3)
```

---

## 🎯 **Use Case Examples**

### **Web Application Logging**

```go
package main

import (
    "context"
    "net/http"
    "tushartemplategin/pkg/logger"
)

func main() {
    // Initialize logger
    config := &logger.Config{
        Level:      "info",
        Format:     "json",
        Output:     "file",
        FilePath:   "./logs",
        MaxSize:    100,
        MaxBackups: 5,
        MaxAge:     7,
        Compress:   true,
        AddCaller:  true,
        AddStack:   false,
    }
    
    appLogger, _ := logger.NewLogger(config)
    
    // HTTP handler with logging
    http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        ctx := context.Background()
        
        // Log request
        appLogger.Info(ctx, "Request received", map[string]interface{}{
            "method": r.Method,
            "path":   r.URL.Path,
            "ip":     r.RemoteAddr,
            "user_agent": r.UserAgent(),
        })
        
        // Process request
        // ...
        
        // Log response
        appLogger.Info(ctx, "Request completed", map[string]interface{}{
            "method": r.Method,
            "path":   r.URL.Path,
            "status": 200,
            "duration_ms": 150,
        })
    })
    
    appLogger.Info(context.Background(), "Server started", map[string]interface{}{
        "port": ":8080",
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### **Database Operation Logging**

```go
package main

import (
    "context"
    "database/sql"
    "tushartemplategin/pkg/logger"
)

func main() {
    appLogger, _ := logger.NewLogger(&logger.Config{
        Level:      "info",
        Format:     "json",
        Output:     "file",
        FilePath:   "./logs",
        MaxSize:    50,
        MaxBackups: 3,
        MaxAge:     14,
        Compress:   true,
        AddCaller:  true,
        AddStack:   false,
    })
    
    // Database operation with logging
    func createUser(ctx context.Context, db *sql.DB, user User) error {
        appLogger.Info(ctx, "Creating user", map[string]interface{}{
            "user_id": user.ID,
            "email":   user.Email,
        })
        
        query := "INSERT INTO users (id, email, name) VALUES (?, ?, ?)"
        result, err := db.ExecContext(ctx, query, user.ID, user.Email, user.Name)
        if err != nil {
            appLogger.Error(ctx, "Failed to create user", map[string]interface{}{
                "user_id": user.ID,
                "error":   err.Error(),
            })
            return err
        }
        
        id, _ := result.LastInsertId()
        appLogger.Info(ctx, "User created successfully", map[string]interface{}{
            "user_id": user.ID,
            "inserted_id": id,
        })
        
        return nil
    }
}
```

### **Error Handling and Recovery**

```go
package main

import (
    "context"
    "tushartemplategin/pkg/logger"
)

func main() {
    appLogger, _ := logger.NewLogger(&logger.Config{
        Level:      "info",
        Format:     "json",
        Output:     "file",
        FilePath:   "./logs",
        MaxSize:    100,
        MaxBackups: 3,
        MaxAge:     0,
        Compress:   true,
        AddCaller:  true,
        AddStack:   true,  // Enable stack traces for errors
    })
    
    // Error handling with logging
    func processData(ctx context.Context, data []byte) error {
        appLogger.Info(ctx, "Processing data", map[string]interface{}{
            "data_size": len(data),
        })
        
        // Simulate processing
        if len(data) == 0 {
            err := errors.New("empty data")
            appLogger.Error(ctx, "Data processing failed", map[string]interface{}{
                "error": err.Error(),
                "data_size": len(data),
            })
            return err
        }
        
        // Process data
        // ...
        
        appLogger.Info(ctx, "Data processed successfully", map[string]interface{}{
            "data_size": len(data),
            "processed_bytes": len(data),
        })
        
        return nil
    }
}
```

---

## 🔧 **Configuration Examples by Environment**

### **Development Environment**
```yaml
log:
  level: "debug"
  format: "console"
  output: "file"
  filePath: "./logs"
  maxSize: 10
  maxBackup: 2
  maxAge: 1
  compress: false
  addCaller: true
  addStack: true
```

### **Testing Environment**
```yaml
log:
  level: "info"
  format: "json"
  output: "file"
  filePath: "./logs"
  maxSize: 5
  maxBackup: 1
  maxAge: 1
  compress: false
  addCaller: true
  addStack: false
```

### **Production Environment**
```yaml
log:
  level: "info"
  format: "json"
  output: "file"
  filePath: "/var/log/app"
  maxSize: 100
  maxBackup: 5
  maxAge: 30
  compress: true
  addCaller: true
  addStack: false
```

### **High-Traffic Environment**
```yaml
log:
  level: "warn"
  format: "json"
  output: "file"
  filePath: "/var/log/app"
  maxSize: 500
  maxBackup: 10
  maxAge: 7
  compress: true
  addCaller: false
  addStack: false
```

---

## 📊 **Monitoring and Analysis Examples**

### **Log Analysis Commands**

```bash
# Check log directory size
du -sh ./logs/

# Count log files
ls -la ./logs/ | wc -l

# Find errors
grep -r "ERROR" ./logs/

# Count errors by hour
grep -r "ERROR" ./logs/ | cut -d' ' -f2 | cut -d':' -f1 | sort | uniq -c

# Find specific error patterns
grep -r "connection timeout" ./logs/

# Monitor real-time logs
tail -f ./logs/app_2025-01-19_14-30-25.log

# Search across all log files
grep -r "user_id.*12345" ./logs/

# Check file ages
find ./logs/ -name "*.gz" -mtime +7 -ls
```

### **Log Rotation Monitoring**

```bash
# Check current file size
ls -lah ./logs/app_2025-01-19_14-30-25.log

# Check backup files
ls -lah ./logs/app_2025-01-19_14-30-25.log.*.gz

# Monitor rotation events
grep -r "rotation" ./logs/ | tail -10
```

---

## 🚀 **Performance Examples**

### **High-Throughput Logging**

```go
package main

import (
    "context"
    "sync"
    "tushartemplategin/pkg/logger"
)

func main() {
    // High-performance configuration
    config := &logger.Config{
        Level:      "info",
        Format:     "json",
        Output:     "file",
        FilePath:   "./logs",
        MaxSize:    500,  // Larger files for high throughput
        MaxBackups: 10,   // More backups
        MaxAge:     3,    // Shorter retention
        Compress:   true,
        AddCaller:  false, // Disable caller info for performance
        AddStack:   false,
    }
    
    appLogger, _ := logger.NewLogger(config)
    
    // Concurrent logging
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                appLogger.Info(context.Background(), "High throughput log", map[string]interface{}{
                    "goroutine_id": id,
                    "iteration": j,
                })
            }
        }(i)
    }
    wg.Wait()
}
```

---

## ✅ **Summary**

This document provides comprehensive examples for:

- **Basic Setup**: Simple logger configuration and usage
- **File Management**: How files are created, rotated, and managed
- **Log Formats**: JSON and console format examples
- **Use Cases**: Web applications, database operations, error handling
- **Environment Configs**: Development, testing, production, high-traffic
- **Monitoring**: Commands and techniques for log analysis
- **Performance**: High-throughput logging examples

These examples demonstrate the flexibility and power of the timestamp-based logging system for various production scenarios.
