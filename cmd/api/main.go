package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"greenlight.vishaaxl.net/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   db
}

type db struct {
	dsn string
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cfg := config{
		port: 8080,
		env:  "development",
		db: db{
			dsn: "postgres://greenlight:mysecretpassword@localhost/greenlight?sslmode=disable",
		},
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Println("database connection established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on port %d", cfg.env, cfg.port)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(time.Minute * 15)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
