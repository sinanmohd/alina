package main

import (
	"os"

	"sinanmohd.com/alina/db"
	"sinanmohd.com/alina/internal/config"
	"sinanmohd.com/alina/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		os.Exit(1)
	}

	queries, pool, err := db.NewWithSetup(cfg.Db)
	if err != nil {
		os.Exit(1)
	}
	defer pool.Close()

	err = server.Run(cfg.Server, queries)
	if err != nil {
		os.Exit(1)
	}
}
