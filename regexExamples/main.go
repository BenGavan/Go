package main

import (
	"fmt"
	"regexp"
)

func main() {
	matched, err := regexp.MatchString(`\s.\s.\S.\S`, "a a x b b")
	fmt.Println(matched) // true
	fmt.Println(err)     // nil (regexp is valid)
}
