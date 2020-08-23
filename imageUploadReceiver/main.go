package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Server struct {
	router *http.ServeMux
}

func newServer() *Server {
	router := http.NewServeMux()

	s := &Server{router: router}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Handle("/upload", s.handleUpload())
}

func (s *Server) printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v | %v | %v\n", time.Now(), r.URL.Path, r.Method)
}

func (s *Server) writeResponse(w http.ResponseWriter, jsonData interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(jsonData)
	if err != nil {
		http.Error(w, "ERROR: Failed to encode and write data", http.StatusInternalServerError)
	}
}

func (s *Server) handleUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.printInfo(w, r)
		//w.Write([]byte("Hey from upload"))

		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		form := r.Form["file"]
		fmt.Printf("Form: %v\n", form)

		s := r.FormValue("meta-json")
		fmt.Printf("key = meta-json, value = %v\n", s)

		v := r.Form["meta-json"]
		fmt.Printf("key = meta-json, value (v) = %v\n", v)

		//var buffer bytes.Buffer

		fmt.Printf("Here 1\n")

		filename := "image.jpg"
		file, fileHeader, err := r.FormFile(filename)
		//fmt.Printf("File header: %v\n", fileHeader)
		fmt.Printf("File: %v\n", file)
		if err != nil {
			fmt.Printf("Here 2x\n")
			fmt.Printf("Error: %v\n", err)
			http.Error(w, "error reading file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fmt.Printf("Here 3\n")
		//fmt.Printf("fileHeader: %v\n", string(fileHeader))
		fmt.Printf("file: %v\n", fileHeader.Filename)
		_ = fileHeader


		key := "meta-json"
		getS := r.Form.Get(key)
		fmt.Printf("get key = %v, value = %v\n", key, getS)

		for key, value := range r.Form {
			fmt.Printf("key: %v, value: %v\n", key, value)
		}

		w.Write([]byte("it must have worked"))

	}
}

func run() error {
	s := newServer()

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	err := httpServer.ListenAndServe()
	return err
}

func main() {
	fmt.Println("Starting Image upload server")

	rand.Seed(time.Now().UTC().UnixNano())

	if err := run(); err != nil {
		fmt.Printf("auth Server failed to start: %v\n", err)
	}
}
