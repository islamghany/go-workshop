package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodGet, "/", app.hello)

	return router
}
