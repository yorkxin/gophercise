package urlshort

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type redirectionEntry struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

const urlKeyRegex = "^[a-z0-9_-]{5,20}$"

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

func httpWriteError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func dbLookupRedirectKey(db *sql.DB, key string) (url *string, err error) {
	stmt, err := db.Prepare("select url from urls where key = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(key)
	err = row.Scan(&url)

	if err == sql.ErrNoRows {
		// no row, okay
		return nil, nil
	} else if err == nil {
		return url, err
	} else {
		// unhandled error from db
		return nil, err
	}
}

func dbInsertRedirection(db *sql.DB, key string, url string) (err error) {
	stmt, err := db.Prepare("insert into urls(key, url) values(?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(key, url); err != nil {
		return err
	} else {
		// OK
		return nil
	}
}

func isValidURL(rawURL string) bool {
	parsedURL, urlParseErr := url.Parse(rawURL)
	return urlParseErr == nil && parsedURL.IsAbs() == true
}

// DBRedirectHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to look up any
// paths (keys in the db) to their corresponding URL (values
// that each key in the db points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func DBRedirectHandler(db *sql.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/")
		url, err := dbLookupRedirectKey(db, key)

		if err != nil {
			httpWriteError(w, err)
			return
		} else if url == nil {
			log.Printf("No entry matches %q in db, fallback to default.", key)
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, *url, http.StatusFound)
		}
	}
}

func DBCreateHandler(db *sql.DB, w http.ResponseWriter, r http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot Prase Request Body"))
	}

	key := r.Form.Get("key")
	rawURL := r.Form.Get("url")

	existingUrl, err := dbLookupRedirectKey(db, key)
	if err != nil {
		httpWriteError(w, err)
		return
	}
	if existingUrl != nil {
		errorMessage := fmt.Sprintf("This key is unavailable for a new redirection: %q", key)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(errorMessage))
		return
	}

	if !(regexp.MustCompile(urlKeyRegex).MatchString(key) && isValidURL(rawURL)) {
		errorMessage := fmt.Sprintf("Invalid key or url. key must match %v and url must be absolute", urlKeyRegex)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(errorMessage))
		return
	}

	if err := dbInsertRedirection(db, key, rawURL); err != nil {
		httpWriteError(w, err)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))
	}
}
