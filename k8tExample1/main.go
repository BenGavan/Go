package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting test service")

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		printInfo(w, r)
		w.Write([]byte("This si the test service router = /"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v | %v | %v\n", time.Now(), r.URL.Path, r.Method)
}
