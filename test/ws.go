package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func sendDataToClient(client *websocket.Conn, data []byte) {
	mut.Lock()
	err := client.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	}
	mut.Unlock()
}
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	// Add client to clients map
	clients[conn] = true
	// send logs to client
	for k, v := range serviceLogs {
		for _, log_ := range v {
			data := map[string]interface{}{
				"method":  "service",
				"service": k,
				"output":  log_,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("%s json marshal error: %s\n", k, err)
				break
			}
			mut.Lock()
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Printf("%s write error: %s\n", k, err)
				mut.Unlock()
				continue
			}
			mut.Unlock()
		}
	}

	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// Broadcast message to all clients
		for c := range clients {
			mut.Lock()
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("write error:", err)
				delete(clients, c)
				mut.Unlock()
				break
			}
			mut.Unlock()
		}
	}

	delete(clients, conn)
}
