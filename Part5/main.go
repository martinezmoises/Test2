package main

import (
	"log"
	"net/http"
)

//Uses the same pattern for constructing a handler

func middlewareA(next http.Handler) http.Handler {
	// It accepts a handler as a parameter and returns a handler
	//This is usefull because
	//1. it returns a handler we can register the middleware function directly with the standard http.ServeMux router in Go's net/http package.
	//2. We can create an arbitrarily long handler chain by nesting middleware functions inside each other.

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//closes-over the message variable to form a closure.

		//this is executed on the way down to the handler
		log.Println("Executing middleware A")

		next.ServeHTTP(w, r)

		log.Println("Executing middleware A again") //this is executed on the way up to the client
	})
}

func middlewareB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//this is executed on the way down to the handler
		log.Println("Executing middleware B")

		if r.URL.Path == "/cherry" { //Conditional return statement
			return
		}
		next.ServeHTTP(w, r)                        // pass the next handler in the chain as a variable. The ServeHTTP method writes out the HTTP response
		log.Println("Executing middleware B again") //this is executed on the way up to the client
	})
}

func ourHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing the handler...")
	w.Write([]byte("Carrots\n"))

}

func main() {
	//multiplexer act lik a router
	mux := http.NewServeMux() // Use the http.NewServeMux() function to create an empty servemux.

	mux.Handle("/", middlewareA(middlewareB(http.HandlerFunc(ourHandler)))) //function to register this with our new servemux, so it acts as the
	//handler for all incoming requests with the URL path stated

	log.Print("starting server on :4000")    //printing port number
	err := http.ListenAndServe(":4000", mux) //we create a server
	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second parameter.
	log.Fatal(err)
}
