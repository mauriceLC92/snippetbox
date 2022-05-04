package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mauriceLC92/snippetbox/internal/models"
)

type config struct {
	httpAddr string
}

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger

	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	var conf config
	// addr := flag.String("addr", ":4000", "HTTP network address")
	flag.StringVar(&conf.httpAddr, "addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://web:password@localhost/snippetbox?sslmode=disable", "PostgreSQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     conf.httpAddr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s \n", conf.httpAddr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
