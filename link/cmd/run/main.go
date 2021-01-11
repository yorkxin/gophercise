package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yorkxin/Gophercise/link"
)

func main() {
	flag.Parse()
	filename := flag.Arg(0)

	if filename == "" {
		fmt.Fprintln(os.Stderr, "First argument must be a filename")
		flag.PrintDefaults()
		os.Exit(1)
	}

	reader, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	links, _ := link.Parse(reader)

	for _, link := range links {
		fmt.Printf("%+v\n", link)
	}
}
