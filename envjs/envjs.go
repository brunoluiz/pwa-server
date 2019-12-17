package envjs

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// EnvToJS create which can be used by applications to access server specified env variables
// Example:
// - prefix: CONFIG_
// - key: app
// - Env variables: PATH=/bin/bash,OS=linux,CONFIG_FOO=bar
// - Output: window.app={FOO:"bar"}
func EnvToJS(prefix string, key string) string {
	values := []string{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		value := pair[1]

		if !strings.HasPrefix(key, prefix) {
			continue
		}

		if value == "true" || value == "false" {
			values = append(values, key+":"+value)
			continue
		}

		values = append(values, key+":\""+value+"\"")
	}

	return "window." + key + "={" + strings.Join(values, ",") + "}"
}

// Handler Outputs EnvToJS string
func Handler(prefix string, key string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprint(w, EnvToJS(prefix, key))
	})
}
