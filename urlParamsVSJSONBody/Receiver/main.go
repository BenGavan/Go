package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	ID        string `json:"id"`
	TimeStamp int64  `json:"time_stamp"`
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/json", handleJSONBody())
	router.HandleFunc("/params", handleParams())

	server := http.Server{
		Addr:    ":8010",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func handleJSONBody() http.HandlerFunc {
	type Request struct {
		ID string `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		requestBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		defer r.Body.Close()

		var request Request
		err = json.Unmarshal(requestBodyBytes, &request)

		err = makeStandardResponse(w, request.ID)
		if err != nil {
			panic(err)
		}
	}
}

func handleParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParams := r.URL.Query()

		id := urlParams.Get("id")

		err := makeStandardResponse(w, id)
		if err != nil {
			panic(err)
		}
	}
}

func makeStandardResponse(w http.ResponseWriter, id string) error {
	resp := Response{
		ID:        id,
		TimeStamp: time.Now().UTC().UnixNano(),
	}

	bs, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(bs)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	return err
}
