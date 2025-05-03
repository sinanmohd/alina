package db

import (
	"context"
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"sinanmohd.com/alina/internal/config"
)

//go:embed migrations/*.sql
var migrations embed.FS

func NewWithSetup(cfg config.DatabaseConfig) (*Queries, *pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.Url)
	if err != nil {
		log.Println("Error creating pool:", err)
		return nil, nil, err
	}

	driver, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Println("Error creating iofs:", err)
		return nil, nil, err
	}
	m, err := migrate.NewWithSourceInstance("iofs", driver, cfg.Url)
	if err != nil {
		log.Println("Error creating migrate:", err)
		return nil, nil, err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println("Error running migrations:", err)
		return nil, nil, err
	}

	return New(pool), pool, nil
}
