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

type NotFoundRedirectRespWr struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
}

func (w *NotFoundRedirectRespWr) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *NotFoundRedirectRespWr) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}
	return len(p), nil // Lie that we successfully written it
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
