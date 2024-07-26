package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Create a new application type
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	//add a new flag named addr with default port number 4000 to change at run time
	addr := flag.String("addr", ":4000", "HTTP network addr")

	//parse the command-line flags from os.Args. Must be call after all flags are defined
	//and before flags are accessed by the program
	flag.Parse()

	//create new logger with log.New(out io.Writer,prefix String, flags int)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Create a new *http.ServeMux which act as a router
	mux := http.NewServeMux()

	//Initialize a new instance of application named app
	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
	}
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  mux,
	}
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/crate", app.snippetCreate)

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}
