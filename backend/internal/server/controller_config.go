package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PublicConfigResp struct {
	FileSizeLimit int `json:"file_size_limit"`
	ChunkSize     int `json:"chunk_size"`
}

func publicConfig(rw http.ResponseWriter, req *http.Request) {
	sCfg := PublicConfigResp{
		FileSizeLimit: server.cfg.FileSizeLimit,
		ChunkSize:     server.cfg.ChunkSize,
	}

	resp, err := json.Marshal(sCfg)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error marshaling response:", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
	return
}
