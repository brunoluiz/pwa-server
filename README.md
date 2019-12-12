# go-pwa-server

Fast static server for your PWA applications

## Usage

Run `go-pwa-server --dir ./your/static/dir` with the following options (all can be set as ENV variables as well)

```
   --dir value            Static files directory [$DIR]
   --address value        Server address (default: ":80") [$ADDRESS]
   --base-url value       If set, adds <base href=value> on HTML heads [$BASE_URL]
   --env-js-prefix value  Dynamic JS env variables prefix (default: "CONFIG_") [$ENV_JS_PREFIX]
   --env-js-key value     Which key to use when exposing the config to window (default: "config") [$ENV_JS_WINDOW_KEY]
   --env-js-route value   JS config route [$ENV_JS_ROUTE]
   --allow-cache          Disable no-cache headers [$ALLOW_CACHE]
   --no-compression       Enable gzip compression [$NO_COMPRESSION]
   --cors                 Add CORS Origin, Method and Headers as * [$CORS]
   --no-helmet            Disable security headers (helmet) [$NO_HELMET]
```

## Features

- Adds security headers (helmet) **by default**
- Enable compression **by default**
- Disable cache **by default**
- Supports adding a base url as `<base href={base-url} />` in `<head>`, if specified
- Support allowing CORS *
- Supports exposing env variables with a configured prefix in a specific route

### Env to JS

Exposes enviroment variables in a specific route.

Example: `go-pwa-server --env-js-prefix CONFIG_ --env-js-key app --env-js-route /config.js` loads:
- Loads all env variables prefixed as `CONFIG_`
- Exposes as `window.app={CONFIG_...="some-value"}`
- At `http://0.0.0.0/config.js`

## To-do

- Allow extension
- Docker image
- Auto-publishing
- Tests
