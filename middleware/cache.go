package middleware

import "net/http"

type NoCacheMiddleware struct {
	next http.Handler
}

func NoCache(next http.Handler) *NoCacheMiddleware {
	return &NoCacheMiddleware{next}
}

func (h *NoCacheMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")
	h.next.ServeHTTP(w, r)
}
