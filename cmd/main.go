package main

import (
	"net/http"
	"os"
	"strings"

	htmlmod "github.com/brunoluiz/go-pwa-server/htmlmod/html.go"
	"github.com/brunoluiz/go-pwa-server/js"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/net/html"
)

func main() {
	app := &cli.App{
		Usage: "PWA static file server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "js-env-prefix",
				Usage:  "Value to get JS env variables",
				Value:  "CONFIG",
				EnvVar: "SERVER_JS_ENV_PREFIX",
			},
			&cli.StringFlag{
				Name:   "js-env-window-key",
				Usage:  "Which key to use when exposing the config to window",
				Value:  "config",
				EnvVar: "SERVER_JS_ENV_WINDOW_KEY",
			},
			&cli.StringFlag{
				Name:   "js-env-route",
				Usage:  "Where window.config js is exposed",
				Value:  "/__/config.js",
				EnvVar: "SERVER_JS_ENV_ROUTE",
			},
			&cli.StringFlag{
				Name:   "dir",
				Usage:  "Static files directory",
				Value:  ":80",
				EnvVar: "SERVER_DIR",
			},
			&cli.StringFlag{
				Name:   "address",
				Usage:  "Server address",
				Value:  ":80",
				EnvVar: "SERVER_ADDRESS",
			},
			&cli.StringFlag{
				Name:   "base-url",
				Usage:  ".",
				EnvVar: "SERVER_BASE_URL",
			},
		},
		Action: serve,
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func serve(c *cli.Context) error {
	fs := http.FileServer(http.Dir(c.String("dir")))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hasTrailing := r.URL.Path[len(r.URL.Path)-1:] == "/"
		if !strings.Contains(r.URL.Path, "htm") && !hasTrailing {
			fs.ServeHTTP(w, r)
			return
		}

		path := c.String("dir") + r.URL.Path
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

		htmlmod.EnhanceHTML(doc, &htmlmod.PWASettings{
			BaseURL: c.String("base-url"),
		})
		if err := html.Render(w, doc); err != nil {
			panic(err)
		}
	})

	http.Handle(c.String("js-env-route"), js.Serve(c.String("js-env-prefix"), c.String("js-env-window-key")))

	return http.ListenAndServe(c.String("address"), nil)
}
