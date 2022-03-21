package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

/*
	Just like JSON encoding, there are two approaches that you can take to decode JSON
	into a native Go object: using a json.Decoder type or using the json.Unmarshal() function.

	Both approaches have their pros and cons, but for the purpose of decoding JSON from
	a HTTP request body, using json.Decoder is generally the best choice.
	It’s more efficient than json.Unmarshal(), requires less code, and offers
	some helpful settings that you can use to tweak its behavior.

	***** managing bad reqeust *****
	What if the client sends something that isn’t JSON, like XML or some random bytes?
	What happens if the JSON is malformed or contains an error?
	What if the JSON types don’t match the types we are trying to decode into?
	What if the request doesn’t even contain a body?


	***** Triaging the Decode error *****
	json.SyntaxError and io.ErrUnexpectedEOF    ---> here is a syntax problem with the JSON being decoded.

	json.UnmarshalTypeError  					---> A JSON value is not appropriate for the destination Go type.

	json.InvalidUnmarshalError					---> The decode destination is not valid (usually because it is not a pointer). This is actually a problem with our application code, not the JSON itself.

	io.EOF 										---> The JSON being decoded is empty.




	********* Restricting Inputs **********
	One such thing is dealing with unknown fields. For example, you can try sending a request containing the unknown field

	Fortunately, Go’s json.Decoder provides a DisallowUnknownFields() setting that we can
	use to generate an error when this happens.
*/

func readJSON(w http.ResponseWriter, r *http.Request, input interface{}) error {

	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	// Decode the request body to the destination.
	err := dec.Decode(input)
	if err != nil {
		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

			// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
			// for syntax errors in the JSON. So we check for this using errors.Is() and
			// return a generic error message. There is an open issue regarding this at
			// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

			// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
			// JSON value is the wrong type for the target destination. If the error relates
			// to a specific field, then we include that in our error message to make it
			// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

			// An io.EOF error will be returned by Decode() if the request body is empty. We
			// check for this with errors.Is() and return a plain-english error message
			// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

			// If the JSON contains a field which cannot be mapped to the target destination
			// then Decode() will now return an error message in the format "json: unknown
			// field "<name>"". We check for this, extract the field name from the error,
			// and interpolate it into our custom error message. Note that there's an open
			// issue at https://github.com/golang/go/issues/29035 regarding turning this
			// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request body exceeds 1MB in size the decode will now fail with the
		// error "http: request body too large". There is an open issue about turning
		// this into a distinct error type at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
			// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
			// pointer to Decode(). We catch this and panic, rather than returning an error
			// to our handler. At the end of this chapter we'll talk about panicking
			// versus returning errors, and discuss why it's an appropriate thing to do in
			// this specific situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

			// For anything else, return the error message as-is.
		default:
			return err
		}
	}
	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Title   string   `json:"title"`
			Year    int32    `json:"year"`
			Runtime int32    `json:"runtime"`
			Genres  []string `json:"genres"`
		}
		err := readJSON(w, r, &input)
		if err != nil {
			log.Print(err)
			return
		}
		w.Write([]byte(fmt.Sprintf("\v", input)))
	})
	http.ListenAndServe(":8000", nil)
}
