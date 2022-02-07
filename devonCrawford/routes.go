package main

import "net/http"

func (s *server) routes() {

	s.router.Handle("/", s.handleHome())
	s.router.Handle("/physics", s.handlePhysics())
	//s.router.Handle("/software", s.handleSoftware())

	s.router.Handle("/static/", http.StripPrefix("/static/", s.fileServer))
}