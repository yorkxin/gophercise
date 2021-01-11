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

	fmt.Println(*urlFlag)

	resp, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	base := resolveBaseURL(*resp.Request.URL)
	hrefs := extractHrefs(resp.Body, base)

	for _, href := range hrefs {
		fmt.Println(href)
	}
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
