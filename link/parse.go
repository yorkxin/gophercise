package link

import "io"

// Link holds parsed data of <a href="...">text</a>
type Link struct {
	Href string
	Text string
}

// Parse reads an HTML document, then returns links in the doc.
func Parse(reader io.Reader) ([]Link, error) {
	// TODO: implementation
	return nil, nil
}
