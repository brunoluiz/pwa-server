package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/brunoluiz/go-pwa-server/envjs"
	"github.com/brunoluiz/go-pwa-server/handler"
	"github.com/brunoluiz/go-pwa-server/middleware"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Usage: "PWA static file server",
		Flags: []cli.Flag{
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
				Name:   "allow-cache",
				Usage:  "Disable no-cache headers",
				EnvVar: "ALLOW_CACHE",
			},
			&cli.BoolFlag{
				Name:   "no-compression",
				Usage:  "Enable gzip compression",
				EnvVar: "NO_COMPRESSION",
			},
			&cli.BoolFlag{
				Name:   "cors",
				Usage:  "Add CORS Origin, Method and Headers as *",
				EnvVar: "CORS",
			},
			&cli.BoolFlag{
				Name:   "no-helmet",
				Usage:  "Disable security headers (helmet)",
				EnvVar: "NO_HELMET",
			},
		},
		Action: serve,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func serve(c *cli.Context) error {
	dir := c.String("dir")
	if dir == "" {
		return errors.New("no static file directory set")
	}

	mux := http.NewServeMux()

	interceptors := []handler.InterceptorConfig{
		{"HTML Dynamic", middleware.HTMLBaseURL(dir, c.String("base-url")), c.String("base-url") == ""},
		{"Compression", gziphandler.GzipHandler, c.Bool("no-compression")},
		{"No-cache", middleware.NoCache, c.Bool("allow-cache")},
		{"CORS", middleware.Cors, c.Bool("cors")},
		{"Helmet", middleware.Helmet, c.Bool("no-helmet")},
	}
	mux.Handle("/", handler.Static(dir, interceptors...))

	if c.String("env-js-route") != "" {
		logrus.Infof(
			"JS config at %s (env %s, window.%s)",
			c.String("env-js-route"), c.String("env-js-prefix"), c.String("env-js-key"),
		)

		js := envjs.Handler(c.String("env-js-prefix"), c.String("env-js-key"))
		js = middleware.Helmet(js)
		mux.Handle(c.String("env-js-route"), js)
	}

	return http.ListenAndServe(c.String("address"), mux)
}
