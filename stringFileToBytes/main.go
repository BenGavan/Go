package main

import (
	"fmt"
	"io"
	"os"
)

func run() {
	path := "r2.txt"
	bytes := readFile(path)
	fmt.Println(bytes)
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == byte('Z') {
			bytes[i] = byte('z')
		}
	}
	fmt.Println(string(bytes))
}

func main() {
	run()
}

func readFile(path string) []byte {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

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

		bytes = append(bytes, text[0])
	}
	return bytes
}