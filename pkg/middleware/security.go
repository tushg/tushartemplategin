package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds comprehensive security-related HTTP headers to all responses
// This middleware implements industry-standard security headers to protect against
// common web vulnerabilities and comply with security frameworks like OWASP Top 10
//
// Security Headers Added:
// 1. HSTS (HTTP Strict Transport Security) - Forces HTTPS usage
// 2. X-Content-Type-Options - Prevents MIME type sniffing attacks
// 3. X-Frame-Options - Prevents clickjacking attacks
// 4. X-XSS-Protection - Enables browser XSS filtering
// 5. Referrer-Policy - Controls referrer information leakage
// 6. Content-Security-Policy - Prevents XSS and injection attacks
// 7. Permissions-Policy - Restricts browser feature access
//
// Usage:
// router.Use(middleware.SecurityHeaders())
//
// Production Considerations:
// - HSTS should only be enabled when HTTPS is available
// - CSP policy should be customized based on application needs
// - Headers should be tested with security scanning tools
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// ===== HSTS (HTTP Strict Transport Security) =====
		// Purpose: Forces browsers to use HTTPS for all future requests
		// How it works: Browser remembers this directive and automatically converts
		// HTTP requests to HTTPS for the specified duration
		//
		// Header Value Breakdown:
		// - max-age=31536000: Tells browser to enforce HTTPS for 1 year (31,536,000 seconds)
		// - includeSubDomains: Applies HSTS to all subdomains (e.g., api.example.com)
		// - preload: Allows inclusion in browser HSTS preload lists for maximum security
		//
		// Real-world Example:
		// User visits http://example.com → Browser sees HSTS header →
		// Browser automatically redirects to https://example.com for 1 year
		//
		// Security Impact: Prevents man-in-the-middle attacks, protocol downgrade attacks
		// DRP Issue Fixed: Missing_HSTS_Header vulnerability

		// Production HSTS with preload support
		if c.Request.TLS != nil {
			// HTTPS request - set strong HSTS
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		} else {
			// HTTP request - set HSTS without preload (allows HTTP to HTTPS redirect)
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// ===== X-Content-Type-Options =====
		// Purpose: Prevents MIME type sniffing attacks
		// How it works: Tells browser to trust the Content-Type header and not
		// try to guess the file type based on content
		//
		// Attack Scenario:
		// Attacker uploads malicious JavaScript file with .jpg extension
		// Without this header: Browser might execute the JS code
		// With this header: Browser respects Content-Type and treats it as image
		//
		// Security Impact: Prevents execution of malicious files disguised as safe types
		// Real-world Example: Prevents .exe files from running when disguised as .txt
		c.Header("X-Content-Type-Options", "nosniff")

		// ===== X-Frame-Options =====
		// Purpose: Prevents clickjacking attacks
		// How it works: Controls whether the page can be embedded in iframes
		// on other domains
		//
		// Attack Scenario:
		// Attacker creates malicious site with invisible iframe containing your login page
		// User thinks they're clicking on attacker's button, but actually clicking
		// on your login form, potentially revealing credentials
		//
		// Values Explained:
		// - DENY: Page cannot be embedded in any iframe (most secure)
		// - SAMEORIGIN: Page can only be embedded in iframes on same domain
		// - ALLOW-FROM uri: Page can only be embedded from specific URI
		//
		// Security Impact: Prevents clickjacking, UI redressing attacks
		// Real-world Example: Prevents Facebook login from being embedded in phishing sites
		c.Header("X-Frame-Options", "DENY")

		// ===== X-XSS-Protection =====
		// Purpose: Enables browser's built-in XSS filtering
		// How it works: Browser detects and blocks reflected XSS attacks
		// Note: This is a legacy header, modern protection comes from CSP
		//
		// Values Explained:
		// - 1: Enable XSS filtering
		// - 0: Disable XSS filtering
		// - 1; mode=block: Enable filtering and block the page if attack detected
		//
		// Attack Scenario:
		// URL: https://example.com/search?q=<script>alert('xss')</script>
		// Without protection: Script executes in user's browser
		// With protection: Browser detects and blocks the script
		//
		// Security Impact: Provides additional layer of XSS protection
		// Real-world Example: Blocks reflected XSS in search forms, comment systems
		c.Header("X-XSS-Protection", "1; mode=block")

		// ===== Referrer-Policy =====
		// Purpose: Controls how much referrer information is sent with requests
		// How it works: Determines what referrer data is included when user
		// navigates from your site to another site
		//
		// Values Explained:
		// - strict-origin-when-cross-origin: Send full referrer to same origin,
		//   only origin (not path) to cross-origin, no referrer to less secure destinations
		//
		// Privacy Impact: Prevents sensitive information leakage in referrer headers
		// Real-world Example: Prevents user IDs, search terms from appearing in
		// analytics on third-party sites
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// ===== Content-Security-Policy (CSP) =====
		// Purpose: Prevents XSS, injection attacks, and controls resource loading
		// How it works: Defines allowed sources for scripts, styles, images, etc.
		// Browser blocks any resources that don't match the policy
		//
		// Policy Breakdown:
		// - default-src 'self': Only allow resources from same origin (most restrictive)
		// - script-src 'self' 'unsafe-inline': Allow scripts from same origin + inline scripts
		// - style-src 'self' 'unsafe-inline': Allow styles from same origin + inline styles
		// - img-src 'self' data: https:: Allow images from same origin, data URIs, and HTTPS
		// - font-src 'self': Only allow fonts from same origin
		// - connect-src 'self': Only allow AJAX/fetch to same origin
		//
		// Attack Scenarios Prevented:
		// 1. XSS: Malicious scripts from external sources are blocked
		// 2. Data Injection: External resources can't be loaded
		// 3. Clickjacking: External iframes are blocked
		//
		// Security Impact: Most effective protection against XSS and injection attacks
		// Real-world Example: Prevents malicious ads from loading scripts on banking sites
		//
		// Customization Needed: Adjust based on your application's requirements
		// - CDN usage: Add CDN domains to appropriate src directives
		// - Third-party integrations: Add trusted domains
		// - Analytics: Add analytics domains to script-src
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self';")

		// ===== Permissions-Policy =====
		// Purpose: Restricts browser features and APIs that could be abused
		// How it works: Controls access to sensitive browser features like
		// geolocation, camera, microphone, etc.
		//
		// Features Restricted:
		// - geolocation=(): No access to user's location
		// - microphone=(): No access to microphone
		// - camera=(): No access to camera
		//
		// Attack Scenarios Prevented:
		// 1. Location Tracking: Malicious sites can't access user location
		// 2. Audio/Video Recording: Sites can't secretly record users
		// 3. Privacy Invasion: Prevents unauthorized access to sensitive devices
		//
		// Security Impact: Protects user privacy and prevents device abuse
		// Real-world Example: Prevents news sites from accessing user's camera
		//
		// Customization: Add more features as needed:
		// - payment=(): Restrict payment APIs
		// - usb=(): Restrict USB device access
		// - bluetooth=(): Restrict Bluetooth access
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Continue to next middleware/handler
		c.Next()
	})
}

// HSTSOnly adds only HSTS header (minimal security approach)
// Use this when you only want to fix the Missing_HSTS_Header DRP issue
// without adding other security headers
//
// Use Cases:
// - Quick security fix for HSTS requirement
// - When other headers are handled by reverse proxy (nginx, CloudFlare)
// - Development/testing environments
//
// Security Level: Basic (only HTTPS enforcement)
func HSTSOnly() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// HSTS header with 1-year duration and subdomain coverage
		// This is the minimum required to fix Missing_HSTS_Header DRP issue
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Next()
	})
}

// CORS adds Cross-Origin Resource Sharing headers
// Enables your API to be accessed from different origins (domains)
//
// Use Cases:
// - Frontend applications on different domains
// - Mobile apps consuming your API
// - Third-party integrations
// - Development with separate frontend/backend
//
// Security Considerations:
// - Access-Control-Allow-Origin: "*" allows any domain (not recommended for production)
// - Consider restricting to specific domains in production
// - Credentials should only be allowed for trusted origins
//
// Production Recommendations:
// - Restrict origins to specific domains
// - Use environment-based configuration
// - Monitor CORS usage for security analysis
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Allow requests from any origin (customize for production)
		// In production, consider restricting to specific domains:
		// c.Header("Access-Control-Allow-Origin", "https://trusted-domain.com")
		c.Header("Access-Control-Allow-Origin", "*")

		// Allow common HTTP methods
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow common headers including custom ones
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Expose response headers to client
		c.Header("Access-Control-Expose-Headers", "Content-Length")

		// Allow credentials (cookies, authorization headers)
		// Note: When using credentials, Access-Control-Allow-Origin cannot be "*"
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS requests
		// Preflight requests are sent by browsers before actual requests
		// to check if the cross-origin request is allowed
		if c.Request.Method == "OPTIONS" {
			// Return 204 No Content for preflight requests
			// This tells the browser the request is allowed
			c.AbortWithStatus(204)
			return
		}

		// Continue to next middleware/handler
		c.Next()
	})
}
