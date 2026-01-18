package server

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/molu/stock-management-system/internal/config"
	db "github.com/molu/stock-management-system/internal/db/sqlc"
	"github.com/molu/stock-management-system/internal/router"
)

type Server struct {
	config  *Config
	router  http.Handler
	httpSrv *http.Server
	db      *pgxpool.Pool
	queries *db.Queries
}

type Config struct {
	Address   string
	DB        *pgxpool.Pool
	Queries   *db.Queries
	Env       string
	JWTSecret string
}

func New(cfg *config.Config) (*Server, error) {
	ctx := context.Background()

	// Create database connection pool
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	// Create queries instance
	queries := db.New(pool)

	// Create server config
	srvCfg := &Config{
		Address:   cfg.ServerAddress,
		DB:        pool,
		Queries:   queries,
		Env:       cfg.Environment,
		JWTSecret: cfg.JWTSecret,
	}

	// Create router
	r := router.New(srvCfg.Queries, srvCfg.JWTSecret)

	srv := &Server{
		config:  srvCfg,
		router:  r,
		db:      pool,
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
	s.db.Close()
	return s.httpSrv.Shutdown(ctx)
}
