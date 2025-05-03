package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/blake2b"
	"sinanmohd.com/alina/db"
	"toolman.org/encoding/base56"
)

type ChunkedStartReq struct {
	FileSize int    `json:"file_size" validate:"required,numeric,gt=1"`
	Name     string `json:"name" validate:"required"`
}

type ChunkedStartResp struct {
	ChunkToken string `json:"chunk_token"`
}

type ChunkedJwtClaims struct {
	ChunkedId int32 `json:"chunk_id"`
	jwt.RegisteredClaims
}

func uploadChunkedStart(rw http.ResponseWriter, req *http.Request) {
	var data ChunkedStartReq

	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if data.FileSize > server.cfg.FileSizeLimit {
		http.Error(rw, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	ipAddr, err := ipFromReq(req)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error parsing ip:", err)
		return
	}

	chunkCount := int32(math.Ceil(float64(data.FileSize) / float64(server.cfg.ChunkSize)))
	chunkedId, err := server.queries.ChunkedCreate(context.Background(), db.ChunkedCreateParams{
		FileSize:   int64(data.FileSize),
		Name:       data.Name,
		IpAddr:     *ipAddr,
		ChunksLeft: chunkCount,
	})

	claims := ChunkedJwtClaims{
		ChunkedId: chunkedId,
	}
	chunkToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(server.cfg.SecretKey))
	if err != nil {
		server.queries.ChunkedDelete(context.Background(), chunkedId)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error creating token:", err)
		return
	}

	resp, err := json.Marshal(ChunkedStartResp{
		ChunkToken: chunkToken,
	})
	if err != nil {
		server.queries.ChunkedDelete(context.Background(), chunkedId)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error marshaling response:", err)
		return
	}

	fmt.Fprint(rw, string(resp))
	return
}

func uploadChunkedProgress(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(0)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var claims ChunkedJwtClaims
	_, err = jwt.ParseWithClaims(req.FormValue("chunk_token"), claims, func(token *jwt.Token) (any, error) {
		return server.cfg.SecretKey, nil
	})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	row, err := server.queries.ChunkedFromId(context.Background(), claims.ChunkedId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(rw, http.StatusText(http.StatusGone), http.StatusGone)
			return
		} else {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("Error querying db:", err)
			return
		}
	}
	if time.Now().Sub(row.CreatedAt.Time).Hours() >= 24*3 ||
		time.Now().Sub(row.LastAccess.Time).Hours() >= 24*1 {
		http.Error(rw, http.StatusText(http.StatusGone), http.StatusGone)
		return
	}

	chunkDir := path.Join(server.chunkedPath, fmt.Sprint(claims.ChunkedId))
	err = os.MkdirAll(chunkDir, 0700)
	if err != nil {
		log.Println("Error creating directory: ", err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	chunkIndex, err := strconv.Atoi(req.FormValue("chunk_index"))
	if err != nil || int32(chunkIndex) > row.ChunksTotal || chunkIndex < 1 {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	chunk, header, err := req.FormFile("file")
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer chunk.Close()
	if int32(chunkIndex) == row.ChunksTotal {
		if header.Size != int64(server.cfg.FileSizeLimit)%row.FileSize {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else {
		if header.Size != int64(server.cfg.ChunkSize) {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	chunkPath := path.Join(server.chunkedPath, fmt.Sprint(chunkIndex))
	_, err = os.Stat(chunkPath)
	if err == nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else if !errors.Is(err, os.ErrNotExist) {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error stating file:", err)
		return
	}

	dst, err := os.OpenFile(chunkPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error creating file:", err)
		return
	}
	_, err = io.Copy(dst, chunk)
	dst.Close()
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error copying file:", err)
		os.Remove(chunkPath)
		return
	}

	chunksLeft, err := server.queries.ChunkedLeftDecrement(context.Background(), claims.ChunkedId)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error querying db:", err)
		os.Remove(chunkPath)
		return
	}
	if chunksLeft == 0 {

	}

	return
}

func ChunkedToFile(chunkedId int32, chunksTotal int, rw http.ResponseWriter) error {

	chunkDir := path.Join(server.chunkedPath, fmt.Sprint(chunkedId))
	chunkFullFilePath := path.Join(chunkDir, "fullfile")
	fullFile, err := os.OpenFile(chunkFullFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("Error creating file:", err)
		return err
	}

	for chunkIndex := 1; chunkIndex <= chunksTotal; chunkIndex++ {
		chunkPath := path.Join(server.chunkedPath, fmt.Sprint(chunkIndex))
		chunk, err := os.Open(chunkPath)
		if err != nil {
			log.Println("Error creating file:", err)
			fileCloseAndRemove(fullFile, chunkFullFilePath)
			fullFile.Close()
			os.Remove(chunkFullFilePath)
			return err
		}
		_, err = io.Copy(fullFile, chunk)
		if err != nil {
			log.Println("Error copying file:", err)
			fileCloseAndRemove(fullFile, chunkFullFilePath)
			return err
		}
		chunk.Close()
	}

	hasher, err := blake2b.New256(nil)
	if err != nil {
		log.Println("Error hashing file:", err)
		fileCloseAndRemove(fullFile, chunkFullFilePath)
		return err
	}
	_, err = io.Copy(hasher, fullFile)
	if err != nil {
		log.Println("Error hashing file:", err)
		fileCloseAndRemove(fullFile, chunkFullFilePath)
		return err
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	row, err := server.queries.FileFromHash(context.Background(), hash)
	if err == nil {
		fileCloseAndRemove(fullFile, chunkDir)
		fileId56 := base56.Encode(uint64(row.ID))
		fileName := fmt.Sprintf("%v%v", fileId56, mimetype.Lookup(row.MimeType).Extension())

		fmt.Fprintf(rw, "%v/%v\n", server.cfg.PublicUrl, fileName)
		return nil
	}

	fullFile.Seek(0, 0)
	mtype, err := mimetype.DetectReader(fullFile)
	if err != nil {
		log.Println("Error detecting mimetype:", err)
		return err
	}
	fileId, err := server.queries.ChunkedToFile(context.Background(), db.ChunkedToFileParams{
		MimeType: mtype.String(),
		Hash:     hash,
		ID:       chunkedId,
	})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error querying db:", err)
		return err
	}
	fileId56 := base56.Encode(uint64(fileId))
	fileName := fmt.Sprintf("%v%v", fileId56, mimetype.Lookup(row.MimeType).Extension())
	filePath := path.Join(server.storagePath, fileName)
	os.Rename(chunkFullFilePath, filePath)

	fmt.Fprintf(rw, "%v/%v\n", server.cfg.PublicUrl, fileName)
	return nil
}

func fileCloseAndRemove(file *os.File, path string) {
	err := file.Close()
	if err != nil {
		log.Println("Error closing file:", err)
	}

	err = os.RemoveAll(path)
	if err != nil {
		log.Println("Error deleting file:", err)
	}
}

func uploadChunkedCancel(rw http.ResponseWriter, r *http.Request) {
}
