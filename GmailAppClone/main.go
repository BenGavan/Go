package main

import (
	"encoding/json"
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"net"
	"net/http"
	"time"
)

type user struct {
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"address,omitempty"`
}

type email struct {
	Sender         user   `json:"sender,omitempty"`
	Date           string `json:"date,omitempty"`
	Subject        string `json:"subject,omitempty"`
	Content        string `json:"content,omitempty"`
	IsStared       bool   `json:"isStared"`
	IsImportant    bool   `json:"isImportant"`
	HasBeenOpened  bool   `json:"hasBeenOpened"`
	HasAttachments bool   `json:"hasAttachments"`
}

func main() {
	fmt.Println("Hey")
	getIP()
	http.HandleFunc("/", handleEmails)
	http.HandleFunc("/string", handleStringDownload)
	http.HandleFunc("/favicon.ico", handleFavicon)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func getIP() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Error(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		log.Error(err)

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			fmt.Println(ip)
		}
	}
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
}

func handleEmails(w http.ResponseWriter, r *http.Request) {
	printInfo(r)

	switch r.Method {
	case http.MethodGet:
		u := user{"Stanford University", "about@stanford.edu"}
		e := email{Sender:u, Date:"7 Sep", Subject:"Physics", Content:"This is the email content", IsStared:false, IsImportant:false, HasBeenOpened:false, HasAttachments:false}
		e2 := email{Sender:u, Date:"7 Sep", Subject:"Physics", Content:"This is the email content", IsStared:true, IsImportant:true, HasBeenOpened:false, HasAttachments:true}
		var es []email
		es = append(append(es, e), e2)

		err := json.NewEncoder(w).Encode(es)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleStringDownload(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	outString := "Hey, How's it going??"
	w.Write([]byte(outString))
}

func printInfo(r *http.Request) {
	fmt.Println(time.Now(), "|", r.URL.Path, "|", r.Method, "|")
}
