package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type subject struct {
	Title    string    `json:"title"`
	Articles []article `json:"articles"`
}

type article struct {
	Title     string    `json:"title"`
	Date      string    `json:"date"`
	VideoURL  string    `json:"video_url"`
	SubjectID string    `json:"subject_id"`
	Elements  []element `json:"elements"`
}

type element struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

func createArticle() {
	_ = article{
		Title:    "title here",
		Date:     "date",
		Elements: nil,
	}
}

func main() {
	fmt.Println("hey")
	getData()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/json", handleJSON)
	mux.HandleFunc("/interpreted", handleInterpreted)
	httpServer := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	err := httpServer.ListenAndServe()
	return err
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	els := []element{
		{
			Tag:     "p",
			Content: "first",
		}, {
			Tag:     "p",
			Content: "second",
		},
	}
	a1 := article{
		Title:    "title here",
		Date:     "date",
		VideoURL: "http://www.youtube.com/watch...",
		Elements: els,
	}
	a2 := article{
		Title:    "Article 2",
		Date:     "Second date",
		VideoURL: "second url",
		Elements: []element{
			{
				Tag:     "p",
				Content: "This si the second article in this subject",
			},
		},
	}
	s := subject{
		Title: "First Subject",
		Articles: []article{
			a1,
			a2,
		},
	}
	json.NewEncoder(w).Encode(&s)
}

func handleInterpreted(w http.ResponseWriter, r *http.Request) {
	a := getData()

	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(w, a)
	if err != nil {
		fmt.Println("ERROR: Index template could not be executed.")
	}
}

func getData() article {
	plan, _ := ioutil.ReadFile("data.json")
	var data article
	err := json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data.Title)
	return data
}
