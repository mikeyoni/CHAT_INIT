package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main333() {
	// 1. Get the Pirate Name first!
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Pirate Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name) // Remove the "Enter" key newline

	// 2. Build the URL with the name as a query parameter
	// This tells the server exactly who is boarding
	url := fmt.Sprintf("ws://localhost:5050/ws?user=%s", name)

	fmt.Println("Attempting to board the ship at:", url)

	// 3. Dial the Server
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Handshake Failed! Did you start the server? Error:", err)
	}
	defer conn.Close()

	fmt.Printf("--- WELCOME ABOARD, %s! ---\n", strings.ToUpper(name))

	// 4. Background Worker: Listen for others talking
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("\n[System] Connection lost.")
				return
			}
			// \r clears the "> " prompt so the message looks clean
			fmt.Printf("\r%s\n> ", string(message))
		}
	}()

	// 5. Foreground Worker: Your typing station
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				continue
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				fmt.Println("Send Error:", err)
				break
			}
		}
	}
}