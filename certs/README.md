# SSL/TLS Certificates for Production

This directory contains SSL/TLS certificates for secure HTTPS communication.

## Certificate Files

- `server.crt` - SSL certificate file (PEM format)
- `server.key` - SSL private key file (PEM format)

## Production Setup

### 1. Obtain SSL Certificate

#### Option A: Let's Encrypt (Free)
```bash
# Install certbot
sudo apt-get install certbot

# Get certificate for your domain
sudo certbot certonly --standalone -d yourdomain.com

# Copy certificates to this directory
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ./certs/server.crt
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ./certs/server.key

# Set proper permissions
sudo chown $USER:$USER ./certs/server.crt ./certs/server.key
chmod 600 ./certs/server.key
chmod 644 ./certs/server.crt
```

#### Option B: Commercial Certificate
1. Purchase SSL certificate from CA (Comodo, DigiCert, etc.)
2. Download certificate files
3. Convert to PEM format if needed
4. Place in this directory

### 2. Update Configuration

Update your production configuration file:
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

### 3. Certificate Renewal (Let's Encrypt)

Let's Encrypt certificates expire every 90 days. Set up auto-renewal:

```bash
# Test renewal
sudo certbot renew --dry-run

# Set up cron job for auto-renewal
sudo crontab -e

# Add this line (runs twice daily)
0 0,12 * * * /usr/bin/certbot renew --quiet
```

## Security Best Practices

1. **Private Key Security**
   - Keep private key secure and never share
   - Use file permissions 600 (owner read/write only)
   - Consider using hardware security modules (HSM) in production

2. **Certificate Validation**
   - Ensure certificate matches your domain
   - Check certificate expiration dates
   - Use strong key sizes (2048+ bits for RSA)

3. **TLS Configuration**
   - Minimum TLS 1.2 (TLS 1.3 recommended)
   - Use strong cipher suites
   - Enable HSTS headers

## Testing SSL

```bash
# Test certificate validity
openssl x509 -in ./certs/server.crt -text -noout

# Test private key
openssl rsa -in ./certs/server.key -check

# Test TLS connection
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com
```

## Troubleshooting

### Common Issues

1. **Permission Denied**
   ```bash
   chmod 600 ./certs/server.key
   chmod 644 ./certs/server.crt
   ```

2. **Certificate Not Found**
   - Verify file paths in configuration
   - Check file permissions
   - Ensure files exist and are readable

3. **TLS Handshake Failed**
   - Verify certificate validity
   - Check certificate chain
   - Ensure private key matches certificate

## Production Checklist

- [ ] Valid SSL certificate obtained
- [ ] Certificate files placed in `./certs/` directory
- [ ] File permissions set correctly
- [ ] Configuration updated with correct paths
- [ ] SSL enabled in production config
- [ ] HTTP to HTTPS redirect enabled
- [ ] HSTS headers configured
- [ ] Certificate renewal process established
- [ ] Security headers implemented
- [ ] TLS 1.2+ enforced
