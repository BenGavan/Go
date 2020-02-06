package main

import "fmt"

func main() {
	r, e := testReturn("Hey :)")
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(r)
}

func testReturn(test string) (result string, err error) {
	result = "We recieved a string = " + test
	err = nil
	return
}


func testArray() {
	fmt.Printf("hello, world\n")

	var twoDArray [][]byte
	var subArray []byte


	for i := 0; i < 100; i++ {
		subArray = append(subArray, byte(i))
		twoDArray = append(twoDArray, subArray)
	}

	fmt.Println(twoDArray)
	fmt.Println(subArray)
}
