// internal/server/server.go
package server

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql

	"github.com/molu/stock-management-system/internal/config"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
	"github.com/molu/stock-management-system/internal/router"
)

type Server struct {
	config  *Config
	router  http.Handler
	httpSrv *http.Server
	db      *sql.DB
	queries *db.Queries
}

type Config struct {
	Address   string
	DB        *sql.DB
	Queries   *db.Queries
	Env       string
	JWTSecret string
}

func New(cfg *config.Config) (*Server, error) {
	ctx := context.Background()

	// Open database connection using database/sql + pgx driver
	dbConn, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Configure connection pool (optional but recommended)
	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(25)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := dbConn.PingContext(ctx); err != nil {
		return nil, err
	}

	// Create sqlc queries instance
	queries := db.New(dbConn)

	// Create server config
	srvCfg := &Config{
		Address:   cfg.ServerAddress,
		DB:        dbConn,
		Queries:   queries,
		Env:       cfg.Environment,
		JWTSecret: cfg.JWTSecret,
	}

	// Create router
	r := router.New(srvCfg.Queries, srvCfg.JWTSecret)

	// Create HTTP server
	srv := &Server{
		config:  srvCfg,
		router:  r,
		db:      dbConn,
		queries: queries,
		httpSrv: &http.Server{
			Addr:         cfg.ServerAddress,
			Handler:      r,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	return srv, nil
}

func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.db != nil {
		_ = s.db.Close()
	}
	return s.httpSrv.Shutdown(ctx)
}
