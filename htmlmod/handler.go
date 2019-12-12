package httpmod

import (
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func ServeStatic(dir string, baseUrl string) http.Handler {
	fs := http.FileServer(http.Dir(dir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hasTrailing := r.URL.Path[len(r.URL.Path)-1:] == "/"
		if !strings.Contains(r.URL.Path, "htm") && !hasTrailing {
			fs.ServeHTTP(w, r)
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

		EnhanceHTML(doc, &PWASettings{BaseURL: baseUrl})
		if err := html.Render(w, doc); err != nil {
			panic(err)
		}
	})
}
