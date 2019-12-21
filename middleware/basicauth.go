package middleware

import (
	"fmt"
	"net/http"
)

// BasicAuth adds a middleare for basic auth
func BasicAuth(
	realm string,
	credentials map[string][]string,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, headerFound := r.BasicAuth()
		if !headerFound {
			unauthorised(w, realm)
			return
		}

		validPasswords, userFound := credentials[user]
		if !userFound {
			unauthorised(w, realm)
			return
		}

		for _, validPassword := range validPasswords {
			if password == validPassword {
				next.ServeHTTP(w, r)
				return
			}
		}

		unauthorised(w, realm)
	})
}

func unauthorised(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}
