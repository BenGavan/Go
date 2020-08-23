package main

// Keep routes in a single file - in one place - so can basically see the map of the server in one place (it's glancible)
// Most code maintenance starts with a URL so it's handy to have one place to look.
/*
routes.go is the high-level map of the routes
 */
func (s *server) routes() {
	// set up routes (mux.HandleFunc(...))
	//s.router.HandleFunc()
	s.router.HandleFunc("/something/else", s.handleSomethingElse("Hello %s"))
	s.router.HandleFunc("/admin", s.adminOnly(s.handleSomething()))
}