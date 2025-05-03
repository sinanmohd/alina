package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"sinanmohd.com/alina/db"
)

type ChunkedStartReq struct {
	FileSize int    `json:"file_size" validate:"required,numeric,gt=1"`
	Name     string `json:"name" validate:"required"`
}

type ChunkedStartResp struct {
	ChunkToken string `json:"chunk_token"`
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

	chunkedId, err := server.queries.ChunkedCreate(context.Background(), db.ChunkedCreateParams{
		FileSize:   int64(data.FileSize),
		Name:       data.Name,
		IpAddr:     *ipAddr,
		ChunksLeft: int32(math.Ceil(float64(data.FileSize) / float64(server.cfg.ChunkSize))),
	})

	chunkToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": chunkedId,
	}).SignedString([]byte(server.cfg.SecretKey))
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

func uploadChunkedProgress(rw http.ResponseWriter, r *http.Request) {
}

func uploadChunkedCancel(rw http.ResponseWriter, r *http.Request) {
}
