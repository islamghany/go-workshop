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
	mux := http.NewServeMux()

	mux.HandleFunc("/home", Home)
	mux.HandleFunc("/", app.notFoundResponse)
	app.logger.PrintInfo("the server is running", nil)

	app.logger.PrintFatal(http.ListenAndServe(":8000", app.recoverPanic(mux)), nil)

}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome to the site")
}

/*

At the moment any panics in our API handlers will be recovered automatically by Go’s
http.Server. This will unwind the stack for the affected goroutine (calling any
deferred functions along the way), close the underlying HTTP connection, and
log an error message and stack trace.

This behavior is OK, but it would be better for the client if we
could also send a 500 Internal Server Error response to explain that
something has gone wrong — rather than just closing the HTTP connection
with no context.


*/
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type interface{}, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
