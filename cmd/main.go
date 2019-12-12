package main

import (
	"errors"
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
				Name:   "env-js-prefix",
				Usage:  "Dynamic JS env variables prefix",
				Value:  "CONFIG_",
				EnvVar: "ENV_JS_PREFIX",
			},
			&cli.StringFlag{
				Name:   "env-js-key",
				Usage:  "Which key to use when exposing the config to window",
				Value:  "config",
				EnvVar: "ENV_JS_WINDOW_KEY",
			},
			&cli.StringFlag{
				Name:   "env-js-route",
				Usage:  "JS config route",
				EnvVar: "ENV_JS_ROUTE",
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
				Usage:  "If set, adds <base href=value> on HTML heads",
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
	if c.String("dir") == "" {
		return errors.New("No static file directory set")
	}

	var h http.Handler
	if c.String("base-url") != "" {
		logrus.Info("HTML dynamic modification enabled")
		h = htmlmod.Serve(c.String("dir"), c.String("base-url"))
	} else {
		h = http.FileServer(http.Dir(c.String("dir")))
	}

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

	if c.String("env-js-route") != "" {
		logrus.Infof(
			"Env to JS route registered at %s (env %s, window.%s)",
			c.String("env-js-route"),
			c.String("env-js-prefix"),
			c.String("env-js-key"),
		)

		http.Handle(
			c.String("env-js-route"),
			js.Handler(c.String("env-js-prefix"), c.String("env-js-key")),
		)
	}

	return http.ListenAndServe(c.String("address"), nil)
}
