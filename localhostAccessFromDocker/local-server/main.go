package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		printInfo(w, r)
		w.Write([]byte("Hello from local server (not inside a docker container)"))
	})
	router.HandleFunc("/cl", func(w http.ResponseWriter, r *http.Request) {
		printInfo(w, r)
		type Response struct {
			Test string `json:"test"`
		}
		resp := Response{Test: "Hello from local server (not inside a docker container)"}
		json.NewEncoder(w).Encode(resp)
	})
	httpServer := &http.Server{
		Addr:    ":8880",
		Handler: router,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v | %v | %v\n", time.Now(), r.URL.Path, r.Method)
}

