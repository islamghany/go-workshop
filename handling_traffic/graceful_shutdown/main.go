package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
 when we stop our API application (usually by pressing Ctrl+C) it is terminated immediately with no opportunity
 for in-flight HTTP requests to complete. This isn’t ideal for two reasons:

 it means that clients won’t receive responses to their in-flight requests —
 all they will experience is a hard closure of the HTTP connection.

Any work being carried out by our handlers may be left in an incomplete state.


----------------------

When our application is running, we can terminate it at any time by sending it a specific signal.
 A common way to do this, which you’ve probably been using,
 is by pressing Ctrl+C on your keyboard to send an interrupt signal — also known as a SIGINT.

Signal		Description	Keyboard 					shortcut	Catchable
SIGINT		Interrupt from keyboard					Ctrl+C		Yes
SIGQUIT		Quit from keyboard						Ctrl+\		Yes
SIGKILL		Kill process (terminate immediately)	-			No
SIGTERM		Terminate process in orderly manner		-			Yes
*/
func main() {
	// Create a shutdownError channel. We will use this to receive any errors returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)

	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.Write([]byte("hello there"))
	})

	mux.HandleFunc("/hangout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Minute)
		w.Write([]byte("i am hanging out"))
	})
	srv := http.Server{
		Addr:         ":8000",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start a background goroutine.
	go func() {
		// Create a quit channel which carries os.Signal values.
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel. Any other signals will not be caught by
		// signal.Notify() and will retain their default behavior.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from the quit channel. This code will block until a signal is
		// received.
		s := <-quit

		// Log a message to say that the signal has been caught. Notice that we also
		// call the String() method on the signal to get the signal name and include it
		// in the log entry properties.
		log.Println("caught signal\n", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		defer cancel()

		// Call Shutdown() on our server, passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because of a problem closing the listeners, or
		// because the shutdown didn't complete before the 5-second context deadline is
		// hit). We relay this return value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)
	}()

	log.Println("The Server is running on addr: ", srv.Addr)
	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Println("gracfuly shutteddown")
		log.Fatal(err)
	}

	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		log.Println("forced interupt")
		log.Fatal(err)
	}

}
