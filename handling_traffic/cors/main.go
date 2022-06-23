package main

import (
	"log"
	"net/http"

	"github.com/islamghany/go-workshop/helpers"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methos", "GET, POST, OPTIONS, DELETE, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		next.ServeHTTP(w, r)
	})
}
func main() {
	router := helpers.Router()

	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("hello in the main page"))
	})
	router.HandlerFunc(http.MethodGet, "/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello in the about page"))
	})

	router.HandlerFunc(http.MethodOptions, "/", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(204)
	})
	log.Fatal(helpers.Serve(":8000", enableCORS(router)))
}
