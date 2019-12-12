package middlewares

import "net/http"

type CorsMiddleware struct {
	handler http.Handler
}

func Cors(handler http.Handler) *CorsMiddleware {
	return &CorsMiddleware{handler}
}

func (h *CorsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	h.handler.ServeHTTP(w, r)
}
