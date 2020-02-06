package main

import (
	"fmt"
	"net/http"
)

type CustomHandler struct{}

func main() {
	fmt.Println("hi")
	registerRoutes()
	httpServer := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}


func (handler *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Custom Handler"))
}
