package link

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// Link holds parsed data of <a href="...">text</a>
type Link struct {
	Href string
	Text string
}

// Parse reads an HTML document, then returns links in the doc.
func Parse(reader io.Reader) ([]Link, error) {
	doc, err := html.Parse(reader)

	if err != nil {
		return nil, err
	}

	nodes := linkNodes(doc)

	for _, node := range nodes {
		fmt.Printf("%+v\n", node)
	}
	return nil, nil
}

// returns all <a> nodes under this node. If node itself is an <a>, it'll be
// wrapped in a slice.
func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}

	var childLinks []*html.Node

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		childLinks = append(childLinks, linkNodes(child)...)
	}

	return childLinks
}
