package main

import (
	"../utils"
	"fmt"
)


//func handleHome(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("Hey there"))
//}
//
//func main() {
//	http.HandleFunc("/", handleHome)
//	http.ListenAndServe(":8080", nil)
//}

func main() {
	utils.PrintTest()
	fmt.Printf("Is is successful")
}

//
//func main() {
//	fmt.Println("Hey")
//	switch os.Args[1] {
//	case "run":
//		run()
//	default:
//		panic("what??")
//	}
//}
//
//func run() {
//	fmt.Println("running", os.Args[2:])
//
//	cmd := exec.Command(os.Args[2], os.Args[3:]...)
//	cmd.Stdin = os.Stdin
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//
//	must(cmd.Run())
//}
//
//func must(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
