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

	dfs(doc, "")
	return nil, nil
}

func dfs(node *html.Node, padding string) {
	msg := node.Data
	if node.Type == html.ElementNode {
		msg = "<" + msg + ">"
	}

	fmt.Println(padding, msg)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		dfs(child, padding+"  ")
	}
}
