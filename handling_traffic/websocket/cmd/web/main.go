package main

import (
	"log"
	"net/http"

	"github.com/islamghany/go-workshop/handling_traffic/websocket/internals/handlers"
)

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Running the server.")
	err := srv.ListenAndServe()
	log.Fatal(err)
}
