# 🚀 **Complete Guide: Running Your Microservice with SSL/TLS**

This guide provides step-by-step instructions on how to run your Go microservice with SSL/TLS support and test all REST endpoints.

## 📋 **Prerequisites**

- Go 1.24.5 or later
- OpenSSL (for certificate generation)
- Valid SSL certificate files (for production)
- Access to ports 8080 (HTTP) and 443 (HTTPS)

## 🌐 **Available REST Endpoints**

Your microservice provides the following REST endpoints:

- **`GET /api/v1/health`** - Overall health status check
- **`GET /api/v1/health/ready`** - Kubernetes readiness probe  
- **`GET /api/v1/health/live`** - Kubernetes liveness probe

## 🔧 **Step 1: Prepare SSL Certificates**

### **For Development/Testing (Self-signed Certificate)**

1. **Create the certs directory:**
   ```bash
   mkdir -p certs
   ```

2. **Generate self-signed certificate:**
   ```bash
   # Generate self-signed certificate for testing
   openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes -subj "//CN=localhost" -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
   ```

3. **Set proper permissions:**
   ```bash
   chmod 600 certs/server.key
   chmod 644 certs/server.crt
   ```

### **For Production (Let's Encrypt)**

1. **Install certbot:**
   ```bash
   # Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install certbot
   
   # CentOS/RHEL
   sudo yum install certbot
   ```

2. **Obtain certificate:**
   ```bash
   # Get certificate for your domain
   sudo certbot certonly --standalone -d yourdomain.com
   ```

3. **Copy certificates to your project:**
   ```bash
   # Create certs directory
   mkdir -p certs
   
   # Copy certificates
   sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem certs/server.crt
   sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem certs/server.key
   
   # Set ownership and permissions
   sudo chown $USER:$USER certs/server.crt certs/server.key
   chmod 600 certs/server.key
   chmod 644 certs/server.crt
   ```

## 🏃‍♂️ **Step 2: Run the Microservice**

### **Option A: Run Directly (Recommended for Development)**

```bash
# Navigate to your project directory
cd tushartemplategin

# Run the service directly
go run cmd/server/main.go
```

### **Option B: Build and Run (Recommended for Production)**

```bash
# Navigate to your project directory
cd tushartemplategin

# Build the application
go build cmd/server/main.go

# Run the executable
./main.exe
```

## 🌐 **Step 3: Service Access Points**

Your service will be available on:

- **HTTP Port**: `:8080` (for redirects if enabled)
- **HTTPS Port**: `:443` (SSL/TLS enabled)
- **Base URL**: `https://localhost:443` (or your domain in production)

## 🧪 **Step 4: Test the REST Endpoints**

### **Using curl (Command Line)**

#### **Test Health Check:**
```bash
# Test HTTPS endpoint (ignore certificate warnings for self-signed)
curl -k https://localhost:443/api/v1/health

# Test HTTP endpoint (should redirect to HTTPS if enabled)
curl -L http://localhost:8080/api/v1/health
```

#### **Test Readiness Probe:**
```bash
curl -k https://localhost:443/api/v1/health/ready
```

#### **Test Liveness Probe:**
```bash
curl -k https://localhost:443/api/v1/health/live
```

### **Using PowerShell (Windows)**

#### **Test Health Check:**
```powershell
# Test HTTPS endpoint
Invoke-RestMethod -Uri "https://localhost:443/api/v1/health" -SkipCertificateCheck

# Test HTTP endpoint
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/health"
```

#### **Test Readiness Probe:**
```powershell
Invoke-RestMethod -Uri "https://localhost:443/api/v1/health/ready" -SkipCertificateCheck
```

#### **Test Liveness Probe:**
```powershell
Invoke-RestMethod -Uri "https://localhost:443/api/v1/health/live" -SkipCertificateCheck
```

### **Using Web Browser**

Open these URLs in your browser:
- `https://localhost:443/api/v1/health`
- `https://localhost:443/api/v1/health/ready`
- `https://localhost:443/api/v1/health/live`

**Note**: For self-signed certificates, you'll need to accept the security warning.

### **Using Postman**

1. **Create a new collection**
2. **Add requests for each endpoint:**
   - `GET https://localhost:443/api/v1/health`
   - `GET https://localhost:443/api/v1/health/ready`
   - `GET https://localhost:443/api/v1/health/live`
3. **Disable SSL verification** in Postman settings for self-signed certificates

## 📊 **Expected API Responses**

### **Health Check Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-19T15:30:00Z",
  "service": "tushartemplategin",
  "version": "1.0.0",
  "database": {
    "status": "connected",
    "type": "postgres"
  }
}
```

### **Readiness Probe Response:**
```json
{
  "status": "ready",
  "timestamp": "2025-01-19T15:30:00Z",
  "checks": {
    "database": "ready",
    "external_services": "ready"
  }
}
```

### **Liveness Probe Response:**
```json
{
  "status": "alive",
  "timestamp": "2025-01-19T15:30:00Z",
  "uptime": "5m30s"
}
```

## 🔍 **Step 5: Monitor the Service**

### **Check Service Logs**

The service will output detailed logs showing:
- SSL/TLS configuration
- Certificate loading
- Server startup
- Request handling
- Health check results

### **Check Service Status**

#### **Windows (PowerShell):**
```powershell
# Check if ports are listening
netstat -an | Select-String ":443"
netstat -an | Select-String ":8080"

# Check process
Get-Process | Where-Object {$_.ProcessName -eq "main"}
```

#### **Linux/macOS:**
```bash
# Check if ports are listening
netstat -an | grep ":443"
netstat -an | grep ":8080"

# Check process
ps aux | grep "main"
```

## 🛑 **Step 6: Stop the Service**

Press `Ctrl+C` in the terminal where the service is running to gracefully shut it down.

## 🚨 **Troubleshooting Common Issues**

### **1. Certificate Errors**

#### **Check certificate file existence:**
```bash
# Windows
dir certs\server.crt
dir certs\server.key

# Linux/macOS
ls -la certs/server.crt
ls -la certs/server.key
```

#### **Verify certificate format:**
```bash
openssl x509 -in certs/server.crt -text -noout
```

#### **Check certificate validity:**
```bash
openssl x509 -in certs/server.crt -noout -dates
```

### **2. Port Already in Use**

#### **Check what's using the ports:**
```bash
# Windows
netstat -an | findstr ":443"
netstat -an | findstr ":8080"

# Linux/macOS
netstat -an | grep ":443"
netstat -an | grep ":8080"
```

#### **Kill process if needed:**
```bash
# Windows
taskkill /F /PID <PID>

# Linux/macOS
kill -9 <PID>
```

### **3. SSL Not Working**

- Ensure `ssl.enabled: true` in `config.json`
- Check certificate file paths are correct
- Verify certificate files exist and are readable
- Check certificate expiration dates
- Ensure proper file permissions

### **4. Permission Denied**

```bash
# Fix permissions for private key
chmod 600 certs/server.key

# Fix permissions for certificate
chmod 644 certs/server.crt

# Fix ownership if needed
sudo chown $USER:$USER certs/server.key certs/server.crt
```

### **5. HSTS Not Working**

- Ensure HTTPS is accessible
- Check certificate validity
- Verify security middleware is loaded
- Test with browser developer tools

## 🎯 **Quick Test Commands**

Here's a quick test script you can run:

```bash
# Start the service in one terminal
go run cmd/server/main.go

# In another terminal, test all endpoints
curl -k https://localhost:443/api/v1/health
curl -k https://localhost:443/api/v1/health/ready
curl -k https://localhost:443/api/v1/health/live
```

## 📋 **Production Checklist**

### **Before Deployment**
- [ ] Valid SSL certificate obtained and installed
- [ ] Certificate files have correct permissions (600 for key, 644 for cert)
- [ ] Domain DNS points to your server
- [ ] Firewall allows traffic on ports 80 and 443
- [ ] SSL configuration enabled in `config.json`
- [ ] Certificate expiration monitoring set up

### **Security Verification**
- [ ] HTTPS accessible on port 443
- [ ] HTTP redirects to HTTPS (if enabled)
- [ ] HSTS headers present in HTTPS responses
- [ ] Security headers properly set
- [ ] TLS 1.2+ enforced
- [ ] Strong cipher suites used

### **Monitoring & Maintenance**
- [ ] Certificate expiration monitoring
- [ ] Auto-renewal process (Let's Encrypt)
- [ ] SSL Labs A+ rating achieved
- [ ] Regular security header testing
- [ ] Performance monitoring enabled

## 🔄 **Certificate Renewal (Let's Encrypt)**

### **Automatic Renewal**
```bash
# Test renewal process
sudo certbot renew --dry-run

# Set up cron job for auto-renewal
sudo crontab -e

# Add this line (runs twice daily)
0 0,12 * * * /usr/bin/certbot renew --quiet && systemctl reload your-service
```

### **Manual Renewal**
```bash
# Renew certificates
sudo certbot renew

# Copy new certificates
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem certs/server.crt
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem certs/server.key

# Restart your service
sudo systemctl restart your-service
```

## 📊 **Performance Considerations**

### **TLS Performance**
- **Session Resumption**: Enabled for faster connections
- **OCSP Stapling**: Reduces certificate validation overhead
- **HTTP/2**: Better multiplexing and compression
- **Connection Pooling**: Efficient resource management

### **Monitoring Commands**
```bash
# Monitor SSL connections
netstat -an | grep :443 | wc -l

# Check certificate expiration
openssl x509 -in certs/server.crt -noout -dates

# Monitor service status
systemctl status your-service
```

## 🎉 **Benefits of This Implementation**

1. **No Reverse Proxy Needed**: SSL handled directly in Go
2. **Production Ready**: Enterprise-grade security features
3. **HSTS Compliant**: Resolves DRP security issues
4. **Easy Maintenance**: Simple certificate renewal process
5. **High Performance**: Native Go implementation
6. **Security Focused**: Modern TLS and security headers
7. **Scalable**: Handles high traffic efficiently

## 🔗 **Additional Resources**

- [Go TLS Documentation](https://golang.org/pkg/crypto/tls/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [SSL Labs Best Practices](https://github.com/ssllabs/research/wiki/SSL-and-TLS-Deployment-Best-Practices)
- [Mozilla Security Guidelines](https://infosec.mozilla.org/guidelines/web_security)

## 🚀 **Quick Start Summary**

1. **Generate certificates**: `openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes -subj "/CN=localhost"`
2. **Run service**: `go run cmd/server/main.go`
3. **Test endpoints**: `curl -k https://localhost:443/api/v1/health`
4. **Stop service**: `Ctrl+C`

---

**🎯 Your microservice is now production-ready with enterprise-grade SSL/TLS security!**

**Note**: This guide assumes you're running on Windows. For Linux/macOS, replace Windows-specific commands with their Unix equivalents.
