package main

import "net/http"

func (s *server) routes() {

	s.router.Handle("/", s.handleSpecificArticle())
	s.router.Handle("/software", s.handleSoftware())

	s.router.Handle("/static/", http.StripPrefix("/static/", s.fileServer))
}