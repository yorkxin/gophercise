package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yorkxin/Gophercise/urlshort"
)

func main() {
	dbPathPtr := flag.String("db", "db/db.sqlite3", "SQLite file for redirection mapping")
	yamlPathPtr := flag.String("yaml", "redirection.yml", "YAML file for redirection mapping")
	jsonPathPtr := flag.String("json", "redirection.json", "JSON file for redirection mapping")
	flag.Parse()

	mux := defaultMux()

	// open db
	dbPath, err := filepath.Abs(*dbPathPtr)
	log.Println(dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// POST /url
	mux.HandleFunc("/url", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			urlshort.DBCreateHandler(db, rw, *r)
			return
		} else {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Not Found"))
		}
	})

	dbHandler := urlshort.DBRedirectHandler(db, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := ioutil.ReadFile(*yamlPathPtr)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yaml, dbHandler)
	if err != nil {
		panic(err)
	}

	// reads from JSON file
	json, err := ioutil.ReadFile(*jsonPathPtr)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on 127.0.0.1:8080")
	http.ListenAndServe("127.0.0.1:8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
