package main

import (
	"log"
	"net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello there!"))
}

func main() {

	db, err := OpenDB()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("database connection pool established.")

	http.HandleFunc("/", sayHello)

	err = http.ListenAndServe(":8000", nil)
	log.Fatal(err)
}
