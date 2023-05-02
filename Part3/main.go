package main

import (
	"log"
	"net/http"

	"github.com/goji/httpauth"
)

//Using Third-Party Middleware
//goji/httpauth, which provides HTTP Basic Authentication functionality.

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {

	//When using this package you call a helper function in order to setup the chainable middleware. Specifically, you call the httpauth.SimpleBasicAuth()
	//function, and this returns a middleware function with the signature func(http.Handler) http.Handler — which you can then use in exactly
	//the same way as any custom-built middleware.
	authHandler := httpauth.SimpleBasicAuth("Moises", "password2")

	/*
		Using the proper username (Moises) and password (password2) should yield a "OK" answer, while entering the incorrect username and password should
		cause the prompt to be presented again. Clicking "Cancel" should yield a plain-text "Unauthorized" response.
	*/

	// Use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)    //HandlerFunc serves as an adapter to allow the use of ordinary http.Handlers, in this case used with Final\
	mux.Handle("/", authHandler(finalHandler)) //function to register this with our new servemux, so it acts as the handler for all incoming requests with the URL path stated

	log.Print("Listening on :3000...") //Printing out the port number

	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second parameter.
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
