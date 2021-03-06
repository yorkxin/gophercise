package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
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

	inputReader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Ready? [y] ")
		if inputReader.Scan(); inputReader.Text() == "y" {
			break
		}
	}

	csvReader := csv.NewReader(file)
	problems, err := csvReader.ReadAll()
	check(err)

	rand.
		New(rand.NewSource(time.Now().UnixNano())).
		Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})

	timer := time.NewTimer(time.Duration(*limitPtr) * time.Second)

	var correct int

problemsLoop:
	for i, record := range problems {
		// quiz you
		quiz := record[0]
		ans := record[1]

		inputCh := make(chan string, 1)

		go func(quizId int, quiz string) {
			for {
				fmt.Printf("Quiz %d: %s >> ", quizId, quiz)
				if inputReader.Scan() {
					if input := strings.Trim(inputReader.Text(), " "); input != "" {
						inputCh <- input
						return
					}
				}
			}
		}(i+1, quiz)

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
