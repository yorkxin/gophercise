package main

import (
	"flag"
	"fmt"
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

	links, _ := link.Parse(resp.Body)
	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	var hrefs []string
	for _, aLink := range links {
		if strings.HasPrefix(aLink.Href, "/") {
			hrefs = append(hrefs, base+aLink.Href)
		} else if strings.HasPrefix(aLink.Href, "http") {
			hrefs = append(hrefs, aLink.Href)
		}
	}

	for _, href := range hrefs {
		fmt.Println(href)
	}
}
