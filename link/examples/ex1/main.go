package main

import (
	"fmt"
	"strings"

	"github.com/yorkxin/Gophersise/link"
)

var exampleHTML = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>`

func main() {
	reader := strings.NewReader(exampleHTML)
	links, _ := link.Parse(reader)
	fmt.Printf("%+v\n", links)
}
