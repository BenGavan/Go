package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Printf("Hey")
}

// Good way of doing it is:
type writer interface {
	Write([]byte) (int, error)
}

func writeData(wr writer, data string) {
	wr.Write([]byte(data))
}


// Bad Way of doing it is (implementing different functions for slightly different tasks):
func writeDataToFile(f *os.File, data string){
	f.Write([]byte(data))
}

func writeToTCPConnection(con *net.TCPConn, data string) {
	con.Write([]byte(data ))
}

