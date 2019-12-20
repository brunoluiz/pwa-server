package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/brunoluiz/go-pwa-server/manifestmod"
)

func ManifestBaseURL(
	dir string,
	baseURL string,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "manifest.json") {
				next.ServeHTTP(w, r)
				return
			}

			f, err := os.Open(dir + r.URL.Path)
			if err != nil {
				panic("Fail to open")
			}
			defer f.Close()

			js, err := manifestmod.ChangeBaseURL(f, baseURL)
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(js); err != nil {
				panic("error on sending data")
			}
		})
	}
}
