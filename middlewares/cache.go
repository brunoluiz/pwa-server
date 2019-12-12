package middlewares

import "net/http"

type NoCacheMiddleware struct {
	handler http.Handler
}

func NoCache(handler http.Handler) *NoCacheMiddleware {
	return &NoCacheMiddleware{handler}
}

func (h *NoCacheMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")
	h.handler.ServeHTTP(w, r)
}
