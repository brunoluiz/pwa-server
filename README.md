# pwa-server

🚀 Fast static server for your PWA applications

## Features

- Adds security headers (helmet) **by default**
- Enable compression **by default**
- Disable cache **by default**
- Exposes /__/ready route **by default**
- Exposes /__/metrics route **by default**
- Supports not found redirects -- default is `index.html`, to support HTML5 PushState
- Supports adding a base url as `<base href={base-url} />` in `<head>`, if specified
- Supports adding a base url to `manifest.json`, if specified
- Support allowing CORS *
- Supports exposing env variables with a configured prefix in a specific route
- Supports request logging

## Install

### MacOS

Use `brew` to install it

```
brew tap brunoluiz/tap
brew install pwa-server
```

### Linux and Windows

[Check the releases section](https://github.com/brunoluiz/pwa-server/releases) for more information details 

### go get

Install using `GO111MODULES=off go get github.com/brunoluiz/pwa-server/cmd/pwa-server` to get the latest version. This will place it in your `$GOPATH`, enabling it to be used anywhere in the system.

**⚠️ Reminder**: the command above download the contents of master, which might not be the stable version. [Check the releases](https://github.com/brunoluiz/pwa-server/releases) and get a specific tag for stable versions.

### Docker

The tool is available as a Docker image as well. Please refer to [Docker Hub page](https://hub.docker.com/r/brunoluiz/pwa-server/tags) to pick a release

```
docker run -p 80:80 \
  --env-file .env.sample \
  -v $(PWD)/test/static:/static \
  brunoluiz/pwa-server
```

## Usage

Run `pwa-server --dir ./your/static/dir` with the following options (all can be set as ENV variables as well)

```
   --dir value                Static files directory [$DIR]
   --address value            Server address (default: ":80") [$ADDRESS]
   --ready-route value        Ready probe route (default: "/__/ready") [$READY_ROUTE]
   --metrics-route value      Metrics route (default: "/__/metrics") [$METRICS_ROUTE]
   --env-js-route value       JS config route (default: "/__/config.js") [$ENV_JS_ROUTE]
   --env-js-prefix value      Dynamic JS env variables prefix (default: "CONFIG_") [$ENV_JS_PREFIX]
   --env-js-key value         Which key to use when exposing the config to window (default: "config") [$ENV_JS_WINDOW_KEY]
   --base-url value           If set, adds <base href=value> on HTML heads [$BASE_URL]
   --no-manifest-base-url     Disables base-url manipulations for manifest.json [$NO_MANIFEST_BASE_URL]
   --allow-cache              Disable no-cache headers [$ALLOW_CACHE]
   --no-compression           Enable gzip compression [$NO_COMPRESSION]
   --not-found-file value     Redirect request to specific file if nothing was found on the route (index.html enables HTML Push State) (default: "index.html") [$NOT_FOUND_FILE]
   --cors                     Add CORS Origin, Method and Headers as * [$CORS]
   --no-helmet                Disable security headers (helmet) [$NO_HELMET]
   --req-logger               Enable request logger [$REQ_LOGGER]
   --req-logger-format value  Request logger format (apache) (default: "%h %l %u %t \"%r\" %>s %b") [$REQ_LOGGER_FORMAT]
   --debug                    Turn on debug mode [$DEBUG]
   --help, -h                 show help
   --version, -v              print the version
```

### Env to JS

Exposes enviroment variables in a specific route.

Example: `pwa-server --env-js-prefix CONFIG_ --env-js-key app --env-js-route /config.js` loads:
- Loads all env variables prefixed as `CONFIG_`
- Exposes as `window.app={CONFIG_...="some-value"}`
- At `http://0.0.0.0/config.js`
