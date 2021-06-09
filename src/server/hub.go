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

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Collection of websocket connections
	connections map[string]*websocket.Conn
}

func newHub() *Hub {
	return &Hub{
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
			log.Println("new subject registration, sending verification details...")

		case client := <-h.unregister:
			log.Printf("subject %s unregister\n", client.ID)
			client.saveClientData()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)  // delete active client when disconnect
			}
		}
	}
}
