package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"path"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/crypto/blake2b"
	"sinanmohd.com/alina/db"
	"toolman.org/encoding/base56"
)

func uploadSimple(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(0)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("Error parsing form:", err)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("Error fetching file:", err)
		return
	}
	defer file.Close()

	hasher, err := blake2b.New256(nil)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error hashing file:", err)
		return
	}
	_, err = io.Copy(hasher, file)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error hashing file:", err)
		return
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	row, err := server.queries.GetFileFromHash(context.Background(), hash)
	if err == nil {
		fileId56 := base56.Encode(uint64(row.ID))
		fileName := fmt.Sprintf("%v%v", fileId56, mimetype.Lookup(row.MimeType).Extension())

		fmt.Fprintf(rw, "%v/%v\n", server.cfg.PublicUrl, fileName)
		return
	}

	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error detecting mimetype:", err)
		return
	}

	ipAddrString, _, _ := net.SplitHostPort(req.RemoteAddr)
	if value := req.Header.Get("X-Forwarded-For"); value != "" {
		ipAddrString = value
	}
	ipAddr, err := netip.ParseAddr(ipAddrString)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error parsing ip:", err)
		return
	}

	fileId, err := server.queries.CreateFile(context.Background(), db.CreateFileParams{
		MimeType: mtype.String(),
		Name:     header.Filename,
		FileSize: header.Size,
		IpAddr:   ipAddr,
		Hash:     hash,
	})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error querying db:", err)
		return
	}

	fileId56 := base56.Encode(uint64(fileId))
	fileName := fmt.Sprintf("%v%v", fileId56, mtype.Extension())
	filePath := path.Join(server.storagePath, fileName)

	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		server.queries.DeleteFile(context.Background(), fileId)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error creating file:", err)
		return
	}
	defer dst.Close()

	file.Seek(0, 0)
	_, err = io.Copy(dst, file)
	if err != nil {
		server.queries.DeleteFile(context.Background(), fileId)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error copying file:", err)
		return
	}

	fmt.Fprintf(rw, "%v/%v\n", server.cfg.PublicUrl, fileName)
	return
}

func uploadChunkedStart(rw http.ResponseWriter, r *http.Request) {
}

func uploadChunkedProgress(rw http.ResponseWriter, r *http.Request) {
}

func uploadChunkedCancel(rw http.ResponseWriter, r *http.Request) {
}
