package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string] *Client

	// Inbound messages from the clients.
	msg chan message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	connections map[string] *websocket.Conn
}

func newHub() *Hub {
	return &Hub{
		msg:        make(chan message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string] *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Println("new subject registration, sending verification details")
			h.clients[client.ID] = client

			jsonFile, err := os.Open(subjectsIDFilePath)
			if err != nil {
				log.Fatalf("error in reading files! %s", err)
			} else {
				log.Println("finish loading user verification information")
			}
			// read our opened jsonFile as a byte array.
			byteValue, _ := ioutil.ReadAll(jsonFile)

			msg := map[string]interface{}{}
			err = json.Unmarshal([]byte(strings.TrimSuffix(string(byteValue), "\r\n")), &msg)
			if err != nil {
				log.Fatalf("error in unmarshal IDs! %s", err)
			}

			message := msgToSubject{"verified_users", msg}

			err = client.conn.WriteJSON(message) // send verified users information to the newly registered client

			if err != nil {
				log.Fatalf("error in sending verified users! %s", err)
			}

		case client := <-h.unregister:
			log.Printf("subject %s unregister\n", client)
			if _, ok := h.clients[client.ID]; ok {
				delete(h.connections, client.ID)
				close(client.hubToClientMsg)
			}

		case message := <-h.msg:
			if client, ok := h.clients[message.ID]; ok {
				select {
				case client.hubToClientMsg <- message.data:
				default:
					close(client.hubToClientMsg)
					delete(h.connections, client.ID)
				}
			}
			//json.Unmarshal(message.data, &topLayerMsg) // decode the top layer incoming msg
		}
	}
}
