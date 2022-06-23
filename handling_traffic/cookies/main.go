package main

import (
	"log"
	"net/http"

	"github.com/islamghany/go-workshop/helpers"
)

/*
Cookies are a way to store information at the client end. The client can be a browser,
a mobile application, or anything which makes an HTTP request. Cookies are basically some
files that are stored in the cache memory of your browser. When you are browsing any website
which supports cookies will drop some kind of information related to your activities in the cookie.
This information could be anything. Cookies in short store historical information about t
he user activities. This information is stored on the client’s computer. Since a cookie
is stored in a file,  hence this information is not lost even when the user closes a browser
window or restarts the computer. A cookie can also store the login information. In fact, login
information such as tokens is generally stored in cookies only. Cookies are stored per domain.
Cookies stored locally belonging to a particular domain are sent in each request to that domain.
They are sent in each request as part of headers.
So essentially cookie is nothing but a header.


type Cookie struct {
	Name  string
	Value string

	// cookie scope
	Path       string    // optional
	Domain     string    // optional the domain that this cookie can set in it

	// if you does't provide a expire the cookie will destroy after
	// the browser close, it called the cookie session
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool // if set to true, the cookie only will sent to the https origins only
	HttpOnly bool // the cookir can not be accessed from the browser
	SameSite SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
*/
func main() {
	router := helpers.Router()

	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:   "name",
			Value:  "mostafa",
			Domain: ".localhost",
			MaxAge: 10,
		}

		http.SetCookie(w, &cookie)
		w.Write([]byte("hello in the main page"))
	})
	router.HandlerFunc(http.MethodGet, "/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello in the about page"))
	})
	log.Fatal(helpers.Serve(":8000", router))
}
