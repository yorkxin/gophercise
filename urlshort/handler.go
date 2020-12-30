package urlshort

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type redirectionEntry struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if target, ok := pathsToUrls[r.URL.Path]; ok == true {
			http.Redirect(w, r, target, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var list []redirectionEntry
	if err := yaml.Unmarshal(yml, &list); err != nil {
		return fallback.ServeHTTP, err
	}

	laMap := map[string]string{}

	for _, entry := range list {
		laMap[entry.Path] = entry.URL
	}

	return MapHandler(laMap, fallback), nil
}

// JSONHandler serves the same purpose of YAMLHandler, but reads from a JSON payload
func JSONHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var list []redirectionEntry
	if err := json.Unmarshal(yml, &list); err != nil {
		return fallback.ServeHTTP, err
	}

	laMap := map[string]string{}

	for _, entry := range list {
		laMap[entry.Path] = entry.URL
	}

	return MapHandler(laMap, fallback), nil
}

// DBHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to look up any
// paths (keys in the db) to their corresponding URL (values
// that each key in the db points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func DBHandler(db *sql.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stmt, err := db.Prepare("select url from urls where key = ?")

		httpWriteError := func(err error) {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		if err != nil {
			httpWriteError(err)
			return
		}

		defer stmt.Close()

		name := strings.TrimPrefix(r.URL.Path, "/")
		row := stmt.QueryRow(name)

		if err != nil {
			httpWriteError(err)
			return
		}

		var url string
		err = row.Scan(&url)

		if err == sql.ErrNoRows {
			log.Printf("No entry matches %q in db, fallback to default.", name)
			fallback.ServeHTTP(w, r)
		} else if err == nil {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			httpWriteError(err)
		}
	}
}
