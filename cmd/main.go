package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/brunoluiz/go-pwa-server/htmlmod"
	"github.com/brunoluiz/go-pwa-server/js"
	"github.com/brunoluiz/go-pwa-server/middlewares"
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
				Value:  "CONFIG_",
				EnvVar: "JS_ENV_PREFIX",
			},
			&cli.StringFlag{
				Name:   "js-env-key",
				Usage:  "Which key to use when exposing the config to window",
				Value:  "config",
				EnvVar: "JS_ENV_WINDOW_KEY",
			},
			&cli.StringFlag{
				Name:   "js-env-route",
				Usage:  "Where window.config js is exposed",
				Value:  "/__/config.js",
				EnvVar: "JS_ENV_ROUTE",
			},
			&cli.BoolFlag{
				Name:   "no-cache",
				Usage:  "No cache headers",
				EnvVar: "NO_CACHE",
			},
			&cli.BoolFlag{
				Name:   "compress",
				Usage:  "Allow response gzip compression",
				EnvVar: "COMPRESS",
			},
			&cli.BoolFlag{
				Name:   "cors",
				Usage:  "Set CORS Origin, Method and Headers to be *",
				EnvVar: "CORS",
			},
			&cli.StringFlag{
				Name:   "dir",
				Usage:  "Static files directory",
				Value:  ":80",
				EnvVar: "DIR",
			},
			&cli.StringFlag{
				Name:   "address",
				Usage:  "Server address",
				Value:  ":80",
				EnvVar: "ADDRESS",
			},
			&cli.StringFlag{
				Name:   "base-url",
				Usage:  ".",
				EnvVar: "BASE_URL",
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

	var h http.Handler
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	if c.Bool("compress") {
		logrus.Info("Compression enabled")
		h = gziphandler.GzipHandler(h)
	}

	if c.Bool("no-cache") {
		logrus.Info("No-cache headers enabled")
		h = middlewares.NoCache(h)
	}

	if c.Bool("cors") {
		logrus.Info("CORS headers enabled")
		h = middlewares.Cors(h)
	}

	http.Handle("/", h)

	logrus.Infof(
		"Loading js script at %s, using env variables with prefix %s and exposing as window.%s",
		c.String("js-env-route"),
		c.String("js-env-prefix"),
		c.String("js-env-key"),
	)
	http.Handle(c.String("js-env-route"), js.Serve(c.String("js-env-prefix"), c.String("js-env-key")))

	return http.ListenAndServe(c.String("address"), nil)
}
