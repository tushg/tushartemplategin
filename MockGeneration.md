# Mock Generation Process for Unit Testing

This document outlines the step-by-step process for generating mocks using `gomock` for unit testing in the Tushar Template Gin project.

## Prerequisites

### 1. Install mockgen tool
```bash
go install github.com/golang/mock/mockgen@latest
```

### 2. Verify installation
```bash
mockgen --version
```

**Note**: The `gomock` library is already included in your `go.mod` file, so you don't need to run `go get github.com/golang/mock/gomock`.

## Step-by-Step Mock Generation Process

### Step 1: Generate Mocks for Database Interfaces
```bash
mockgen -source=E:\tushartemplategin\pkg\interfaces\database.go -destination=mocks/mock_database.go
```

### Step 2: Generate Mocks for Logger Interface
```bash
mockgen -source=E:\tushartemplategin\pkg\interfaces\logger.go -destination=mocks/mock_logger.go
```

### Step 3: Generate Mocks for Config Interfaces
```bash
mockgen -source=E:\tushartemplategin\pkg\interfaces\config.go -destination=mocks/mock_config.go
```

### Step 4: Fix Package Declarations (Important!)
After generation, you need to manually fix the package declarations in each mock file:

**In `mocks/mock_database.go`:**
```go
// Change from:
// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

// To:
// Package mocks is a generated GoMock package.
package mocks
```

**In `mocks/mock_logger.go`:**
```go
// Change from:
package mock_interfaces

// To:
package mocks
```

**In `mocks/mock_config.go`:**
```go
// Change from:
package mock_interfaces

// To:
package mocks
```

### Step 5: Verify Mock Generation
Check that all mock files were created:
```bash
ls mocks/
# Should show:
# mock_database.go
# mock_logger.go
# mock_config.go
```

### Step 6: Test the Mocks
Build the project to ensure mocks work:
```bash
go build ./...
```

### Step 7: Run Tests to Verify
```bash
go test ./... -v
```

## Complete Script (Copy-Paste Ready)

### Windows PowerShell Script
```powershell
# Navigate to your project directory
cd E:\tushartemplategin

# Remove old mocks (if they exist)
Remove-Item -Recurse -Force mocks/ -ErrorAction SilentlyContinue

# Create mocks directory
mkdir mocks

# Generate mocks
mockgen -source=E:\tushartemplategin\pkg\interfaces\database.go -destination=mocks/mock_database.go
mockgen -source=E:\tushartemplategin\pkg\interfaces\logger.go -destination=mocks/mock_logger.go
mockgen -source=E:\tushartemplategin\pkg\interfaces\config.go -destination=mocks/mock_config.go

# Verify files were created
ls mocks/

# Build to check for errors
go build ./...

# Run tests
go test ./... -v
```

### Linux/Mac Script
```bash
#!/bin/bash

# Navigate to your project directory
cd /path/to/your/project

# Remove old mocks (if they exist)
rm -rf mocks/

# Create mocks directory
mkdir mocks

# Generate mocks
mockgen -source=pkg/interfaces/database.go -destination=mocks/mock_database.go
mockgen -source=pkg/interfaces/logger.go -destination=mocks/mock_logger.go
mockgen -source=pkg/interfaces/config.go -destination=mocks/mock_config.go

# Verify files were created
ls -la mocks/

# Fix package declarations
sed -i 's/package mock_interfaces/package mocks/g' mocks/*.go

# Build to check for errors
go build ./...

# Run tests
go test ./... -v
```

## Makefile Alternative (Recommended for Linux)

Create a `Makefile` in your project root:

```makefile
.PHONY: mocks test build clean

# Generate all mocks
mocks:
	@echo "🔧 Generating mocks..."
	rm -rf mocks/
	mkdir mocks
	mockgen -source=pkg/interfaces/database.go -destination=mocks/mock_database.go
	mockgen -source=pkg/interfaces/logger.go -destination=mocks/mock_logger.go
	mockgen -source=pkg/interfaces/config.go -destination=mocks/mock_config.go
	@echo "🔧 Fixing package declarations..."
	sed -i 's/package mock_interfaces/package mocks/g' mocks/*.go
	@echo "✅ Mocks generated successfully!"

# Build project
build: mocks
	@echo "🔨 Building project..."
	go build ./...
	@echo "✅ Build successful!"

# Run tests
test: build
	@echo "🧪 Running tests..."
	go test ./... -v

# Clean generated files
clean:
	@echo "🧹 Cleaning up..."
	rm -rf mocks/
	@echo "✅ Cleanup complete!"

# Full workflow
all: clean mocks build test
```

### Usage with Makefile
```bash
# Generate mocks only
make mocks

# Build project (includes mock generation)
make build

# Run tests (includes build and mock generation)
make test

# Full workflow (clean, generate, build, test)
make all

# Clean up
make clean
```

## When to Regenerate Mocks

### Regenerate mocks when you:
1. **Add new methods** to interfaces
2. **Change method signatures** in interfaces
3. **Add new interfaces** to `pkg/interfaces/`
4. **Update Go version** (mockgen compatibility)
5. **Change interface package structure**

### You DON'T need to regenerate when:
1. **Adding new implementations** (concrete types)
2. **Changing business logic** in implementations
3. **Adding new domains** (unless they have new interfaces)
4. **Updating tests** (mocks stay the same)

## Troubleshooting Common Issues

### Issue 1: "package mock_interfaces is not in std"
**Solution**: Fix package declarations (Step 4 above)

### Issue 2: "undefined: interfaces.HealthStatus"
**Solution**: Regenerate mocks after interface changes

### Issue 3: "mockgen failed with import cycle"
**Solution**: Ensure interfaces are in `pkg/interfaces/` not in domain packages

### Issue 4: "cannot use mock as interface"
**Solution**: Check that mock package name is `mocks` not `mock_interfaces`

### Issue 5: "mockgen command not found"
**Solution**: Install mockgen with `go install github.com/golang/mock/mockgen@latest`

## Best Practices

1. **✅ Always regenerate mocks** after interface changes
2. **✅ Keep interfaces centralized** in `pkg/interfaces/`
3. **✅ Use consistent package naming** (`package mocks`)
4. **✅ Test mocks work** by building the project
5. **✅ Commit mocks to version control** (they're generated code)
6. **✅ Use Makefile** for consistent workflow (Linux/Mac)

## Quick Reference Commands

```bash
# Generate all mocks at once
mockgen -source=pkg/interfaces/database.go -destination=mocks/mock_database.go && mockgen -source=pkg/interfaces/logger.go -destination=mocks/mock_logger.go && mockgen -source=pkg/interfaces/config.go -destination=mocks/mock_config.go

# Check mock status
go build ./... && go test ./... -v

# Quick mock generation (Linux/Mac)
make mocks

# Full test workflow (Linux/Mac)
make test
```

## Project Structure After Mock Generation

```
tushartemplategin/
├── pkg/
│   └── interfaces/           # Interface definitions
│       ├── database.go       # Database interfaces
│       ├── logger.go         # Logger interface
│       └── config.go         # Config interfaces
├── mocks/                    # Generated mock files
│   ├── mock_database.go      # Database mocks
│   ├── mock_logger.go        # Logger mocks
│   └── mock_config.go        # Config mocks
└── internal/                 # Domain implementations
    └── health/               # Health domain
        ├── repository_test.go # Tests using mocks
        └── ...
```

## Summary

Follow these steps to generate mocks for your unit tests:

1. **Install mockgen** (if not already installed)
2. **Generate mocks** for each interface file
3. **Fix package declarations** to use `package mocks`
4. **Verify generation** by building the project
5. **Test mocks** by running unit tests
6. **Use Makefile** for consistent workflow (recommended)

This process ensures your mocks are always up-to-date with your interfaces and your unit tests can run successfully! 🚀
