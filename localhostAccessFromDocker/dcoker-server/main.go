package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const baseURL = "http://docker.for.mac.localhost:8882"

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		printInfo(w, r)
		w.Write([]byte("Hello from docker server (from inside a docker container)"))
	})
	// Request content from a localhost server
	router.HandleFunc("/cl", func(w http.ResponseWriter, r *http.Request) {
		printInfo(w, r)

		type Response struct {
			This string `json:"this"`
			Test string `json:"test"`
		}

		url := "http://docker.for.mac.localhost:8000/cl"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error Getting url: %v: %v\n", url, err)
		}
		fmt.Printf("response from %v: %v\n", url, resp)

		defer resp.Body.Close()

		// Get response data from body
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		var response Response
		err = json.Unmarshal(responseBody, &response)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		response.This = "This is the docker container speeking"

		json.NewEncoder(w).Encode(response)
	})
	router.HandleFunc("/test", handleContactDocker())

	httpServer := &http.Server{
		Addr:    ":8881",
		Handler: router,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Makes a GET request to to 8002
func handleContactDocker() http.HandlerFunc {
	type Response struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		printInfo(w, r)

		url := baseURL + "/docker-to-docker"

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		defer resp.Body.Close()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		var response Response
		err = json.Unmarshal(responseBody, &response)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		jsonBytes, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Response: %v\n", string(jsonBytes))
		return
	}
}

func printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v | %v | %v\n", time.Now(), r.URL.Path, r.Method)
}
