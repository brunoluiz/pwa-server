package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/brunoluiz/pwa-server/htmlmod"
	"golang.org/x/net/html"
)

// HTMLBaseURL add <base href='baseUrl'> tag to html->head, with the specified base url
func HTMLBaseURL(dir string, baseURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isHTML := strings.Contains(r.URL.Path, "htm")
			hasTrailing := r.URL.Path[len(r.URL.Path)-1:] == "/"
			if !isHTML && !hasTrailing {
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
}
