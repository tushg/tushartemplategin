# Security Middleware

This package provides security-related middleware for the Gin web framework to fix common DRP (Dynamic Risk Profile) issues and improve application security.

## Available Middleware

### 1. SecurityHeaders()
**Purpose:** Adds comprehensive security headers to all responses
**Fixes:** Multiple DRP issues including Missing_HSTS_Header

**Headers Added:**
- `Strict-Transport-Security` - HSTS header (fixes Missing_HSTS_Header)
- `X-Content-Type-Options` - Prevents MIME type sniffing
- `X-Frame-Options` - Prevents clickjacking attacks
- `X-XSS-Protection` - Enables XSS filtering
- `Referrer-Policy` - Controls referrer information
- `Content-Security-Policy` - Prevents XSS and injection attacks
- `Permissions-Policy` - Restricts browser features

### 2. HSTSOnly()
**Purpose:** Adds only HSTS header (minimal security)
**Use Case:** When you only want to fix the Missing_HSTS_Header issue

### 3. CORS()
**Purpose:** Adds Cross-Origin Resource Sharing headers
**Use Case:** When your API needs to be accessed from different origins

## Usage

### Basic Security Headers
```go
import "tushartemplategin/pkg/middleware"

// Add to your router
router.Use(middleware.SecurityHeaders())
```

### HSTS Only (Minimal)
```go
router.Use(middleware.HSTSOnly())
```

### CORS Support
```go
router.Use(middleware.CORS())
```

### Combined Usage
```go
router.Use(middleware.SecurityHeaders())
router.Use(middleware.CORS())
```

## DRP Issues Fixed

1. **Missing_HSTS_Header** ✅ - Fixed by `Strict-Transport-Security` header
2. **Clickjacking** ✅ - Fixed by `X-Frame-Options` header
3. **XSS Attacks** ✅ - Fixed by `X-XSS-Protection` and `Content-Security-Policy` headers
4. **MIME Type Sniffing** ✅ - Fixed by `X-Content-Type-Options` header

## Security Headers Explained

### HSTS (HTTP Strict Transport Security)
- **Value:** `max-age=31536000; includeSubDomains; preload`
- **max-age=31536000:** Tells browsers to use HTTPS for 1 year
- **includeSubDomains:** Applies to all subdomains
- **preload:** Allows inclusion in browser HSTS preload lists

### Content Security Policy
- **default-src 'self':** Only allow resources from same origin
- **script-src 'self' 'unsafe-inline':** Allow inline scripts (customize as needed)
- **style-src 'self' 'unsafe-inline':** Allow inline styles (customize as needed)

## Customization

You can customize the security headers by modifying the middleware functions or creating your own based on your specific security requirements.

## Testing

After adding the middleware, you can test the headers using:
```bash
curl -I http://localhost:8080/api/v1/health
```

You should see the security headers in the response.

## Production Considerations

1. **HSTS:** Only enable in production with HTTPS
2. **CSP:** Customize based on your application's resource requirements
3. **CORS:** Restrict origins in production, don't use `*` for sensitive APIs
4. **Testing:** Use security scanning tools to verify headers are working
