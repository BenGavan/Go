package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"time"

	//"github.com/gorilla/websocket"
)

/*
A standard practice is to have  one go routine handling all the incoming and one handling all the writing to the sockets.
*/

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/v1/ws", handleV1WS)
	http.HandleFunc("/v2/ws", handleV2WS)
	http.HandleFunc("/v3/ws", handleV3WS)
	http.HandleFunc("/v4/ws", handleV4WS)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hey, this is the home page of the socket testing")
	if err != nil {
		fmt.Printf("Error service /: %v", err)
		return
	}
}

func handleV1WS(w http.ResponseWriter, r *http.Request) {
	printInfo(w, r)
	var conn, _ = upgrader.Upgrade(w, r, nil) // Upgrades the incoming http request to a websocket.

	go func(conn *websocket.Conn) {
		for {
			mType, msg, _ := conn.ReadMessage()

			err := conn.WriteMessage(mType, msg)
			if err != nil {
				fmt.Printf("Error writing to websocket: %v\n", err)
				continue
			}
		}
	}(conn)
}

func handleV2WS(w http.ResponseWriter, r *http.Request) {
	printInfo(w, r)
	conn, _ := upgrader.Upgrade(w, r, nil)

	go func(conn *websocket.Conn) {
		for {
			_, msg, _ := conn.ReadMessage()
			println(string(msg))

		}
	}(conn)
}

func handleV3WS(w http.ResponseWriter, r *http.Request) {
	printInfo(w, r)
	conn, _ := upgrader.Upgrade(w, r, nil)

	go func(conn *websocket.Conn) {
		ch := time.Tick(5 * time.Second)

		for range ch {
			conn.WriteJSON(myStruct{
				Username:  "ben_gavan",
				FirstName: "Ben",
				LastName:  "Gavan",
			})
		}
	}(conn)
}

/*
var ws = new WebSocket("ws://localhost:8000/v4/ws")

ws.addEventListener("message", function(e) {console.log(e);});
 */

func handleV4WS(w http.ResponseWriter, r *http.Request) {
	printInfo(w, r)
	conn, _ := upgrader.Upgrade(w, r, nil)

	go func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
			}
		}
	}(conn)

	go func(conn *websocket.Conn) {
		ch := time.Tick(5 * time.Second)

		for range ch {
			conn.WriteJSON(myStruct{
				Username:  "ben_gavan",
				FirstName: "Ben",
				LastName:  "Gavan",
			})
		}
	}(conn)
}

func printInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v | %v | %v\n", time.Now(), r.URL.Path, r.Method)
}

type myStruct struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
