package main

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"io"
	"os"
)

func main() {
	fmt.Println("Hey")

	//line := "\\title{This is the Title}"
	//
	//command := extractCommandFromString(line)
	//value := extractValueFromCurlyBrackets(line)
	//fmt.Println(command, value)
	//fmt.Println(findIndexesOfOccurrencesOfByte("\\qwee{hi{", []byte("{")[0]))
	//
	//var lines []string
	//lines = append(lines, line)
	//
	//out := convertLATEXFileToHTML(lines)
	//fmt.Println(out)

	test()
}

func test() {
	filePath := "template_Article.tex"
	htmlFilePath := "out.html"
	//filePath := "test.txt"
	lines := openLATEXFile(filePath)
	htmlLines := convertLATEXFileToHTML(lines)
	saveHTMLLines(htmlFilePath, htmlLines)
}

func convertLATEXFileToHTML(lines []string) []string {
	var htmlLines []string
	
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line[0] != []byte("\\")[0] || (len(line) < 1) {
			continue
		}
		commandString := extractCommandFromString(line)
		switch commandString {
		case "title":
			val := extractValueFromCurlyBrackets(line)
			tag := "<h1>" + val + "</h1>"
			htmlLines = append(htmlLines, tag)
		}
	}

	return htmlLines
}

func openLATEXFile(filePath string) []string {

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	var lines []string
	var line []byte
	var linesBytes [][]byte

	var bytes []byte
	var text = make([]byte, 1)
	for {
		_, err = file.Read(text)

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(text[0])
		line = append(line, text[0])
		if len(line) > 1 {
			lastIndex := len(line) - 1

			if line[lastIndex] == 10 {
				fmt.Println(string(line))
				line = line[:lastIndex]

				if line[0] == 10 {
					line = line[1:]
				}

				lines = append(lines, string(line))
				linesBytes = append(linesBytes, line)
				line = nil
			}
		}

		bytes = append(bytes, text[0])
	}

	return lines
}

func saveHTMLLines(filepath string, lines []string) {
	var bytes []byte

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		var byteLine []byte

		for j := 0; j < len(line); j++ {
			byteLine = append(byteLine, line[j])
		}

		for j := 0; j < len(byteLine); j++ {
			bytes = append(bytes, byteLine[j])
		}
	}
	writeToFile(filepath, bytes)
}

func writeToFile(filepath string, bytes []byte) {
	fileCreated := createFile(filepath)
	if !fileCreated {
		fmt.Println("ERROR: File not created")
		return
	}

	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Error(err)
		return
	}

	_, err = file.Write(bytes)
	if err != nil {
		log.Error(err)
		return
	}

	err = file.Sync()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println("File successfully uploaded")
}

func createFile(filepath string) bool {
	if fileDoesExist(filepath) {
		fmt.Println("File", filepath, "does exist.")
		return true
	}

	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		fmt.Println("File has not been created")
		return false
	}

	fmt.Println("File", filepath, "has been created")
	return true
}

func fileDoesExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func removeNewLineByte(line []byte) []byte {
	lastIndex := len(line) - 1
	if line[lastIndex] == 10 {
		line = line[:lastIndex - 1]
		return removeNewLineByte(line)
	} else {
		return line
	}
}

func extractCommandFromString(s string) string {
	if s[0] != []byte("\\")[0] {
		return ""
	}

	endIndex := findIndexesOfOccurrencesOfByte(s, []byte("{")[0])[0]

	return s[1: endIndex]
}

func extractValueFromCurlyBrackets(s string) string {
	startIndex := findIndexesOfOccurrencesOfByte(s, charToByte("{"))[0] + 1
	endIndex := findIndexesOfOccurrencesOfByte(s, charToByte("}"))[0]
	return s[startIndex : endIndex]
}

func findIndexesOfOccurrencesOfByte(s string, b byte) []int {
	var indexes []int
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func charToByte(c string) byte {
	return []byte(c)[0]
}

