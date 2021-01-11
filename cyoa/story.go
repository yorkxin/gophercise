package cyoa

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"text/template"
)

var html = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Choose Your Own Adventure!</title>
</head>
<body>
  <h1>{{.Title}}</h1>

  <article>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
  </article>

  <ul>
    {{range .Options}}
    <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
  </ul>
</body>
</html>`

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").Parse(html))
}

type storyHandler struct {
	story Story
}

func (h storyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var storyName string
	path := strings.Trim(request.URL.Path, " ")

	if path == "" || path == "/" {
		storyName = "intro"
	} else {
		// "/intro" => "intro"
		storyName = path[1:]
	}

	chapter, found := h.story[storyName]

	if !found {
		http.Error(writer, "Chapter not found.", http.StatusNotFound)
		return
	}

	if err := tmpl.Execute(writer, chapter); err != nil {
		http.Error(writer, "Something went wrong...", http.StatusInternalServerError)
	}

	// success
}

// NewHandler returns an http.Handler for story navigation
func NewHandler(story Story) http.Handler {
	return storyHandler{story: story}
}

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
