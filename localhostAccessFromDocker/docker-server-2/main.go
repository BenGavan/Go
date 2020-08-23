package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "http://docker.for.mac.localhost:8881"

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		printInfo(w, r)
		w.Write([]byte("Hello from docker server (from inside a docker container)"))
	})
	router.Handle("/docker-to-docker", handleRequest())

	httpServer := &http.Server{
		Addr:    ":8882",
		Handler: router,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Handles the response to get requests to baseURL/docker-to-docker by returning some dummy data.
func handleRequest() http.HandlerFunc {
	type Response struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		printInfo(w, r)

		// TODO: Make repsonse
		resp := Response{
			ID:      "id-string",
			Content: "Some content from 8882/docker-to-docker",
		}

		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v | %v | %v\n", time.Now(), r.URL.Path, r.Method)
}
