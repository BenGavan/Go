package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hey")
	if err := run(); err != nil {
		fmt.Println("Error", err)
		return
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleMain)
	httpServer := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	err := httpServer.ListenAndServe()
	return err
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	type page struct {
		Title string
		Article template.HTML
	}

	var p page
	a, err := readFile("article1.html")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	p.Article = template.HTML(a)
	p.Title = "This is the Article title."
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(w, p)
}

func readFile(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	var bytes []byte
	var text = make([]byte, 1)
	for {
		_, err = file.Read(text)

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		bytes = append(bytes, text[0])
	}

	textString := string(bytes)
	return textString, nil
}