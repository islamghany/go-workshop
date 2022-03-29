package main

import (
	"net/http"

	"github.com/islamghany/go-workshop/handling_traffic/websocket/internals/handlers"

	"github.com/julienschmidt/httprouter"
)

func routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", handlers.Home)
	router.HandlerFunc(http.MethodGet, "/ws", handlers.WsEndpoint)
	fileServer := http.FileServer(http.Dir("./handling_traffic/websocket/static/"))
	router.Handler(http.MethodGet, "/static/", http.StripPrefix("./handling_traffic/websocket/static", fileServer))
	return router
}
