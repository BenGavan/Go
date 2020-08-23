package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("iOS Cookie Test")

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/create", handleCreateCookie)
	http.HandleFunc("/get", handleGetCookie)
	http.HandleFunc("/delete", handleDeleteCookie)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	printInfo(w, r)
	//expire := time.Now().Add(20 * time.Minute) // Expires in 20 minutes
	//cookie := http.Cookie{Name: "username", Value: "nonsecureuser", Path: "/", Expires: expire, MaxAge: 86400}
	//http.SetCookie(w, &cookie)
	//cookie = http.Cookie{Name: "secureusername", Value: "secureuser", Path: "/", Expires: expire, MaxAge: 86400, HttpOnly: true, Secure: true}
	//http.SetCookie(w, &cookie)

	expires := time.Now().AddDate(1, 0, 0)

	ck := http.Cookie{
		Name:    "JSESSION_ID",
		Domain:  "localhost",
		Path:    "/",
		Expires: expires,
	}

	// value of cookie
	ck.Value = "value of this awesome cookie"
	ck.HttpOnly = false

	// write the cookie to response
	http.SetCookie(w, &ck)

	w.Write([]byte("Home"))
}

func handleCreateCookie(w http.ResponseWriter, r *http.Request) {
	printInfo(w, r)
	params := r.URL.Query()
	fmt.Println(params)

	name := params.Get("name")
	value := params.Get("value")

	createCookie(w, r, name, value)

	s := fmt.Sprintf("name = %v, value = %v\n", name, value)
	write(s, w)
}

func handleGetCookie(w http.ResponseWriter, r *http.Request) {
	printInfo(w, r)
	fmt.Println("Getting cookie")
	params := r.URL.Query()

	name := params.Get("name")
	value := getCookieValue(w, r, name)

	s := fmt.Sprintf("name = %v, value = %v\n", name, value)
	write(s, w)
}

func handleDeleteCookie(w http.ResponseWriter, r *http.Request) {
	printInfo(w, r)
	params := r.URL.Query()

	name := params.Get("name")
	deleteCookie(w, r, name)

	s := fmt.Sprintf("deleted name = %v\n", name)
	write(s, w)
}

func createCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	expires := time.Now().AddDate(1, 0, 0) // Expires one year from now
	c := http.Cookie{
		Name:    name,
		Value:   value,
		MaxAge:  360000,
		Expires: expires,
	}
	http.SetCookie(w, &c)
}

func getCookieValue(w http.ResponseWriter, r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		http.Error(w, "Cookie not found for "+name, http.StatusNotFound)
		return ""
	}

	value := c.Value
	s := fmt.Sprintf("Cookie: name = %v, value = %v\n", name, value)
	w.Write([]byte(s))
	return value
}

func deleteCookie(w http.ResponseWriter, r *http.Request, name string) {

	c := http.Cookie{
		Name:   name,
		MaxAge: -1,
	}
	http.SetCookie(w, &c)

	s := fmt.Sprintf("Cookie with name = %v deleted\n", name)
	write(s, w)
}

func write(s string, w http.ResponseWriter) {
	_, err := w.Write([]byte(s))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now(), "|", r.URL.Path, "|", r.Method, "|")
}

func writeResponse(w http.ResponseWriter, response interface{}) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	_, err = w.Write(responseBytes)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}