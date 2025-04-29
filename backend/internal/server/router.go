package server

import "net/http"

func Run() {
	http.HandleFunc("POST /", uploadSimple)
	http.HandleFunc("POST /_alina/upload/simple", uploadSimple)

	http.HandleFunc("PUT /_alina/upload/chunked", uploadChunkedStart)
	http.HandleFunc("POST /_alina/upload/chunked", uploadChunkedProgress)
	http.HandleFunc("PATCH /_alina/upload/chunked", uploadChunkedEnd)
	http.HandleFunc("DELETE /_alina/upload/chunked", uploadChunkedCancel)

	http.ListenAndServe(":8008", nil)
}
