package urlshort

import (
	"encoding/json"
	"net/http"

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
