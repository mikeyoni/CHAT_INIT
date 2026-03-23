package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{

	ReadBufferSize:  2040,
	WriteBufferSize: 2040,

	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clientes   = make(map[string]*websocket.Conn)
	clientesMu sync.Mutex
)

func handlecannection(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("user")

	if r.URL.Path != "/chat" {
		fmt.Printf("\n Server error : %v ", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Printf("\n Cannection problem : %v ", err)
		return
	}

	defer conn.Close()

	clientesMu.Lock()

	clientes[username] = conn

	clientesMu.Unlock()

	fmt.Printf(" \n the %s has join the server \n ", username)

	defer func() {
		clientesMu.Lock()

		delete(clientes, username)

		clientesMu.Unlock()
		fmt.Printf("\n %s has left the server \n", username)
	}()

	for {

		_, massages, err := conn.ReadMessage()
		if err != nil {
			break
		}
		messagess := fmt.Sprintf("[%s] : > %s ", username, string(massages))
		fmt.Println(messagess)
		broadcastee(messagess, username)
	}

}

func broadcastee(msg string, sendere string) {
	clientesMu.Lock()
	defer clientesMu.Unlock()

	for username, conn := range clientes {

		if username != sendere {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				conn.Close()
				delete(clientes, username)
		
			}
		}
	}

}

func main() {

	http.HandleFunc("/chat", handlecannection)
	fmt.Printf("\n Server is started on : 6060 \n")
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		fmt.Printf(" server is not started : ", err)
		return
	}

	fmt.Printf(" \n hello ewvery one ")
}
