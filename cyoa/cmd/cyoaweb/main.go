package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/yorkxin/Gophercise/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "port to listen on")
	filename := flag.String("file", "gopher.json", "File to read stories from")

	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.ParseStory(f)
	if err != nil {
		panic(err)
	}

	myTemplate := template.Must(template.New("myTemplate").Parse(myHTML))

	handler := cyoa.NewHandler(
		story,
		cyoa.WithTemplate(myTemplate),
		cyoa.WithPathFunction(myPathFunc),
	)
	address := fmt.Sprintf("localhost:%d", *port)
	fmt.Println("Starting server at:", address)
	log.Fatal(http.ListenAndServe(address, handler))
}

func myPathFunc(request *http.Request) string {
	path := strings.Trim(request.URL.Path, " ")

	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}

	// "/story/intro" => "intro"
	return path[len("/story/"):]
}

var myHTML = `
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
    <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
  </ul>
</body>
</html>`
