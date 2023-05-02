package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

//the LoggingHandler middleware from the gorilla/handlers package, which records request logs using the Apache Common Log Format.

//Instead of using the standard middleware signature, this middleware has the signature func(out io.Writer, h http.Handler)
// http.Handler, so it takes not only the next handler but also the io.Writer that the log will be written to.

func newLoggingHandler(dst io.Writer) func(http.Handler) http.Handler {

	//constructor function which wraps the LoggingHandler() middleware and returns a standard func(http.Handler) http.Handler
	// function that we can nest neatly with other middleware.
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(dst, h)
	}
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)

	/*
		Creates a new file called "server.log" or opens an existing file with the same name if it exists.
		The file permission bits 0664 specify read and write permission for the owner and group, and read permission for others.

		The file is opened in write-only mode and with the O_CREATE and O_APPEND flags set.
		The O_CREATE flag specifies that the file should be created if it does not already exist, and the O_APPEND flag specifies that data
		should be appended to the end of the file rather than overwriting its contents.
	*/

	if err != nil {
		log.Fatal(err) //If an error occurs while opening or creating the file, the program will log the error and terminate.
	}

	//creating a new logging handler using the newLoggingHandler function, passing in the logFile as the destination for the log output.
	loggingHandler := newLoggingHandler(logFile)

	// Use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()

	//HandlerFunc serves as an adapter to allow the use of ordinary http.Handlers, in this case used with Final\
	finalHandler := http.HandlerFunc(final)

	//function to register this with our new servemux, so it acts as the handler for all incoming requests with the URL path stated
	mux.Handle("/", loggingHandler(finalHandler))

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", mux)
	// Then we create a new server and start listening for incoming requests
	// with the http.ListenAndServe() function, passing in our servemux for it to
	// match requests against as the second parameter.
	log.Fatal(err)
}
