package kbserver

import (
	"net/http"
	"strings"
)

func AllowSubdomainCORS(domain string, server http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow same domain-origin requests
		origin := r.Header.Get("Origin")
		if origin == domain || strings.HasSuffix(origin, "."+domain) {
			w.Header().Set("Access-Control-Allow-Methods", "PUT, GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		server.ServeHTTP(w, r)
	})
}
