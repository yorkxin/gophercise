package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fileNamePtr := flag.String("csv", "problems.csv", "input file name")

	flag.Parse()

	fullPath, err := filepath.Abs(*fileNamePtr)
	check(err)

	file, err := os.Open(fullPath)
	check(err)

	csvReader := csv.NewReader(file)

	var questions, correct int

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		check(err)

		// quiz you
		questions++

		quiz := record[0]
		ans := record[1]

		fmt.Printf("%s >> ", quiz)

		inputReader := bufio.NewScanner(os.Stdin)

		if inputReader.Scan() {
			input := inputReader.Text()

			if input == ans {
				correct++
			}
		}
	}

	fmt.Println("You got", correct, "correct of total", questions, "questions.")
}
