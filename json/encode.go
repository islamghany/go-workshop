package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
 beacause of json is already a string we can write json response using io.WriteString(),
  w.Write() or fmt.Fprint()

in fact the only special thing we have do is to set the Content-Type header to application/json
so the client knows he is receiving json so he can interpret it.

at a high level Go’s encoding/json package us with two options for encoding things to json
 we can either choose json.Marshal or json.Encoder type.


 *********** How Go types are encoded **********
   bool  							---> 	json boolean
   string							---> 	json string
   int*, uint*,float*, rune 		---> 	json number
   array, slice						--->  	json array
   struct, map						---> 	json object
   nil pointer, interface values,
	   								---> json null
	   slices, maps, etc.
   chan, func, complex*	   			---> not supported
   time.Time 						---> RFC3339-format JSON string ex:  "2020-11-08T06:27:59+01:00"
   []byte							---> Base64-encoded JSON string

   hence :
	  - Encoding of nested objects is supported. So, for example,
	  if you have a slice of structs in Go that will encode to an array of objects in JSON.

	  - Channels, functions and complex number types cannot be encoded.
	  If you try to do so, you’ll get a json.UnsupportedTypeError error at runtime.



1- json.Marshal():
	the way the json.marshal works is strightforward, we pass a Go natvie object as parameter
	and it return json representation of that object in []byte slice, the function looks loke that
		func Marshal(v interface{}) ([]byte, error)

2- json.encoder type:
	it allow to encode an object to json and write that json to an output stream in a single step.
	for example
			func writeJSON(w http.ResponseWriter){
				data := map[string]string{
					"hello": "world",
				}

				w.Header().Set("Content-Type", "application/json")

				// Use the json.NewEncoder() function to initialize a json.Encoder instance that
				// writes to the http.ResponseWriter. Then we call its Encode() method, passing in
				// the data that we want to encode to JSON (which in this case is the map above). If
				// the data can be successfully encoded to JSON, it will then be written to our
				// http.ResponseWriter.
				err := json.NewEncoder(w).Encode(data)
				if err != nil {
					fmt.Println("an error has occur ", err)
					return
				}
			}
	When we call json.NewEncoder(w).Encode(data) the JSON is created and written to the http.ResponseWriter
	 in a single step, which means there’s no opportunity to set HTTP response headers conditionally based
	  on whether the Encode() method returns an error or not.

	Imagine, for example, that you want to set a Cache-Control header on a successful response,
	 but not set a Cache-Control header if the JSON encoding fails
	  and you have to return an error response.

	  You could set the Cache-Control header and then delete it from the header map again in
	  the event of an error — but that’s pretty hacky.

	  Another option is to write the JSON to an interim bytes.Buffer instead of directly to
	  the http.ResponseWriter. You can then check for any errors, before setting the Cache-Control
	  header and copying the JSON from the bytes.Buffer to http.ResponseWriter.
	  But once you start doing that, it’s simpler and cleaner (as well as slightly faster)
	  o use the alternative json.Marshal() approach instead

	  -regarding the performance:
		json.Marshal() requires ever so slightly more memory (B/op) than json.Encoder,
		 and also makes one extra heap memory allocation (allocs/op).

		There’s no obvious observable difference in the average runtime (ns/op) between
		the two approaches. Perhaps with a larger benchmark sample or a larger data set a difference
		might become clear, but it’s likely to be in the order of
		microseconds, rather than anything larger.


		-- encoding a struct
			type Movie struct {
					ID        int64     `json:"id"`
					CreatedAt time.Time `json:"-"` // Use the - directive = delete it from parsing
					Title     string    `json:"title"`
					Year      int32     `json:"year,omitempty"`    // Add the omitempty directive = delete it if empty
					Runtime   int32     `json:"runtime,omitempty,string"` / Add the string directive
					Genres    []string  `json:"genres,omitempty"`  // Add the omitempty directive
					Version   int32     `json:"version"`
			}


		// Use the json.MarshalIndent() function so that whitespace is added to the encoded
		// JSON. Here we use no line prefix ("") and tab indents ("\t") for each element.
		js, err := json.MarshalIndent(data, "", "\t")
		 ===
				{
						"environment": "development",
						"status": "available",
						"version": "1.0.0"
				}
*/

func writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	jsonData, err := json.Marshal(data)

	if err != nil {
		fmt.Println("an error has occur ", err)
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	jsonData = append(jsonData, '\n')

	// Go does't throw an error when iterating on a nil object so its fine.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application")
	w.WriteHeader(status)
	w.Write(jsonData)
	return nil
}
func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"environment": "development",
			"status":      "available",
			"version":     "1.0.0",
			"day":         25,
		}
		err := writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			log.Print(err)
			return
		}
	})
	http.ListenAndServe(":8000", nil)
}
