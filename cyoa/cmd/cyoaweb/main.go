package main

import (
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

	story, err := cyoa.ParseStory(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)
}
