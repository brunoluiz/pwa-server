package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/brunoluiz/go-pwa-server/htmlmod"
	"golang.org/x/net/html"
)

func HTMLBaseURL(dir string, baseURL string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hasTrailing := r.URL.Path[len(r.URL.Path)-1:] == "/"
		if !strings.Contains(r.URL.Path, "htm") && !hasTrailing {
			next.ServeHTTP(w, r)
			return
		}

		path := dir + r.URL.Path
		if hasTrailing {
			path += "index.html"
		}

		f, err := os.Open(path)
		if err != nil {
			panic("Fail to open")
		}
		defer f.Close()

		doc, err := html.Parse(f)
		if err != nil {
			panic("Fail to parse")
		}

		htmlmod.EnhanceHTML(doc, &htmlmod.Settings{BaseURL: baseURL})
		if err := html.Render(w, doc); err != nil {
			panic(err)
		}
	})
}
