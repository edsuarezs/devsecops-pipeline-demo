package middleware

import "net/http"

// SecurityHeaders adds OWASP-recommended security headers to every response.
// These headers mitigate common web vulnerabilities:
//   - X-Content-Type-Options: prevents MIME-type sniffing
//   - X-Frame-Options: prevents clickjacking
//   - X-XSS-Protection: legacy XSS filter (still useful for older browsers)
//   - Content-Security-Policy: restricts resource loading
//   - Cache-Control: prevents sensitive data caching
//   - Strict-Transport-Security: enforces HTTPS
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}
