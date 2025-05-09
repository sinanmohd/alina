package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"sinanmohd.com/alina/db"
	"sinanmohd.com/alina/internal/config"
)

var server struct {
	queries *db.Queries
	cfg     config.ServerConfig

	storagePath string
	chunkedPath string
}

//go:embed  all:frontend
var frontendFs embed.FS

func Run(cfg config.ServerConfig, queries *db.Queries) error {
	mux := http.NewServeMux()
	server.queries = queries
	server.cfg = cfg
	server.storagePath = path.Join(cfg.Data, "storage")
	server.chunkedPath = path.Join(cfg.Data, "chunked")

	err := os.MkdirAll(server.storagePath, 0755)
	if err != nil {
		log.Println("Error creating directory: ", err)
		return err
	}

	if cfg.CorsAllowAll {
		corsOptionsHandler := middlewareCorsOnFlag(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.Header().Set("Access-Control-Allow-Methods", "*")
			rw.Header().Set("Access-Control-Allow-Headers", "*")
		}))

		mux.Handle("OPTIONS /", corsOptionsHandler)
	}

	mux.Handle("GET /metrics", promhttp.Handler())

	mux.HandleFunc("GET /", func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, "/home/", http.StatusMovedPermanently)
	})
	frontend, err := fs.Sub(fs.FS(frontendFs), "frontend")
	if err != nil {
		log.Println("Error traversing fs: ", err)
		return err
	}
	httpFs := http.FileServer(http.FS(frontend))
	mux.Handle("GET /home/", middlewareCorsAlways(http.StripPrefix("/home/", httpFs)))
	mux.Handle("GET /favicon.ico", httpFs)
	mux.Handle("GET /robots.txt", httpFs)

	fs := middlewareCorsOnFlag(http.FileServer(http.Dir(server.storagePath)))
	mux.Handle("GET /{fileId}", fs)
	mux.Handle("GET /files/{fileId}", http.StripPrefix("/files/", fs))

	mux.HandleFunc("GET /notes/{fileId}", notes)
	mux.HandleFunc("GET /notes/styles.css", notesCSS)

	publicConfigHandler := middlewareCorsOnFlag(http.HandlerFunc(publicConfig))
	mux.Handle("GET /_alina/config", publicConfigHandler)

	uploadSimpleHandler := middlewareCorsOnFlag(http.HandlerFunc(uploadSimple))
	mux.Handle("POST /", uploadSimpleHandler)
	mux.Handle("POST /_alina/upload/simple", uploadSimpleHandler)

	uploadChunkedStartHandler := middlewareCorsOnFlag(http.HandlerFunc(uploadChunkedStart))
	uploadChunkedProgressHandler := middlewareCorsOnFlag(http.HandlerFunc(uploadChunkedProgress))
	uploadChunkedCancelHandler := middlewareCorsOnFlag(http.HandlerFunc(uploadChunkedCancel))
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
