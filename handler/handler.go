package handler

import (
	"fmt"
	"net/http"
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
		logrus.Infof("%s enabled", name)
	}

	return h
}

// Static Exposes static files through HTTP
func Static(dir string, interceptors ...InterceptorConfig) http.Handler {
	return ApplyInterceptors(http.FileServer(http.Dir(dir)), interceptors...)
}

// Ready Disable CORS (add * headers)
func Ready() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
}
