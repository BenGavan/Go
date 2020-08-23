package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	PHYSICS = "physicsbooklist.txt"
	CS      = "csbooklist.txt"
	MATH    = "MathBookList.txt"
)

func main() {
	fmt.Println("Hey")

	books := searchForString("Basic", MATH)
	for i := 0; i < len(books); i++ {
		fmt.Println(books[i])
	}
}

func countLinesForFile(path string) int {
	lines := openFile(MATH)
	return len(lines)
}

func openFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return lines
}

func searchForString(s, path string) []string {
	books := openFile(path)
	var booksMatched []string

	for i := 0; i < len(books); i++ {
		book := books[i]
		if strings.Contains(book, s) {
			booksMatched = append(booksMatched, book)
		}
	}
	return booksMatched
}
