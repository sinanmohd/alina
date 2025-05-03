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
	partialPath string
}

func Run(cfg config.ServerConfig, queries *db.Queries) error {
	server.queries = queries
	server.cfg = cfg
	server.storagePath = path.Join(cfg.Data, "storage")
	server.partialPath = path.Join(cfg.Data, "partial")

	err := os.MkdirAll(server.storagePath, 0700)
	if err != nil {
		log.Println("Error creating directory: ", err)
		return err
	}
	err = os.MkdirAll(server.partialPath, 0700)
	if err != nil {
		log.Println("Error creating directory: ", err)
		return err
	}

	http.HandleFunc("POST /", uploadSimple)
	http.HandleFunc("POST /_alina/upload/simple", uploadSimple)

	http.HandleFunc("PUT /_alina/upload/chunked", uploadChunkedStart)
	http.HandleFunc("POST /_alina/upload/chunked", uploadChunkedProgress)
	http.HandleFunc("DELETE /_alina/upload/chunked", uploadChunkedCancel)

	bindAddr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	log.Printf("alina is listening on http://%v\n", bindAddr)
	http.ListenAndServe(bindAddr, nil)

	return nil
}
