package main

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// Represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	if err := http.ListenAndServe(":8080", nil);  err != nil {
		log.Fatal("ListenAndServe", err)
	}
}





func printInfo(r *http.Request) {
	fmt.Println(time.Now(), "|", r.URL.Path, "|", r.Method, "|")
}
