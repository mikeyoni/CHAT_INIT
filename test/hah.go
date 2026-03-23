package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgreader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	// The "Address Book" of all connected pirates
	clients   = make(map[string]*websocket.Conn)
	clientsMu sync.Mutex
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 1. Get the username from the URL (e.g., /ws?user=Mikey)
	username := r.URL.Query().Get("user")
	if username == "" {
		return
	}

	conn, err := upgreader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 2. Add the new pirate to our Address Book
	clientsMu.Lock()
	clients[username] = conn
	clientsMu.Unlock()

	fmt.Printf("--- %s has boarded the ship! ---\n", username)

	// 3. Remove them when they leave
	defer func() {
		clientsMu.Lock()
		delete(clients, username)
		clientsMu.Unlock()
		fmt.Printf("--- %s has left the ship ---\n", username)
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// LOGIC: Broadcast the message to EVERYONE else
		message := fmt.Sprintf("[%s]: %s", username, string(p))
		fmt.Println(message)

		broadcast(message, username)
	}
}

// 4. The "Postman" function
func broadcast(msg string, sender string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for username, conn := range clients {
		// Don't send the message back to the person who sent it!
		if username != sender {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Printf("Could not send to %s, closing conn\n", username)
				conn.Close()
				delete(clients, username)
			}
		}
	}
}

func maineeeeeeeeeeeeeeeee() { // deforming the thing ok 
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Pirate Hub started on :5050")
	http.ListenAndServe(":5050", nil)
}
