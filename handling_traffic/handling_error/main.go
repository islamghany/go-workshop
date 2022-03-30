package main

import (
	"fmt"
	"net/http"
	"os"
)

type application struct {
	logger *Logger
}

func main() {
	app := &application{
		logger: New(os.Stdout, LevelInfo),
	}

	http.HandleFunc("/home", Home)
	http.HandleFunc("*", app.notFoundResponse)
	app.logger.PrintInfo("the server is running", nil)
	app.logger.PrintFatal(http.ListenAndServe(":8000", nil), nil)

}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome to the site")
}
