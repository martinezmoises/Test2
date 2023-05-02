package main

import (
	"log"
	"net/http"
)

//writes some log messages in your terminal window

// Illustrating the Flow of Control

//In web application development there is a lot shared functionality that you want to run
//Setting up this shared functionality as middleware—self-contained code that operates on a request independently before or after your standard
//application handlers—is one method to organize it.

//Uses the same pattern for constructing a handler
func middlewareOne(next http.Handler) http.Handler {

	// It accepts a handler as a parameter and returns a handler
	//This is usefull because
	//1. it returns a handler we can register the middleware function directly with the standard http.ServeMux router in Go's net/http package.
	//2. We can create an arbitrarily long handler chain by nesting middleware functions inside each other.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//an anonymous function which closes-over the message variable to form a closure.
		//w.Write([]byte(message)
		//We will use the same patter for the one above, Instead of passing a string into the closure, we would pass the next handler
		// in the chain as a variable, and then transfer control to this next handler by calling it's ServeHTTP() method.
		log.Print("Executing middlewareOne")
		next.ServeHTTP(w, r) // pass the next handler in the chain as a variable. The ServeHTTP method writes out the HTTP response
		log.Print("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {

	// We can stop control propagating through the chain at any point by issuing a return from a middleware handler
	// Uses the same idea of the first middleware except that this middleware conditional return
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareTwo")

		if r.URL.Path == "/foo" { //Conditional return statement

			return
		}

		/*
			 checks if the requested URL path is equal to "/foo". If the condition is true, the function immediately returns and exits
			the middlewareTwo function, without calling the next handler in the chain of middleware. This means that the subsequent middleware
			and the final handler will not be executed.
			However, if the URL path is not equal to "/foo", the middlewareTwo function calls the next handler in the chain of middleware,
			which in turn executes the final handler. The final handler simply writes a "OK" response to the client.

		*/

		next.ServeHTTP(w, r) // pass the next handler in the chain as a variable. The ServeHTTP method writes out the HTTP response
		log.Print("Executing middlewareTwo again")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {

	// Use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final) //HandlerFunc serves as an adapter to allow the use of ordinary http.Handlers, in this case used with Final\

	//function to register this with our new servemux, so it acts as the handler for all incoming requests with the URL path stated
	mux.Handle("/", middlewareOne(middlewareTwo(finalHandler)))
	// We can start visualize the middleware chain and what it will output depending in the order we nested them,
	// and then back up again in the reverse direction.
	log.Print("Listening on :3000...")
	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second parameter.
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)

}
