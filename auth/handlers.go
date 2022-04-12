package main


import (
	"net/http"
)

func (app *application) hello(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hello world"))
}