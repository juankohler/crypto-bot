package common

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Dependencies struct {
	Mux *http.ServeMux
	DB  *sqlx.DB
}

func BuildDependencies(cfg *Config) (*Dependencies, error) {
	db, err := sqlx.Connect("sqlite3", cfg.Database)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	return &Dependencies{
		Mux: mux,
		DB:  db,
	}, nil
}
