package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

	handler := cyoa.NewHandler(story)
	address := fmt.Sprintf("localhost:%d", *port)
	fmt.Println("Starting server at:", address)
	log.Fatal(http.ListenAndServe(address, handler))
}
