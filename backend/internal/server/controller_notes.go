package server

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"toolman.org/encoding/base56"
)

//go:embed assets/simple.min.css
var css []byte

func notes(rw http.ResponseWriter, req *http.Request) {
	fileId56 := req.PathValue("fileId")
	fileId, err := base56.Decode(fileId56)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	row, err := server.queries.FileFromId(context.Background(), int32(fileId))
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else if row.FileSize > 1024*1024 {
		http.Error(rw, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	} else if !strings.HasPrefix(row.MimeType, "text/") {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mtype := mimetype.Lookup(row.MimeType)
	if mtype == nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error detecting mimetype:", err)
		return
	}

	fileName := fmt.Sprintf("%v%v", fileId56, mtype.Extension())
	filePath := path.Join(server.storagePath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error reading file:", err)
		return
	}

	unsafeHTML := markdown.ToHTML(fileContent, nil, nil)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)

	htmlFmt := `
		<head>
			<title>Note %v</title>
			<link rel="stylesheet" href="/notes/styles.css">
			<meta charset="utf-8">

			<meta name="author" content="alina">
			<meta name="description" content="Your frenly neighbourhood file sharing website">

			<meta property="og:title" content="Alina">
			<meta property="og:description" content="Your frenly neighbourhood file sharing website">
		</head>
		<body>%v</body>
	`
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = fmt.Fprintf(rw, htmlFmt, fileId56, string(html))
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error writing html:", err)
		return
	}

	return
}

func notesCSS(rw http.ResponseWriter, req *http.Request) {
	_, err := rw.Write(css)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error writing css:", err)
		return
	}

	return
}
