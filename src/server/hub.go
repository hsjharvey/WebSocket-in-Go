package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	msg chan message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	connections map[string]*websocket.Conn
}

func newHub() *Hub {
	return &Hub{
		msg:        make(chan message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
			log.Println("new subject registration, sending verification details")

			message := msgToSubject{"verified_users", subjectIDs}
			err := client.conn.WriteJSON(message) // send verified users information to the newly registered client
			if err != nil {
				log.Fatalf("error in sending verified users! %s", err)
			}

		case client := <-h.unregister:
			log.Printf("subject %s unregister\n", client.ID)
			client.saveSubjectData()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.connections, client.ID)
			}
		}
	}
}
