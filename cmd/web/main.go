package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.formme.net/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func newTemplateCache() (map[string]*template.Template, error) {
	//Create an empty map acting as a cache that contain the key-value pair
	//(key as the template name,value as a Template corresponding to the key)
	cache := map[string]*template.Template{}
	//Get all the pages in the pages folder
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	//Looping through the all the pages and create a Template corresponding
	//to the name of the page
	for _, page := range pages {
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		//ParseFiles created a new Template and parses
		//the template definitions from the named file
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		name := filepath.Base(page)
		cache[name] = ts
	}
	return cache, nil
}
func main() {
	//add a new flag named addr with default port number 4000 to change at run time
	addr := flag.String("addr", ":4000", "HTTP network addr")
	//add a new flag named DSN(Data source name)
	dsn := flag.String("dsn", "web:giahao45@/snippetbox?parseTime=true", "MySQL data source name")
	//parse the command-line flags from os.Args. Must be call after all flags are defined
	//and before flags are accessed by the program
	flag.Parse()

	//create new logger(which define how logs should be print) with log.New(out io.Writer,prefix String, flags int)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//open a new connection pool to database
	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()
	//Initialize a new template cache
	tempalteCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}
	//Initialize a new instance of application named app
	app := &application{
		infoLog:       infoLog,
		errorLog:      errLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: tempalteCache,
	}

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}
