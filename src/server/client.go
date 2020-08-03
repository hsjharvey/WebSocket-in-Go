package wsserver

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	hubToClientMsg chan []byte

	ID string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readMsg() {
	topLayerMsg := make(map[string]interface{})
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// c.hub.msg <- message

		json.Unmarshal(message, &topLayerMsg) // decode the top layer incoming msg
		if topLayerMsg["msg_type"] == "game_version" {
			if topLayerMsg["msg"] == "practice" {
				log.Println("game_version: practice")
			} else if topLayerMsg["msg"] == "game_version" {
				log.Println("game_version: actual")
			}
		} else if topLayerMsg["msg_type"] == "register" {
			log.Println(topLayerMsg["msg"])
			if topLayerMsg["msg"] == "not_verified" {
				log.Println("subject not verified")
			} else {
				log.Printf("subject %s verified\n")
			}
		} else if topLayerMsg["msg_type"] == "play_data" {
			log.Println(topLayerMsg["msg"])
		} else if topLayerMsg["msg_type"] == "game_information" {
			log.Println(topLayerMsg["msg"])
		} else {
			log.Fatalf("Warning, unsupported event: %s", topLayerMsg["msg_type"])
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, hubToClientMsg: make(chan []byte, 1024*256), ID: "default_" + string(rand.Intn(1000))}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.readMsg()
}