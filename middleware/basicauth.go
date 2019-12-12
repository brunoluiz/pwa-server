package middleware

import (
	"fmt"
	"net/http"
)

type BasicAuthMiddleware struct {
	realm       string
	credentials map[string][]string
	next        http.Handler
}

func BasicAuth(
	realm string,
	credentials map[string][]string,
	next http.Handler,
) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{realm, credentials, next}
}

func (h *BasicAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, password, headerFound := r.BasicAuth()
	if !headerFound {
		h.unauthorized(w)
		return
	}

	validPasswords, userFound := h.credentials[user]
	if !userFound {
		h.unauthorized(w)
		return
	}

	for _, validPassword := range validPasswords {
		if password == validPassword {
			h.next.ServeHTTP(w, r)
			return
		}
	}

	h.unauthorized(w)
}

func (h *BasicAuthMiddleware) unauthorized(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, h.realm))
	w.WriteHeader(http.StatusUnauthorized)
}
