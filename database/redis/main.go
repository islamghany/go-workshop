package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

/*
	The first is the Dial() function, which returns a new connection to our Redis server.

	The second is the Do() method, which sends a command to our Redis server across
	the connection. This returns the reply from Redis as an interface{} type,
	along with any error if applicable.
*/

// Define a custom struct to hold Album data. Notice the struct tags?
// These indicate to redigo how to assign the data from the reply into
// the struct.
// type Album struct {
// 	Title  string  `redis:"title"`
// 	Artist string  `redis:"artist"`
// 	Price  float64 `redis:"price"`
// 	Likes  int     `redis:"likes"`
// }

// func main() {

// 	// Establish a connection to the Redis server listening on port
// 	// 6379 of the local machine. 6379 is the default port, so unless
// 	// you've already changed the Redis configuration file this should
// 	// work.
// 	conn, err := redis.Dial("tcp", "localhost:6379")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Importantly, use defer to ensure the connection is always
// 	// properly closed before exiting the main() function.
// 	defer conn.Close()

// 	// res, err := conn.Do("HMSET", "album:1", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)

// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// log.Println(res, err)

// 	res, err := redis.Values(conn.Do("HGETALL", "album:1"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create an instance of an Album struct and use redis.ScanStruct()
// 	// to automatically unpack the data to the struct fields. This uses
// 	// the struct tags to determine which data is mapped to which
// 	// struct fields.
// 	var album Album

// 	err = redis.ScanStruct(res, &album)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("%+v\n", album)
// 	// mux := http.NewServeMux()

// 	// mux.HandleFunc("/")

// 	// http.ListenAndServe(":8000",mux)
// }

func main() {
	// Initialize a connection pool and assign it to the pool global
	// variable.
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/album", showAlbum)
	mux.HandleFunc("/like", addLike)
	log.Println("Listening on :8000...")
	http.ListenAndServe(":8000", mux)
}

func showAlbum(w http.ResponseWriter, r *http.Request) {
	// Unless the request is using the GET method, return a 405 'Method
	// Not Allowed' response.
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Retrieve the id from the request URL query string. If there is
	// no id key in the query string then Get() will return an empty
	// string. We check for this, returning a 400 Bad Request response
	// if it's missing.
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// Validate that the id is a valid integer by trying to convert it,
	// returning a 400 Bad Request response if the conversion fails.
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Call the FindAlbum() function passing in the user-provided id.
	// If there's no matching album found, return a 404 Not Found
	// response. In the event of any other errors, return a 500
	// Internal Server Error response.
	bk, err := FindAlbum(id)
	if err == ErrNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Write the album details as plain text to the client.
	fmt.Fprintf(w, "%s by %s: £%.2f [%d likes] \n", bk.Title, bk.Artist, bk.Price, bk.Likes)
}

func addLike(w http.ResponseWriter, r *http.Request) {
	// Unless the request is using the POST method, return a 405
	// Method Not Allowed response.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Retrieve the id from the POST request body. If there is no
	// parameter named "id" in the request body then PostFormValue()
	// will return an empty string. We check for this, returning a 400
	// Bad Request response if it's missing.
	id := r.PostFormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// Validate that the id is a valid integer by trying to convert it,
	// returning a 400 Bad Request response if the conversion fails.
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Call the IncrementLikes() function passing in the user-provided
	// id. If there's no album found with that id, return a 404 Not
	// Found response. In the event of any other errors, return a 500
	// Internal Server Error response.
	err := IncrementLikes(id)
	if err == ErrNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Redirect the client to the GET /album route, so they can see the
	// impact their like has had.
	http.Redirect(w, r, "/album?id="+id, 303)
}
