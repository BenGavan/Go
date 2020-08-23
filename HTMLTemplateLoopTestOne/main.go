package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	fmt.Println("Hey")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleLayout)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func handleLayout(w http.ResponseWriter, r *http.Request) {
	type item struct  {
		PropertyOne string
		Number int
		Complete bool
	}
	type pageData struct {
		Title string
		Items []item
	}
	pd := pageData{
		Title:"This is the title.",
		Items: []item{{
			PropertyOne: "First",
			Number:      1,
			Complete:true,
		}, {
			PropertyOne: "Second",
			Number:      2,
			Complete:false,
		}},
	}
	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(w, pd)
}
