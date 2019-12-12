package main

import (
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/brunoluiz/go-pwa-server/htmlmod"
	"github.com/brunoluiz/go-pwa-server/js"
	"github.com/brunoluiz/go-pwa-server/middlewares"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Usage: "PWA static file server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "js-env-prefix",
				Usage:  "Dynamic JS env variables prefix",
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
				Usage:  "JS config route",
				Value:  "/__/config.js",
				EnvVar: "JS_ENV_ROUTE",
			},
			&cli.BoolFlag{
				Name:   "no-cache",
				Usage:  "Add no-cache headers",
				EnvVar: "NO_CACHE",
			},
			&cli.BoolFlag{
				Name:   "compression",
				Usage:  "Enable gzip compression",
				EnvVar: "COMPRESSION",
			},
			&cli.BoolFlag{
				Name:   "cors",
				Usage:  "Add CORS Origin, Method and Headers as *",
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
	h := htmlmod.Serve(c.String("dir"), c.String("base-url"))

	if c.Bool("compression") {
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
