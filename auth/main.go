package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/islamghany/go-workshop/auth/internals/data"
	_ "github.com/lib/pq"
)

type envelope map[string]interface{}
type config struct {
	port string
	db   struct {
		dsn string
	}
}

type application struct {
	models data.Models
	config config
}

func main() {

	conf := &config{
		port: ":8000",
	}
	conf.db.dsn = "postgres://test:islamghany@localhost/test"
	db, err := openDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	log.Println("database connection pool established")

	app := &application{
		config: *conf,
		models: data.NewModels(db),
	}

	srv := http.Server{
		Addr:         conf.port,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func openDB(conf *config) (*sql.DB, error) {

	db, err := sql.Open("postgres", conf.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 'postgres://test:test@localhost/test'
