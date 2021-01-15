package request

import (
	"net"
	"net/http"
	"strings"
)

// CuvvaClientIPHeader is the customer header, set by our network edge,
// that is expected to be the clients real IP address when X-Forwarded-For
// cannot be trusted.
const CuvvaClientIPHeader = `Cuvva-Client-Ip`

// CuvvaClientIP respects the Cuvva-Client-IP header containing the clients
// real IP address.
func CuvvaClientIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cci := strings.TrimSpace(r.Header.Get(CuvvaClientIPHeader))
		if cci != "" {
			// assert valid IP address
			ip := net.ParseIP(cci)
			if ip != nil {
				// marshal from parse string to restrict to supported encoding
				r.RemoteAddr = ip.String()
			}
		}

		next.ServeHTTP(w, r)
	})
}
