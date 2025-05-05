package server

import "net/http"

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if server.cfg.CorsAllowAll {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
		}

		next.ServeHTTP(rw, req)
	})
}
