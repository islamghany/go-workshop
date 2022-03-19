package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// (DefaultServerMux) is a struct that implement the ServeHTTP
// by dafault when we use DefaultServerMux by itself, it has a method called ServeHTTP
// when we pass the pattern and the Hanlder func to HandlerFun,
// it save this pattern with this func in map so when we go for a specific pattern
// first and call ServerHTTP method then this method mapped us yo the specific Handler Func

// so for any struct that has the method ServeHTTP it can be use as a multiplexer.

// // 1 - start a server with the default multiplexer (DefaultServerMux)

// func main() {
// 	http.HandleFunc("/", SayHello)

// 	// by default when we pass nil it will use the DefaultServerMux struct or we can
// 	// use it explicitly
// 	http.ListenAndServe(":8000", nil)
// }

// // 2 - start a server with a custom multiplexer

// type CustomMuliplexer struct {
// }

// func (p *CustomMuliplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":

// 		SayHello(w, r)
// 		return

// 	case "/goodbye":

// 		SayGoodbye(w, r)
// 		return

// 	default:

// 		SayHello(w, r)

// 	}
// }

// func main() {

// 	mux := &CustomMuliplexer{}
// 	http.ListenAndServe(":8000", mux)
// }

// // 3 - start a server with the use of  ServeMux
// // unlike the last approach, its not simple to create a multiple endpoints
// // we can use a make a struct that has the same attributes of the DefaultServeMux
// // and use its HandleFunc to map the pattern with the handlers.

// func main() {
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/goodbye", SayGoodbye)
// 	mux.HandleFunc("/", SayHello)

// 	http.ListenAndServe(":8000", mux)
// }

// // 4 - start a server with controlling the config
// // when we start a server there are configurations that we want to control of
// // such the header size, handler , errorlog, addr, readtimeout etc..
// // this approuch help us to build this config from scratch.

// func main() {

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/goodbye", SayGoodbye)
// 	mux.HandleFunc("/", SayHello)

// 	s := &http.Server{
// 		Addr:           ":8000",
// 		Handler:        mux, // the default value is the (DefaultServeMux)
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}

// 	s.ListenAndServe()

// }

// this 4 ways is the basics approaches that we can start a server in go
// but go community had a lot of libraries to help us manage this routes traffic effeciently

// 1- using httprouter : a lightweight HTTP router
/*
	httprouter plays well with the inbuilt http.Handler
	httprouter explicitly says that a request can only match to one route or none
	The router's design encourages building sensible, hierarchical RESTful APIs
	You can build efficient static file servers
*/

func main() {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", SayHello)
	router.HandlerFunc(http.MethodGet, "/goodbye", SayGoodbye)
	log.Fatal(http.ListenAndServe(":8000", router))
}

// there is also another library goriall mux : a powerful HTTP router
