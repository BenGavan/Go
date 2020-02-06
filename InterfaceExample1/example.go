package main

type writer interface {
	Write([]byte) (int, error)
}

func writeData(wr writer, data string) {
	wr.Write([]byte(data))
}
