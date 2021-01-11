package cyoa

import (
	"encoding/json"
	"io"
)

// ParseStory parses the input IO as Story structure
func ParseStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// Story is a collection of chapters and their handler
type Story map[string]Chapter

// Chapter is a chapter that contains its title and texts,
// as well as where to go from this chap.
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option is a link to the next chapter
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"` // not to be confused with type Chapter
}
