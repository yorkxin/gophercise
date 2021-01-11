package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/yorkxin/Gophercise/link"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the URL you want to build sitemap for")

	hrefs := bfs(*urlFlag, 5)
	for _, href := range hrefs {
		fmt.Printf("%+v\n", href)
	}
}

type visitMeta struct {
	url   string
	depth int // depth of first discovery
}

func bfs(urlToAccess string, maxDepth int) []visitMeta {
	visited := make(map[string]int) // number is depth

	nextVisit := []string{urlToAccess}

	for depth := 0; depth <= maxDepth; depth++ {
		toVisit := nextVisit
		nextVisit = make([]string, 0)

		for _, visitURL := range toVisit {
			fmt.Printf("[D=%d] ", depth)

			if _, ok := visited[visitURL]; ok == true {
				fmt.Printf("\x1b[34mskip\x1b[m : %s\n", visitURL)
				continue
			}
			fmt.Printf("\x1b[1;32mvisit\x1b[m: %s\n", visitURL)
			visited[visitURL] = depth

			nextVisit = append(nextVisit, getHrefsFromURL(visitURL)...)

			for _, newURL := range nextVisit {
				fmt.Printf("   + %s\n", newURL)
			}
		}
	}

	// transform map to array
	allVisited := make([]visitMeta, 0, len(visited))
	for visitedURL, depth := range visited {
		allVisited = append(allVisited, visitMeta{url: visitedURL, depth: depth})
	}

	return allVisited
}

func getHrefsFromURL(urlToAccess string) []string {
	resp, err := http.Get(urlToAccess)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	base := resolveBaseURL(*resp.Request.URL)
	return filter(extractHrefs(resp.Body, base), func(str string) bool {
		return strings.HasPrefix(str, base)
	})
}

func resolveBaseURL(reqURL url.URL) string {
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	return baseURL.String()
}

func extractHrefs(reader io.Reader, base string) []string {
	links, _ := link.Parse(reader)

	var hrefs []string
	for _, aLink := range links {
		if strings.HasPrefix(aLink.Href, "/") {
			hrefs = append(hrefs, base+aLink.Href)
		} else if strings.HasPrefix(aLink.Href, "http") {
			hrefs = append(hrefs, aLink.Href)
		}
	}

	return hrefs
}

func filter(inputStrings []string, keepFunc func(string) bool) []string {
	var result []string
	for _, href := range inputStrings {
		if keepFunc(href) {
			result = append(result, href)
		}
	}
	return result
}
