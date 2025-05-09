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
	"strings"
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

type ChunkedCancelReq struct {
	ChunkToken string `json:"chunk_token" validate:"required"`
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
		FileSize:    int64(data.FileSize),
		Name:        data.Name,
		IpAddr:      *ipAddr,
		ChunksLeft:  chunkCount,
		ChunksTotal: chunkCount,
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
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error marshaling response:", err)

		err := server.queries.ChunkedDelete(context.Background(), chunkedId)
		if err != nil {
			// scheduled cleanup will catch this even if it fails
			log.Println("Error querying db:", err)
		}

		return
	}

	rw.Header().Set("Content-Type", "application/json")
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
	_, err = jwt.ParseWithClaims(req.FormValue("chunk_token"), &claims, func(token *jwt.Token) (any, error) {
		return []byte(server.cfg.SecretKey), nil
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
	err = os.MkdirAll(chunkDir, 0755)
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
	chunk, header, err := req.FormFile("chunk")
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer chunk.Close()
	if int32(chunkIndex) == row.ChunksTotal {
		if header.Size != row.FileSize%int64(server.cfg.ChunkSize) {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else {
		if header.Size != int64(server.cfg.ChunkSize) {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	chunkPath := path.Join(chunkDir, fmt.Sprint(chunkIndex))
	_, err = os.Stat(chunkPath)
	if err == nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else if !errors.Is(err, os.ErrNotExist) {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error stating file:", err)
		return
	}

	dst, err := os.OpenFile(chunkPath, os.O_WRONLY|os.O_CREATE, 0644)
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
		err := ChunkedToFile(claims.ChunkedId, int(row.ChunksTotal), rw)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			os.Remove(chunkPath)
			return
		}
	}

	return
}

func ChunkedToFile(chunkedId int32, chunksTotal int, rw http.ResponseWriter) error {
	chunkDir := path.Join(server.chunkedPath, fmt.Sprint(chunkedId))
	chunkFullFilePath := path.Join(chunkDir, "fullfile")
	fullFile, err := os.OpenFile(chunkFullFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error creating file:", err)
		return err
	}

	for chunkIndex := 1; chunkIndex <= chunksTotal; chunkIndex++ {
		chunkPath := path.Join(chunkDir, fmt.Sprint(chunkIndex))
		chunk, err := os.Open(chunkPath)
		if err != nil {
			log.Println("Error opening file:", err)
			fileCloseAndRemove(fullFile, chunkFullFilePath)
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
	_, err = fullFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("Error seeking file:", err)
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
		err := server.queries.ChunkedDelete(context.Background(), chunkedId)
		if err != nil {
			fileCloseAndRemove(fullFile, chunkFullFilePath)
			log.Println("Error querying db:", err)
			return err
		}

		fileCloseAndRemove(fullFile, chunkDir)
		fileId56 := base56.Encode(uint64(row.ID))
		fileName := fmt.Sprintf("%v%v", fileId56, mimetype.Lookup(row.MimeType).Extension())

		fmt.Fprintf(rw, "%v/%v\n", server.cfg.PublicUrl, fileName)
		return nil
	}

	_, err = fullFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("Error seeking file:", err)
		os.Remove(chunkFullFilePath)
		return err
	}
	mtype, err := mimetype.DetectReader(fullFile)
	fullFile.Close()
	if err != nil {
		log.Println("Error detecting mimetype:", err)
		os.Remove(chunkFullFilePath)
		return err
	}
	fileId, err := server.queries.ChunkedToFile(context.Background(), db.ChunkedToFileParams{
		MimeType: strings.Split(mtype.String(), ";")[0],
		Hash:     hash,
		ID:       chunkedId,
	})
	if err != nil {
		os.Remove(chunkFullFilePath)
		log.Println("Error querying db:", err)
		return err
	}

	fileId56 := base56.Encode(uint64(fileId))
	fileName := fmt.Sprintf("%v%v", fileId56, mtype.Extension())
	filePath := path.Join(server.storagePath, fileName)
	err = os.Rename(chunkFullFilePath, filePath)
	if err != nil {
		os.Remove(chunkFullFilePath)
		log.Println("Error moving file:", err)

		err = server.queries.FileDelete(context.Background(), fileId)
		if err != nil {
			os.Remove(chunkDir)
			log.Println("Error querying db:", err)
			return err
		}

		return err
	}

	err = server.queries.ChunkedDelete(context.Background(), chunkedId)
	if err != nil {
		// scheduled cleanup will catch this even if it fails
		log.Println("Error querying db:", err)
	}
	err = os.RemoveAll(chunkDir)
	if err != nil {
		// scheduled cleanup will catch this even if it fails
		log.Println("Error removing chunkDir:", err)
	}

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

func uploadChunkedCancel(rw http.ResponseWriter, req *http.Request) {
	var data ChunkedCancelReq

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

	var claims ChunkedJwtClaims
	_, err = jwt.ParseWithClaims(data.ChunkToken, &claims, func(token *jwt.Token) (any, error) {
		return []byte(server.cfg.SecretKey), nil
	})
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err = server.queries.ChunkedDelete(context.Background(), claims.ChunkedId)
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

	chunkDir := path.Join(server.chunkedPath, fmt.Sprint(claims.ChunkedId))
	err = os.RemoveAll(chunkDir)
	if err != nil {
		// scheduled cleanup will catch this even if it fails
		log.Println("Error removing chunkDir:", err)
	}

	return
}
