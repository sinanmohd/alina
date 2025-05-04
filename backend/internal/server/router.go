package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"sinanmohd.com/alina/db"
	"sinanmohd.com/alina/internal/config"
)

var server struct {
	queries *db.Queries
	cfg     config.ServerConfig

	storagePath string
	chunkedPath string
}

func Run(cfg config.ServerConfig, queries *db.Queries) error {
	server.queries = queries
	server.cfg = cfg
	server.storagePath = path.Join(cfg.Data, "storage")
	server.chunkedPath = path.Join(cfg.Data, "chunked")

	err := os.MkdirAll(server.storagePath, 0700)
	if err != nil {
		log.Println("Error creating directory: ", err)
		return err
	}

	fs := http.FileServer(http.Dir(server.storagePath))
	http.Handle("GET /", fs)

	http.HandleFunc("POST /", uploadSimple)
	http.HandleFunc("POST /_alina/upload/simple", uploadSimple)

	http.HandleFunc("PUT /_alina/upload/chunked", uploadChunkedStart)
	http.HandleFunc("PATCH /_alina/upload/chunked", uploadChunkedProgress)
	http.HandleFunc("DELETE /_alina/upload/chunked", uploadChunkedCancel)

	bindAddr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	log.Printf("alina is listening on http://%v\n", bindAddr)
	http.ListenAndServe(bindAddr, nil)

	return nil
}
