package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// generate a `Certificate` struct
	cert, _ := tls.LoadX509KeyPair( "server.crt", "server.key" )

	// create a custom server with `TLSConfig`
	s := &http.Server{
		Addr: ":5050",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{ cert },
		},
	}

	http.HandleFunc( "/", func( res http.ResponseWriter, req *http.Request ) {
		fmt.Fprint( res, "Hello World!" )
	} )

	log.Fatal( s.ListenAndServeTLS("server.crt", "server.key") )
}
