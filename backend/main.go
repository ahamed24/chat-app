package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket Upgrade error: %v", err)
		return
	}
	defer ws.Close()

	messageType, message, err := ws.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
		return
	}

	fmt.Println("Received: ", string(message))

	err = ws.WriteMessage(messageType, message)
	if err != nil {
		fmt.Println("write error: ", err)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)

	http.HandleFunc("/ws", websocketHandler)

	fmt.Println("Server running on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server")
	}

}
