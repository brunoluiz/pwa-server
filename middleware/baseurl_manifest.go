package middleware

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/brunoluiz/go-pwa-server/manifestmod"
)

// ManifestBaseURL changes manifest.json start_url and scope using baseURL
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
				panic("Fail to open manifest.json")
			}
			defer f.Close()

			buf, err := ioutil.ReadAll(f)
			if err != nil {
				panic("error on reading manifest.json")
			}

			js, err := manifestmod.ChangeBaseURL(buf, baseURL)
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(js); err != nil {
				panic("error on sending manifest.json")
			}
		})
	}
}
