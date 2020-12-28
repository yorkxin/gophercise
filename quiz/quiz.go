package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fileNamePtr := flag.String("csv", "problems.csv", "input file name")
	limitPtr := flag.Uint("limit", 30, "time limit to answer one question")

	flag.Parse()

	fullPath, err := filepath.Abs(*fileNamePtr)
	check(err)

	file, err := os.Open(fullPath)
	check(err)

	fmt.Print("Ready? [y] ")
	inputReader := bufio.NewScanner(os.Stdin)
	for inputReader.Scan() {
		if inputReader.Text() == "y" {
			break
		} else {
			fmt.Print("Ready? [y] ")
		}
	}

	csvReader := csv.NewReader(file)
	problems, err := csvReader.ReadAll()
	check(err)

	timer := time.NewTimer(time.Duration(*limitPtr) * time.Second)

	var correct int

problemsLoop:
	for _, record := range problems {
		// quiz you
		quiz := record[0]
		ans := record[1]

		fmt.Printf("%s >> ", quiz)

		inputCh := make(chan string, 1)

		go func() {
			if inputReader.Scan() {
				inputCh <- inputReader.Text()
			}
		}()

		select {
		case input := <-inputCh:
			if input == ans {
				correct++
			}
		case <-timer.C:
			// timeout
			fmt.Println()
			break problemsLoop
		}
	}

	fmt.Println("You got", correct, "correct of total", len(problems), "questions.")
}
