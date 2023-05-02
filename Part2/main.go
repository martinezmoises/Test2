package main

import (
	"log"
	"mime"
	"net/http"
)

/*
We want to create some middleware which a) checks for the existence of a Content-Type header and
b) if the header exists, check that it has the mime type application/json. If either of those checks fail, we want our middleware
to write an error message and to stop the request from reaching our application handlers.

*/

//Uses the same pattern for constructing a handler
func enforceJSONHandler(next http.Handler) http.Handler {
	// It accepts a handler as a parameter and returns a handler

	//an anonymous function which closes-over the message variable to form a closure.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get the content type of the request
		contentType := r.Header.Get("Content-Type")

		//Condition that checks for the existence of a Content-Type header & if the header exists, check that it has the mime type application/json.
		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r) // pass the next handler in the chain as a variable. The ServeHTTP method writes out the HTTP response
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {

	// Use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final) //HandlerFunc serves as an adapter to allow the use of ordinary http.Handlers, in this case used with Final

	//function to register this with our new servemux, so it acts as the handler for all incoming requests with the URL path stated
	mux.Handle("/", enforceJSONHandler(finalHandler))
	//Printing out port number of server
	log.Print("Listening on :3000...")

	err := http.ListenAndServe(":3000", mux) // Then we create a new server and start listening for incoming requests
	log.Fatal(err)
}
