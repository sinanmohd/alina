package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/crypto/blake2b"
	"sinanmohd.com/alina/db"
	"toolman.org/encoding/base56"
)

func uploadSimple(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(0)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if int(header.Size) > server.cfg.FileSizeLimit {
		http.Error(rw, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

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

	txCtx := context.Background()
	tx, err := server.pool.Begin(txCtx)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error creating transcation:", err)
		return
	}
	defer tx.Rollback(txCtx)
	qtx := server.queries.WithTx(tx)

	ipAddr, err := ipFromReq(req)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error parsing ip:", err)
		return
	}

	userAgent := req.Header.Get("User-Agent")
	if userAgent == "" {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userAgentId, err := qtx.UserAgentIdGet(context.Background(), userAgent)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error querying db:", err)
		return
	}

	row, err := server.queries.FileFromHash(context.Background(), hash)
	if err == nil {
		fileId56 := base56.Encode(uint64(row.ID))
		fileName := fmt.Sprintf("%v%v", fileId56, mimetype.Lookup(row.MimeType).Extension())

		_, err := qtx.UploadCreate(context.Background(), db.UploadCreateParams{
			Name:      header.Filename,
			UserAgent: userAgentId,
			File:      int64(row.ID),
			IpAddr:    *ipAddr,
		})
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("Error querying db:", err)
			return
		}

		tx.Commit(txCtx)
		fmt.Fprintf(rw, "%v/%v\n", server.cfg.PublicUrl, fileName)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error seeking file:", err)
		return
	}
	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error detecting mimetype:", err)
		return
	}

	fileId, err := qtx.FileCreate(context.Background(), db.FileCreateParams{
		MimeType: strings.Split(mtype.String(), ";")[0],
		FileSize: header.Size,
		Hash:     hash,
	})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error querying db:", err)
		return
	}
	_, err = qtx.UploadCreate(context.Background(), db.UploadCreateParams{
		Name:      header.Filename,
		IpAddr:    *ipAddr,
		UserAgent: userAgentId,
		File:      int64(fileId),
	})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error querying db:", err)
		return
	}

	fileId56 := base56.Encode(uint64(fileId))
	fileName := fmt.Sprintf("%v%v", fileId56, mtype.Extension())
	filePath := path.Join(server.storagePath, fileName)

	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error creating file:", err)
		return
	}
	defer dst.Close()

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error seeking file:", err)
		return
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error copying file:", err)
		return
	}

	tx.Commit(txCtx)
	fmt.Fprintf(rw, "%v/%v\n", server.cfg.PublicUrl, fileName)
	return
}
