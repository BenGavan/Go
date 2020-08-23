package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type server struct {
	router *http.ServeMux
}

func newServer() *server {
	router := http.NewServeMux()
	s := &server{
		router: router,
	}
	s.setupRoutes()
	return s
}

func (s *server) setupRoutes() {
	s.router.Handle("/", s.handleHome())
}

func (s *server) handleHome() http.HandlerFunc {
	type response struct {
		Properties []Property `json:"properties"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.printInfo(w, r)
		resp := response{Properties: []Property{
			{
				Address: "Address One",
				Price:   123456.99,
				ImageURLs: []string{
					"",
					"",
				},
			},
			{
				Address: "Address Two",
				Price:   7654321.00,
			},
			{
				Address: "Address Three",
				Price:   19847289,
			},
			{
				Address: "Address Four",
				Price:   123,
			},
		}}
		json.NewEncoder(w).Encode(&resp)
	}
}

func (s *server) printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now(), "|", r.URL.Path, "|", r.Method, "|")
}

type Property struct {
	Address   string   `json:"address"`
	Price     float32  `json:"price"`
	ImageURLs []string `json:"image_urls"`
}

func run() error {
	s := newServer()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}
	err := httpServer.ListenAndServe()
	return err
}

func main() {
	fmt.Println("Starting Server")
	if err := run(); err != nil {
		panic(err)
	}
}
