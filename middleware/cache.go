package middleware

import "net/http"

// NoCache set cache headers to no-cache, forcing the client to always request the latest assets
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "-1")
		next.ServeHTTP(w, r)
	})
}
