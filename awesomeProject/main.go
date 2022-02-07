package main

import "fmt"

type testStruct struct {
	x func(string) bool
}

func main() {
	s := testStruct{}

	s.x = validate

	v := s.x("test string")
	fmt.Println(v)

	// or

	d := data{
		s: "test sting",
	}

	v = d.validate()
}

func validate(s string) bool {
	return len(s) > 0
}

type data struct {
	s string
}

func (d *data) validate() bool {
	return len(d.s) > 0
}

