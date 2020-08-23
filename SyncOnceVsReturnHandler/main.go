package main

import (
	"fmt"
	"net/http"
	"sync"
)

type onceHandler struct {
	once sync.Once
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/return", handleReturn())
	mux.Handle("/once", &onceHandler{})

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	httpServer.ListenAndServe()
}

func handleReturn() http.HandlerFunc {
	fmt.Println("Handler Returning - Runs when server is initialized")
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Return: Sending data")
		w.Write([]byte("<h1>Returned Handler Function</h1>"))
	}
}

func (t *onceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		fmt.Println("Handler sync once - Runs only the first time this handler is required (When a user first requests it - NOT on start up)")
	})
	fmt.Println("once: Sending data")
	w.Write([]byte("<h1>Once</h1>"))
}
