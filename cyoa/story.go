package cyoa

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
