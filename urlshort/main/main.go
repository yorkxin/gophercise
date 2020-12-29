package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yorkxin/Gophercise/urlshort"
)

func main() {
	yamlPathPtr := flag.String("yaml", "redirection.yml", "YAML file for redirection mapping")
	jsonPathPtr := flag.String("json", "redirection.json", "JSON file for redirection mapping")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := ioutil.ReadFile(*yamlPathPtr)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
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
