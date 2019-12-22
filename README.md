# go-pwa-server

🚀 Fast static server for your PWA applications

## Features

- Adds security headers (helmet) **by default**
- Enable compression **by default**
- Disable cache **by default**
- Exposes /__/ready route **by default**
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
brew install go-pwa-server
```

### Linux and Windows

[Check the releases section](https://github.com/brunoluiz/go-pwa-server/releases) for more information details 

### go get

Please avoid using `go get`. [Check our releases](https://github.com/brunoluiz/go-pwa-server/releases) for more details.

### Docker

The tool is available as a Docker image as well. Please refer to [Docker Hub page](https://hub.docker.com/r/brunoluiz/go-pwa-server/tags) to pick a release

```
docker run -p 80:80 \
  --env-file .env.sample \
  -v $(PWD)/test/static:/static \
  brunoluiz/go-pwa-server
```

## Usage

Run `go-pwa-server --dir ./your/static/dir` with the following options (all can be set as ENV variables as well)

```
   --dir value                Static files directory [$DIR]
   --address value            Server address (default: ":80") [$ADDRESS]
   --ready-route value        Ready probe route (default: "/__/ready") [$READY_ROUTE]
   --base-url value           If set, adds <base href=value> on HTML heads [$BASE_URL]
   --no-manifest-base-url     Disables base-url manipulations for manifest.json [$NO_MANIFEST_BASE_URL]
   --env-js-prefix value      Dynamic JS env variables prefix (default: "CONFIG_") [$ENV_JS_PREFIX]
   --env-js-key value         Which key to use when exposing the config to window (default: "config") [$ENV_JS_WINDOW_KEY]
   --env-js-route value       JS config route [$ENV_JS_ROUTE]
   --allow-cache              Disable no-cache headers [$ALLOW_CACHE]
   --no-compression           Enable gzip compression [$NO_COMPRESSION]
   --cors                     Add CORS Origin, Method and Headers as * [$CORS]
   --no-helmet                Disable security headers (helmet) [$NO_HELMET]
   --req-logger               Enable request logger [$REQ_LOGGER]
   --req-logger-format value  Request logger format (apache) (default: "%h %l %u %t \"%r\" %>s %b") [$REQ_LOGGER_FORMAT]
   --help, -h                 show help
   --version, -v              print the version
```

### Env to JS

Exposes enviroment variables in a specific route.

Example: `go-pwa-server --env-js-prefix CONFIG_ --env-js-key app --env-js-route /config.js` loads:
- Loads all env variables prefixed as `CONFIG_`
- Exposes as `window.app={CONFIG_...="some-value"}`
- At `http://0.0.0.0/config.js`
