package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/yorkxin/Gophercise/link"
)

const sitemapNS = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	URL string `xml:"loc"`
}

type urlset struct {
	Urlset []loc  `xml:"url"`
	Xmlns  string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the URL you want to build sitemap for")
	depth := flag.Int("depth", 3, "the depth of URLs to visit")
	debugMode := flag.Bool("debug", false, "set to true for debug mode")

	flag.Parse()

	debugOutput := ioutil.Discard
	if *debugMode == true {
		debugOutput = os.Stderr
	}

	hrefs := bfs(*urlFlag, *depth, &debugOutput)

	sitemap := urlset{
		Urlset: make([]loc, len(hrefs)),
		Xmlns:  sitemapNS,
	}

	for i, href := range hrefs {
		sitemap.Urlset[i] = loc{URL: href.url}
	}

	output := os.Stdout

	output.WriteString(xml.Header)
	encoder := xml.NewEncoder(output)
	encoder.Indent("", "  ")
	if err := encoder.Encode(sitemap); err != nil {
		panic(err)
	}
	output.WriteString("\n")
}

type visitMeta struct {
	url   string
	depth int // depth of first discovery
}

func bfs(urlToAccess string, maxDepth int, debugOutput *io.Writer) []visitMeta {
	visited := make(map[string]int) // number is depth

	nextVisit := []string{urlToAccess}

	for depth := 0; depth <= maxDepth; depth++ {
		toVisit := nextVisit
		nextVisit = make([]string, 0)

		for _, visitURL := range toVisit {
			fmt.Fprintf(*debugOutput, "[D=%d] ", depth)

			if _, ok := visited[visitURL]; ok == true {
				fmt.Fprintf(*debugOutput, "\x1b[34mskip\x1b[m : %s\n", visitURL)
				continue
			}
			fmt.Fprintf(*debugOutput, "\x1b[1;32mvisit\x1b[m: %s\n", visitURL)
			visited[visitURL] = depth

			for _, newURL := range getHrefsFromURL(visitURL) {
				if _, ok := visited[newURL]; ok == false {
					fmt.Fprintf(*debugOutput, "           + %s\n", newURL)
					nextVisit = append(nextVisit, newURL)
				} else {
					fmt.Fprintf(*debugOutput, "           - %s\n", newURL)
				}
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
