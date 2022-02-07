package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type server struct {
	router      *http.ServeMux
	fileServer  http.Handler
	baseURL     string
	articlesURL string
}

func newServer() *server {
	router := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("static/"))
	baseURL := "http://localhost:8080"
	s := &server{
		router:      router,
		fileServer:  fileServer,
		baseURL:     baseURL,
		articlesURL: baseURL + "/articles",
	}
	s.routes()
	return s
}

func (s *server) printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now(), "|", r.URL.Path, "|", r.Method, "|")
}

func (s *server) handleHome() http.HandlerFunc {
	var (
		init     sync.Once
		tmplt    *template.Template
		tmpltErr error
	)

	type Page struct {

	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.printInfo(w, r)
		init.Do(func() {
			tmplt, tmpltErr = template.ParseFiles("src/index.html")
		})
		if tmpltErr != nil {
			http.Error(w, tmpltErr.Error(), http.StatusInternalServerError)
			return
		}

		err := tmplt.Execute(w, nil)
		if err != nil {
			http.Error(w, tmpltErr.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handlePhysics() http.HandlerFunc {
	var (
		init sync.Once
		tmplt *template.Template
		templateErr error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		s.printInfo(w, r)
		init.Do(func() {
			tmplt, templateErr = template.ParseFiles("src/physics.html")
		})
		if templateErr != nil {
			http.Error(w, templateErr.Error(), http.StatusInternalServerError)
			return
		}

		err := tmplt.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func run() error {
	s := newServer()
	httpServer := &http.Server{
		Addr:    ":80",
		Handler: s.router,
	}
	err := httpServer.ListenAndServe()
	return err
}

func main() {
	fmt.Println("Hey")
	err := run()
	if err != nil {
		panic(err)
	}
}
