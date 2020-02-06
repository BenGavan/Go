package main

import "net/http"

func handleAbout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About route"))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About home"))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About login"))
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About logout"))
}
