package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

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

		logrus.Infof("%s enabled", interceptor.Name)
		h = interceptor.Handler(h)
	}

	return h
}

func Static(dir string, interceptors ...InterceptorConfig) http.Handler {
	return ApplyInterceptors(http.FileServer(http.Dir(dir)), interceptors...)
}
