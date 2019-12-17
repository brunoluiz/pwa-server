package handler

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/sirupsen/logrus"
)

// InterceptorConfig configs for interceptor
type InterceptorConfig struct {
	Wrapper func(http.Handler) http.Handler
	Disable bool
}

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

func Static(dir string, interceptors ...InterceptorConfig) http.Handler {
	return ApplyInterceptors(http.FileServer(http.Dir(dir)), interceptors...)
}
