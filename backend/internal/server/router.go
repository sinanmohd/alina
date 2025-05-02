package server

import (
	"fmt"
	"log"
	"net/http"

	"sinanmohd.com/alina/db"
	"sinanmohd.com/alina/internal/config"
)

var queries *db.Queries

func Run(cfg config.ServerConfig, q *db.Queries) {
	queries = q

	http.HandleFunc("POST /", uploadSimple)
	http.HandleFunc("POST /_alina/upload/simple", uploadSimple)

	http.HandleFunc("PUT /_alina/upload/chunked", uploadChunkedStart)
	http.HandleFunc("POST /_alina/upload/chunked", uploadChunkedProgress)
	http.HandleFunc("DELETE /_alina/upload/chunked", uploadChunkedCancel)

	bindAddr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	log.Printf("alina is listening on http://%v\n", bindAddr)
	http.ListenAndServe(bindAddr, nil)
}
