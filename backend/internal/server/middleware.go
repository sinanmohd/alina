package server

import (
	"context"
	"log"
	"net/http"
)

func middlewareCorsOnFlag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if server.cfg.CorsAllowAll {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
		}

		next.ServeHTTP(rw, req)
	})
}

func middlewareCorsAlways(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(rw, req)
	})
}

func middlewareIpLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ipAddr, err := ipFromReq(req)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("Error parsing ip:", err)
			return
		}

		count, err := server.queries.UploadsIpCountPerDay(context.Background(), *ipAddr)
		if err != nil {
			log.Println("Error querying db:", err)
			return
		}
		if int(count) >= server.cfg.UploadsPerDay {
			http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(rw, req)
	})
}
