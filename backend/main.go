package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var broadcaster = make(chan ChatMessage)

var clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket Upgrade error: %v", err)
		return
	}

	clients[ws] = true
	log.Println("New client connected")

	for {
		var msg ChatMessage

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading the message: ", err)
			delete(clients, ws)
			break
		}

		broadcaster <- msg
	}

}

func handleMessages() {
	for {
		msg := <-broadcaster

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error writing the message: ", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {

	http.HandleFunc("/", handler)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Server running on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server")
	}

}
