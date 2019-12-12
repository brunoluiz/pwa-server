package middleware

import "net/http"

type CorsMiddleware struct {
	next http.Handler
}

func Cors(next http.Handler) *CorsMiddleware {
	return &CorsMiddleware{next}
}

func (h *CorsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	h.next.ServeHTTP(w, r)
}
