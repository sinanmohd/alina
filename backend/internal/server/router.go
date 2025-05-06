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
	mux := http.NewServeMux()
	server.queries = queries
	server.cfg = cfg
	server.storagePath = path.Join(cfg.Data, "storage")
	server.chunkedPath = path.Join(cfg.Data, "chunked")

	err := os.MkdirAll(server.storagePath, 0700)
	if err != nil {
		log.Println("Error creating directory: ", err)
		return err
	}

	if cfg.CorsAllowAll {
		corsOptionsHandler := middleware(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.Header().Set("Access-Control-Allow-Methods", "*")
			rw.Header().Set("Access-Control-Allow-Headers", "*")
		}))

		mux.Handle("OPTIONS /", corsOptionsHandler)
	}

	fs := middleware(http.FileServer(http.Dir(server.storagePath)))
	mux.Handle("GET /", fs)
	mux.Handle("GET /files/", http.StripPrefix("/files/", fs))

	publicConfigHandler := middleware(http.HandlerFunc(publicConfig))
	mux.Handle("GET /_alina/config", publicConfigHandler)

	uploadSimpleHandler := middleware(http.HandlerFunc(uploadSimple))
	mux.Handle("POST /", uploadSimpleHandler)
	mux.Handle("POST /_alina/upload/simple", uploadSimpleHandler)

	uploadChunkedStartHandler := middleware(http.HandlerFunc(uploadChunkedStart))
	uploadChunkedProgressHandler := middleware(http.HandlerFunc(uploadChunkedProgress))
	uploadChunkedCancelHandler := middleware(http.HandlerFunc(uploadChunkedCancel))
	mux.Handle("POST /_alina/upload/chunked", uploadChunkedStartHandler)
	mux.Handle("PATCH /_alina/upload/chunked", uploadChunkedProgressHandler)
	mux.Handle("DELETE /_alina/upload/chunked", uploadChunkedCancelHandler)

	bindAddr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	log.Printf("alina is listening on http://%v\n", bindAddr)
	err = http.ListenAndServe(bindAddr, mux)
	if err != nil {
		log.Println("Error serving http: ", err)
		return err
	}

	return nil
}
