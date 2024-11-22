package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cfg := config{
		port: 8080,
		env:  "development",
	}

	app := &application{
		config: cfg,
		logger: logger,
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
