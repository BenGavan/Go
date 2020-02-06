package main

import "fmt"

func main() {
	a := 4.639121096
	p := 4.639105182

	percentage := ((a - p) / a) * 100

	fmt.Println(percentage)
}
