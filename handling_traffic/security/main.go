package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

/*
HTTPS is essentially HTTP sent across a TLS (Transport Layer Security) connection.
Because it’s sent over a TLS connection the data is encrypted and signed,
which helps ensure its privacy and integrity during transit.

If you’re not familiar with the term, TLS is essentially the modern version
of SSL (Secure Sockets Layer). SSL now has been officially deprecated due to
security concerns, but the name still lives on in the public consciousness and is
often used interoperably with TLS. For clarity and accuracy, we’ll stick with the
term TLS throughout this book.

Before our server can start using HTTPS, we need to generate a TLS certificate.

For production servers I recommend using Let’s Encrypt to create your TLS certificates,
but for development purposes the simplest thing to do is to generate your own self-signed certificate.

A self-signed certificate is the same as a normal TLS certificate,
except that it isn’t cryptographically signed by a trusted certificate authority.
This means that your web browser will raise a warning the first time it’s used,
but it will nonetheless encrypt HTTPS traffic correctly and is fine for development
and testing purposes



If you open up your web browser and visit https://localhost:8000/
you will probably get a browser warning similar to the screenshot below.
*/

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", sayHello)

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MaxVersion: tls.VersionTLS12,
		MinVersion: tls.VersionTLS12,
	}

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      mux,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute, // keep a live header.
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello There!")
}
