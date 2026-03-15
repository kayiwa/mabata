package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/marcboeker/go-duckdb"

	"github.com/kayiwa/mabata/internal/auth"
	"github.com/kayiwa/mabata/internal/config"
	"github.com/kayiwa/mabata/internal/duck"
	"github.com/kayiwa/mabata/internal/httpui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := sql.Open("duckdb", cfg.DuckDBPath)
	if err != nil {
		log.Fatalf("open duckdb: %v", err)
	}
	defer db.Close()

	if err := duck.Init(db); err != nil {
		log.Fatalf("init duckdb: %v", err)
	}

	authn, err := auth.New(cfg)
	if err != nil {
		log.Fatalf("init oidc: %v", err)
	}

	mux := http.NewServeMux()
	h := httpui.New(cfg, db, authn)
	h.Register(mux)

	log.Printf("listening on %s", cfg.AppAddr)
	if err := http.ListenAndServe(cfg.AppAddr, mux); err != nil {
		log.Fatalf("http server: %v", err)
	}
}
