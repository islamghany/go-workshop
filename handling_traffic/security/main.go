package main

import (
	"fmt"
	"net/http"
)

func main() {
	srv := http.NewServeMux()

	srv.HandleFunc("/", sayHello)

	http.ListenAndServe(":8000", srv)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello There!")
}
