package link

import (
	"io"
	"strings"

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

	var links []Link

	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func buildLink(node *html.Node) Link {
	var link Link

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break // ignore duplicate ones
		}
	}

	link.Text = strings.Join(strings.Fields(extractText(node)), " ")
	return link
}

func extractText(node *html.Node) string {
	switch node.Type {
	case html.TextNode:
		return node.Data
	case html.ElementNode:
		var fullText string
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			fullText += extractText(child) + " "
		}
		return fullText
	default:
		return ""
	}
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
