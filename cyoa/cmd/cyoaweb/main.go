package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yorkxin/Gophercise/cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "File to read stories from")

	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	d := json.NewDecoder(f)
	var story cyoa.Story
	if err := d.Decode(&story); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)
}
