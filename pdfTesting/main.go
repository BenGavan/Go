package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"rsc.io/pdf"
)

var text = ""

func do() {
	filepath := "test.pdf"
	fileReader, err := pdf.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	outline := fileReader.Outline()
	title := outline.Title
	outlineChild := outline.Child

	fmt.Printf("***\nOutline: %v\nOutline title: %v\nOutilne child: %v\n***\n", outline, title, outlineChild)

	numPage := fileReader.NumPage()

	fmt.Printf("***\nNum page: %v\n***", numPage)

	var previousY float64 = 0
	for i := 1; i < numPage; i++ {
		page := fileReader.Page(i)

		pageContent := page.Content()
		pageContentRect := pageContent.Rect
		pageContentText := pageContent.Text
		for i, c := range pageContentText {
			fmt.Printf("Text: %v, Byte: %v, Font: %v, Font Size: %v, width: %v, x: %v, y: %v\n", c.S, []byte(c.S), c.Font, c.FontSize, c.W, c.X, c.Y)
			if (previousY != c.Y) && (c.Y != 0) {
				text += "\n"
			}

			if i != 0 {
				if pageContentText[i-1]  {

				}
			}


			text += c.S
			previousY = c.Y
		}
		pageResourse := page.Resources()

		fmt.Printf("***\nPage: %v\n\nPage Conent: %v\n\nPage Resources: %v\n\nPage Conent Rect: %v\n\nPage Content Text: %v\n\n***\n", page, pageContent, pageResourse, pageContentRect, pageContentText)
	}

	trailer := fileReader.Trailer()
	trailerString := trailer.String()
	trailerKeys := trailer.Keys()

	fmt.Printf("***\nTrailer: %v\nTrailer String: %v\nTrailer Keys: %v\n***\n", trailer, trailerString, trailerKeys)

	fmt.Printf("String: %v\n", text)
}

func readBytesFromFile(filepath string) ([]byte, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func readBytes() {
	fmt.Println("Fuck you")

	filepath := "test.pdf"
	fileBytes, err := readBytesFromFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	fileAsString := string(fileBytes)

	fmt.Printf("File as String: %v\n", fileAsString)
}


func run() error {
	do()
	return nil
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}