package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// InterceptorConfig configs for interceptor
type InterceptorConfig struct {
	Name    string
	Handler func(http.Handler) http.Handler
	Disable bool
}

func ApplyInterceptors(h http.Handler, interceptors ...InterceptorConfig) http.Handler {
	for _, interceptor := range interceptors {
		if interceptor.Disable {
			continue
		}

		h = interceptor.Handler(h)
		logrus.Infof("%s enabled", interceptor.Name)
	}

	return h
}

func Static(dir string, interceptors ...InterceptorConfig) http.Handler {
	return ApplyInterceptors(http.FileServer(http.Dir(dir)), interceptors...)
}
