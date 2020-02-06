package main

import "net/http"

var mux = http.NewServeMux()

func registerRoutes() {
	mux.HandleFunc("/home", handleHome)
	mux.HandleFunc("/about",  handleAbout)
	mux.HandleFunc("/login",  handleLogin)
	mux.HandleFunc("/logout",  handleLogout)
}
