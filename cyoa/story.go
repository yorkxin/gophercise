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

var defaultTemplate *template.Template

func init() {
	defaultTemplate = template.Must(template.New("").Parse(html))
}

type storyHandler struct {
	story    Story
	template *template.Template
	pathFn   func(*http.Request) string
}

func (h storyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	storyName := h.pathFn(request)
	chapter, found := h.story[storyName]

	if !found {
		http.Error(writer, "Chapter not found.", http.StatusNotFound)
		return
	}

	if err := h.template.Execute(writer, chapter); err != nil {
		http.Error(writer, "Something went wrong...", http.StatusInternalServerError)
	}

	// success
}

func defaultPathFn(request *http.Request) string {
	path := strings.Trim(request.URL.Path, " ")

	if path == "" || path == "/" {
		return "intro"
	}

	// "/intro" => "intro"
	return path[1:]
}

// HandlerOption is a return type for functional options.
//
// NOTE: Expect the *storyHandler will be modified by functional options in place.
type HandlerOption func(h *storyHandler)

// WithTemplate allows the user to overide template
func WithTemplate(myTemplate *template.Template) HandlerOption {
	return func(h *storyHandler) {
		h.template = myTemplate
	}
}

// WithPathFunction allows the user to overide path parsing function
func WithPathFunction(myFunction func(*http.Request) string) HandlerOption {
	return func(h *storyHandler) {
		h.pathFn = myFunction
	}
}

// NewHandler returns an http.Handler for story navigation
func NewHandler(story Story, options ...HandlerOption) http.Handler {
	handler := storyHandler{
		story:    story,
		template: defaultTemplate,
		pathFn:   defaultPathFn,
	}

	// "functional options" pattern
	// See https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	for _, option := range options {
		option(&handler)
	}

	return handler
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
