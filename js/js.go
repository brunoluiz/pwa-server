package js

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func RenderGlobalsFromEnv(prefix string, key string) string {
	values := []string{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		value := pair[1]

		logrus.Info(key, prefix)
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

type PwaGlobals struct {
	Prefix string
	Key    string
}

func Handler(prefix string, key string) *PwaGlobals {
	return &PwaGlobals{prefix, key}
}

func (pwa *PwaGlobals) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, RenderGlobalsFromEnv(pwa.Prefix, pwa.Key))
}
