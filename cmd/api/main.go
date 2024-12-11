package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
	"greenlight.vishaaxl.net/internal/data"
	"greenlight.vishaaxl.net/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port    int
	env     string
	db      db
	limiter limiter
}

type limiter struct {
	rps     float64
	burst   int
	enabled bool
}

type db struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	cfg := config{
		port: 8080,
		env:  "development",
		db: db{
			dsn:          "postgres://greenlight:mysecretpassword@localhost/greenlight?sslmode=disable",
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime:  time.Minute * 15,
		},
		limiter: limiter{
			rps:     2,
			burst:   4,
			enabled: true,
		},
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("database connection established", nil)

	app := &application{
		logger: logger,
		config: cfg,
		models: data.NewModels(db),
	}

	// Call app.serve() to start the server.
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
