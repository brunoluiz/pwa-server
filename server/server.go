package server

import "net/http"

type Interceptor func(http.Handler) http.Handler

func StaticHandler(addr string, interceptors ...Interceptor) (h http.Handler) {
	for _, interceptor := range interceptors {
		h = interceptor(h)
	}

	return h
}
