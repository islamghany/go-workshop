package helpers

import (
	"log"
	"net/http"
	"time"
)

func Serve(addr string, handler http.Handler) error {
	srv := http.Server{
		Addr:         addr,
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("Starting a server on port: ", addr)
	err := srv.ListenAndServe()

	return err
}
