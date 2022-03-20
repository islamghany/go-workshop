package main

import (
	"net/http"
)

/*

***************** Intro *********************

we can thing of a Go web application as a chain of ServeHTTP() method being called one after another

-- any object that has the method ServeHTTP(w,r*) is a handler

so we can do that

	type home struct {}

	func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is my home page"))
	}

	mux := http.NewServeMux()
	mux.Handle("/", &home{})

because of making a object that has the method ServeHTTP() with the same signture is long-winded
and a bit confusing

go proveded us with cool feature to overcome this redunant work
 1 - it has the http.HandlerFunc adaptor to take a function with this signture
	func home(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is my home page"))
	}

 and add to it the 	ServeHTTP() method.
 and then we can do this

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(home))

2- the second cool way is a syntactic suger from the first solution that we pass the pattern and the home func as handler
 to a function called

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

***************** MIDDLEWARE *********************

when our server receives an incoming request it calls the servemux's ServeHTTP() method, this looks up
the relevant handler based on the URL path and in turn calls that handler's ServeHTTP() method.

so the idea of middleware that we insert another handler into the chain, the middleware handler executes
some logic and then call the ServeHttp() method  of the next handler chain.

the pattern

func myMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Reqesut){
		// do some login here

		next.ServeHTTP(w,r)
	}

	return http.HandlerFunc(fn)
}
*/

// // middleware on a normal function
// func enableCORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
// 		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
// 		next.ServeHTTP(w, r)
// 	})
// }

// middleware on a handler function
func enableCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		next.ServeHTTP(w, r)
	})
}
func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
func main() {
	mux := http.NewServeMux()

	//mux.HandleFunc("/", enableCORSMiddleware(sayHello))
	//http.ListenAndServe(":8000", mux)

	mux.HandleFunc("/", sayHello)
	http.ListenAndServe(":8000", enableCORSMiddleware(mux))
}
