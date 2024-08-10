package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	//add a new flag named addr with default port number 4000 to change at run time
	addr := flag.String("addr", ":4000", "HTTP network addr")

	//parse the command-line flags from os.Args. Must be call after all flags are defined
	//and before flags are accessed by the program
	flag.Parse()

	//create new logger(which define how logs should be print) with log.New(out io.Writer,prefix String, flags int)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Initialize a new instance of application named app
	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
	}

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}
