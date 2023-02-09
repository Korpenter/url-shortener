package middleware

import (
	"net/http"
	"net/netip"

	"github.com/Mldlr/url-shortener/internal/app/config"
)

// Trusted is a middleware checking if request was made from trusted subnet.
type Trusted struct {
	Config *config.Config
}

// TrustedCheck verifies request if it came from the trusted network
// if request tries to access api/internal/ enpoints
func (t Trusted) TrustCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if subnet was provided
		if t.Config.TrustedSubnet == "" {
			http.Error(w, "untrusted user", http.StatusForbidden)
			return
		}
		// Get the ip header
		xRealIP := r.Header.Get("X-Real-IP")
		// Check if ip is in trusted subnet
		netip, err := netip.ParseAddr(xRealIP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !t.Config.SubnetPrefix.Contains(netip) {
			http.Error(w, "untrusted user", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
