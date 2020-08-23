package main

import (
	"fmt"
	"os"
)

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	cmd := 
}

// docker run <container> cmd args
// vs
// go run main.go run cmd args
func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("What??")
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
