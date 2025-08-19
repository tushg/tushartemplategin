# Production SSL/TLS Setup Guide

This guide explains how to configure your Go microservice with SSL/TLS for production deployment, resolving HSTS header issues and ensuring security compliance.

## 🎯 **What This Implementation Provides**

- **✅ SSL/TLS Support**: Native HTTPS without reverse proxy
- **✅ HSTS Headers**: Proper HTTP Strict Transport Security implementation
- **✅ Certificate Management**: Support for Let's Encrypt and commercial certificates
- **✅ HTTP to HTTPS Redirect**: Automatic redirection for security
- **✅ Production Security**: TLS 1.2+, strong cipher suites, security headers
- **✅ Zero Downtime**: Graceful shutdown and restart capabilities

## 🚀 **Quick Start for Production**

### 1. **Enable SSL in Configuration**

Update your main configuration file (`configs/config.json`):

```json
{
  "server": {
    "port": ":8080",
    "mode": "release",
    "ssl": {
      "enabled": true,
      "port": ":443",
      "certFile": "/etc/ssl/certs/your-domain.crt",
      "keyFile": "/etc/ssl/private/your-domain.key",
      "redirectHTTP": true
    }
  }
}
```

### 2. **Obtain SSL Certificate**

#### **Option A: Let's Encrypt (Recommended for Production)**

```bash
# Install certbot
sudo apt-get update
sudo apt-get install certbot

# Get certificate (replace 'yourdomain.com' with your actual domain)
sudo certbot certonly --standalone -d yourdomain.com

# Copy certificates to secure location
sudo mkdir -p /etc/ssl/certs /etc/ssl/private
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem /etc/ssl/certs/your-domain.crt
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem /etc/ssl/private/your-domain.key

# Set proper permissions
sudo chown root:root /etc/ssl/certs/your-domain.crt /etc/ssl/private/your-domain.key
sudo chmod 644 /etc/ssl/certs/your-domain.crt
sudo chmod 600 /etc/ssl/private/your-domain.key
```

#### **Option B: Commercial Certificate**

1. Purchase certificate from CA (DigiCert, Comodo, etc.)
2. Download and convert to PEM format
3. Place in `/etc/ssl/certs/` and `/etc/ssl/private/`

### 3. **Deploy and Test**

```bash
# Build for production
go build -o server cmd/server/main.go

# Run with main config (SSL enabled)
./server

# Test HTTPS
curl -k https://yourdomain.com/health
```

## 🔒 **Security Features Implemented**

### **TLS Configuration**
- **Minimum TLS Version**: 1.2 (TLS 1.3 supported)
- **Cipher Suites**: Strong, modern ciphers only
- **HTTP/2 Support**: Enabled for better performance
- **Certificate Validation**: Automatic validation on startup

### **Security Headers**
- **HSTS**: HTTP Strict Transport Security with preload support
- **X-Content-Type-Options**: Prevents MIME type sniffing
- **X-Frame-Options**: Prevents clickjacking
- **X-XSS-Protection**: XSS protection
- **Referrer-Policy**: Controls referrer information
- **Content-Security-Policy**: Resource loading restrictions

### **Production Features**
- **Graceful Shutdown**: 30-second timeout for active connections
- **Connection Pooling**: Efficient resource management
- **Error Handling**: Comprehensive error logging and recovery
- **Health Checks**: Built-in health monitoring endpoints

## 📋 **Production Checklist**

### **Before Deployment**
- [ ] Valid SSL certificate obtained and installed
- [ ] Certificate files have correct permissions (600 for key, 644 for cert)
- [ ] Domain DNS points to your server
- [ ] Firewall allows traffic on ports 80 and 443
- [ ] SSL configuration enabled in `config.json`

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

## 🧪 **Testing Your SSL Setup**

### **1. SSL Labs Test**
Visit [SSL Labs](https://www.ssllabs.com/ssltest/) and test your domain for an A+ rating.

### **2. Security Headers Test**
```bash
# Test HTTPS endpoint
curl -I https://yourdomain.com/health

# Verify HSTS header
curl -I https://yourdomain.com/health | grep -i "strict-transport-security"

# Test HTTP to HTTPS redirect
curl -I http://yourdomain.com/health
```

### **3. Certificate Validation**
```bash
# Check certificate details
openssl x509 -in /etc/ssl/certs/your-domain.crt -text -noout

# Test TLS connection
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com
```

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
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem /etc/ssl/certs/your-domain.crt
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem /etc/ssl/private/your-domain.key

# Restart your service
sudo systemctl restart your-service
```

## 🚨 **Troubleshooting Common Issues**

### **1. Certificate Not Found**
```bash
# Check file existence
ls -la /etc/ssl/certs/your-domain.crt
ls -la /etc/ssl/private/your-domain.key

# Check permissions
ls -la /etc/ssl/private/your-domain.key | grep "^-rw-------"
```

### **2. Permission Denied**
```bash
# Fix permissions
sudo chmod 600 /etc/ssl/private/your-domain.key
sudo chmod 644 /etc/ssl/certs/your-domain.crt
sudo chown root:root /etc/ssl/private/your-domain.key /etc/ssl/certs/your-domain.crt
```

### **3. Port Already in Use**
```bash
# Check what's using port 443
sudo netstat -tlnp | grep :443
sudo lsof -i :443

# Kill process if needed
sudo kill -9 <PID>
```

### **4. HSTS Not Working**
- Ensure HTTPS is accessible
- Check certificate validity
- Verify security middleware is loaded
- Test with browser developer tools

## 📊 **Performance Considerations**

### **TLS Performance**
- **Session Resumption**: Enabled for faster connections
- **OCSP Stapling**: Reduces certificate validation overhead
- **HTTP/2**: Better multiplexing and compression
- **Connection Pooling**: Efficient resource management

### **Monitoring**
```bash
# Monitor SSL connections
netstat -an | grep :443 | wc -l

# Check certificate expiration
openssl x509 -in /etc/ssl/certs/your-domain.crt -noout -dates

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

---

**🎯 Your microservice is now production-ready with enterprise-grade SSL/TLS security!**
