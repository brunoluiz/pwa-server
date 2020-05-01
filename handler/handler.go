package handler

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Interceptor wraps an http.Handler with intercepting functionality
type Interceptor func(http.Handler) http.Handler

// InterceptorConfig configs for interceptor
type InterceptorConfig struct {
	Wrapper Interceptor
	Disable bool
}

// ApplyInterceptors apply based on InterceptorConfig
func ApplyInterceptors(h http.Handler, interceptors ...InterceptorConfig) http.Handler {
	for _, interceptor := range interceptors {
		if interceptor.Disable {
			continue
		}

		name := runtime.FuncForPC(reflect.ValueOf(interceptor.Wrapper).Pointer()).Name()
		h = interceptor.Wrapper(h)
		logrus.Infof("enabled: %s", name)
	}

	return h
}

// Static Exposes static files through HTTP
func Static(dir string, interceptors ...InterceptorConfig) http.Handler {
	root := http.Dir(dir)
	fs := http.StripPrefix("/", http.FileServer(root))

	return ApplyInterceptors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(fmt.Sprintf("%s", root) + r.RequestURI); os.IsNotExist(err) {
			http.ServeFile(w, r, dir+"/index.html")
			return
		}

		fs.ServeHTTP(w, r)
	}), interceptors...)
}

// Ready Disable CORS (add * headers)
func Ready() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
}
