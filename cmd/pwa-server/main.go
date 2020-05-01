package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/brunoluiz/pwa-server/envjs"
	"github.com/brunoluiz/pwa-server/handler"
	"github.com/brunoluiz/pwa-server/middleware"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
				Name:   "ready-route",
				Usage:  "Ready probe route",
				Value:  "/__/ready",
				EnvVar: "READY_ROUTE",
			},
			&cli.StringFlag{
				Name:   "metrics-route",
				Usage:  "Metrics route",
				Value:  "/__/metrics",
				EnvVar: "METRICS_ROUTE",
			},
			&cli.StringFlag{
				Name:   "env-js-route",
				Usage:  "JS config route",
				EnvVar: "ENV_JS_ROUTE",
				Value:  "/__/config.js",
			},
			&cli.StringFlag{
				Name:   "env-js-prefix",
				Usage:  "Dynamic JS env variables prefix (eg: window.env.CONFIG_SOMEFLAG = true)",
				Value:  "CONFIG_",
				EnvVar: "ENV_JS_PREFIX",
			},
			&cli.StringFlag{
				Name:   "env-js-window-key",
				Usage:  "Which key to use when exposing the config to window (eg: window.env)",
				Value:  "env",
				EnvVar: "ENV_JS_WINDOW_KEY",
			},
			&cli.StringFlag{
				Name:   "base-url",
				Usage:  "If set, adds <base href=value> on HTML heads",
				EnvVar: "BASE_URL",
			},
			&cli.BoolFlag{
				Name:   "no-manifest-base-url",
				Usage:  "Disables base-url manipulations for manifest.json",
				EnvVar: "NO_MANIFEST_BASE_URL",
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
			&cli.StringFlag{
				Name:   "not-found-file",
				Usage:  "Redirect request to specific file if nothing was found on the route (index.html enables HTML Push State)",
				EnvVar: "NOT_FOUND_FILE",
				Value:  "index.html",
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
			&cli.BoolFlag{
				Name:   "req-logger",
				Usage:  "Enable request logger",
				EnvVar: "REQ_LOGGER",
			},
			&cli.StringFlag{
				Name:   "req-logger-format",
				Usage:  "Request logger format (apache)",
				EnvVar: "REQ_LOGGER_FORMAT",
				Value:  middleware.LogFormatCommon,
			},
			&cli.BoolFlag{
				Name:   "debug",
				Usage:  "Turn on debug mode",
				EnvVar: "DEBUG",
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

	if !c.Bool("debug") {
		logrus.SetLevel(logrus.FatalLevel)
	}

	// Start HTTP mux
	mux := chi.NewRouter()

	// Operational handlers
	mux.Handle(c.String("ready-route"), handler.Ready())
	mux.Handle(c.String("metrics-route"), promhttp.Handler())

	// Serve static files
	mux.Handle("/*", handler.Static(dir, c.String("not-found-file")))

	// Create JS config route
	logrus.Infof(
		"JS config at %s (env %s, window.%s)",
		c.String("env-js-route"), c.String("env-js-prefix"), c.String("env-js-key"),
	)

	js := envjs.Handler(c.String("env-js-prefix"), c.String("env-js-key"))
	mux.Handle(c.String("env-js-route"), js)

	// Create interceptors
	interceptors := []handler.InterceptorConfig{
		{Wrapper: middleware.HTMLBaseURL(dir, c.String("base-url")), Disable: c.String("base-url") == ""},
		{Wrapper: middleware.ManifestBaseURL(dir, c.String("base-url")), Disable: c.Bool("no-manifest-base-url")},
		{Wrapper: gziphandler.GzipHandler, Disable: c.Bool("no-compression")},
		{Wrapper: middleware.NoCache, Disable: c.Bool("allow-cache")},
		{Wrapper: middleware.Cors, Disable: !c.Bool("cors")},
		{Wrapper: middleware.Helmet, Disable: c.Bool("no-helmet")},
		{Wrapper: middleware.Logger(c.String("req-logger-format")), Disable: !c.Bool("req-logger")},
	}

	// Serve
	return http.ListenAndServe(
		c.String("address"),
		handler.ApplyInterceptors(mux, interceptors...),
	)
}
