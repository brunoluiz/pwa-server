package middlewares

import (
	"fmt"
	"net/http"
)

type BasicAuthMiddleware struct {
	realm       string
	credentials map[string][]string
	handler     http.Handler
}

func BasicAuth(
	realm string,
	credentials map[string][]string,
	handler http.Handler,
) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{realm, credentials, handler}
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
			h.handler.ServeHTTP(w, r)
			return
		}
	}

	h.unauthorized(w)
}

func (h *BasicAuthMiddleware) unauthorized(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, h.realm))
	w.WriteHeader(http.StatusUnauthorized)
}
